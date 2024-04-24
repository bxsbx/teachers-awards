package dm

import (
	"context"
	"database/sql"
	_ "dm"
	"fmt"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/atomic"
	"io"
	"strings"
)

func init() {
	database.Register("dm", &DM{})
}

var DefaultMigrationsTable = "schema_migrations"

var (
	ErrNilConfig = fmt.Errorf("no config")
	ErrNoSchema  = fmt.Errorf("no schema")
)

type Config struct {
	MigrationsTable string
	SchemaName      string
}

type DM struct {
	conn     *sql.Conn
	db       *sql.DB
	isLocked atomic.Bool

	config *Config
}

// connection instance must have `multiStatements` set to true
func WithConnection(ctx context.Context, conn *sql.Conn, config *Config) (*DM, error) {
	if config == nil {
		return nil, ErrNilConfig
	}

	if err := conn.PingContext(ctx); err != nil {
		return nil, err
	}

	mx := &DM{
		conn:   conn,
		db:     nil,
		config: config,
	}

	if len(config.MigrationsTable) == 0 {
		config.MigrationsTable = DefaultMigrationsTable
	}

	if err := mx.ensureTableSpace(); err != nil {
		return nil, err
	}

	if err := mx.ensureSchema(); err != nil {
		return nil, err
	}

	if err := mx.ensureVersionTable(); err != nil {
		return nil, err
	}

	return mx, nil
}

// instance must have `multiStatements` set to true
func WithInstance(instance *sql.DB, config *Config) (database.Driver, error) {
	ctx := context.Background()

	if err := instance.Ping(); err != nil {
		return nil, err
	}

	conn, err := instance.Conn(ctx)
	if err != nil {
		return nil, err
	}

	mx, err := WithConnection(ctx, conn, config)
	if err != nil {
		return nil, err
	}
	if config.SchemaName == "" {
		query := `SELECT SYS_CONTEXT ('userenv', 'current_schema') FROM DUAL;`
		var schemaName string
		if err := instance.QueryRow(query).Scan(&schemaName); err != nil {
			return nil, &database.Error{OrigErr: err, Query: []byte(query)}
		}

		if len(schemaName) == 0 {
			return nil, ErrNoSchema
		}

		config.SchemaName = schemaName
	}

	mx.db = instance

	return mx, nil
}

func (m *DM) Open(dns string) (database.Driver, error) {
	cfg, err := ParseDSN(dns)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("dm", fmt.Sprintf("dm://%s:%s@%s", cfg.User, cfg.Passwd, cfg.Addr))
	if err != nil {
		return nil, err
	}

	mx, err := WithInstance(db, &Config{
		MigrationsTable: cfg.Params["x-migrations-table"],
		SchemaName:      cfg.SchemaName,
	})
	if err != nil {
		return nil, err
	}
	return mx, nil
}

func (m *DM) Close() error {
	connErr := m.conn.Close()
	var dbErr error
	if m.db != nil {
		dbErr = m.db.Close()
	}

	if connErr != nil || dbErr != nil {
		return fmt.Errorf("conn: %v, db: %v", connErr, dbErr)
	}
	return nil
}

// DBMS_LOCK 达梦的锁机制包
func (m *DM) Lock() error {
	return database.CasRestoreOnErr(&m.isLocked, false, true, database.ErrLocked, func() error {
		aid := generateAdvisoryLockId(fmt.Sprintf("%s:%s", m.config.SchemaName, m.config.MigrationsTable))

		query := "SELECT DBMS_LOCK.REQUEST(?);"
		var success int
		if err := m.conn.QueryRowContext(context.Background(), query, aid).Scan(&success); err != nil {
			return &database.Error{OrigErr: err, Err: "try lock failed", Query: []byte(query)}
		}

		if success != 0 {
			return database.ErrLocked
		}

		return nil
	})
}

func (m *DM) Unlock() error {
	return database.CasRestoreOnErr(&m.isLocked, true, false, database.ErrNotLocked, func() error {
		aid := generateAdvisoryLockId(fmt.Sprintf("%s:%s", m.config.SchemaName, m.config.MigrationsTable))

		query := `SELECT DBMS_LOCK.RELEASE(?);`
		var success int
		if err := m.conn.QueryRowContext(context.Background(), query, aid).Scan(&success); err != nil {
			return &database.Error{OrigErr: err, Query: []byte(query)}
		}

		if success != 0 {
			return database.ErrNotLocked
		}
		return nil
	})
}

func (m *DM) Run(migration io.Reader) error {
	migr, err := io.ReadAll(migration)
	if err != nil {
		return err
	}

	ctx := context.Background()

	query := string(migr[:])
	split := strings.Split(query, ";")
	tx, err := m.conn.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return &database.Error{OrigErr: err, Err: "transaction start failed"}
	}
	for _, v := range split {
		trim := strings.TrimSpace(v)
		if len(trim) > 0 {
			if _, err := m.conn.ExecContext(ctx, v); err != nil {
				if errRollback := tx.Rollback(); errRollback != nil {
					err = multierror.Append(err, errRollback)
				}
				return &database.Error{OrigErr: err, Query: []byte(query)}
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return &database.Error{OrigErr: err, Err: "transaction commit failed"}
	}
	return nil
}

func (m *DM) SetVersion(version int, dirty bool) error {
	tx, err := m.conn.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return &database.Error{OrigErr: err, Err: "transaction start failed"}
	}

	query := "DELETE FROM " + quoteIdentifier(m.config.SchemaName) + `.` + quoteIdentifier(m.config.MigrationsTable)
	if _, err := tx.ExecContext(context.Background(), query); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			err = multierror.Append(err, errRollback)
		}
		return &database.Error{OrigErr: err, Query: []byte(query)}
	}

	// Also re-write the schema version for nil dirty versions to prevent
	// empty schema version for failed down migration on the first migration
	// See: https://github.com/golang-migrate/migrate/issues/330
	if version >= 0 || (version == database.NilVersion && dirty) {
		query := "INSERT INTO " + quoteIdentifier(m.config.SchemaName) + `.` + quoteIdentifier(m.config.MigrationsTable) + " (\"version\", \"dirty\") VALUES (?, ?)"
		if _, err := tx.ExecContext(context.Background(), query, version, dirty); err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				err = multierror.Append(err, errRollback)
			}
			return &database.Error{OrigErr: err, Query: []byte(query)}
		}
	}

	if err := tx.Commit(); err != nil {
		return &database.Error{OrigErr: err, Err: "transaction commit failed"}
	}

	return nil
}

func (m *DM) Version() (version int, dirty bool, err error) {
	query := "SELECT \"version\", \"dirty\" FROM " + quoteIdentifier(m.config.SchemaName) + `.` + quoteIdentifier(m.config.MigrationsTable) + " LIMIT 1"
	err = m.conn.QueryRowContext(context.Background(), query).Scan(&version, &dirty)
	switch {
	case err == sql.ErrNoRows:
		return database.NilVersion, false, nil

	case err != nil:
		return 0, false, &database.Error{OrigErr: err, Query: []byte(query)}

	default:
		return version, dirty, nil
	}
}

func (m *DM) Drop() (err error) {
	query := "SELECT TABLE_NAME FROM ALL_TABLES WHERE OWNER ='" + m.config.SchemaName + "';"
	tables, err := m.conn.QueryContext(context.Background(), query)
	if err != nil {
		return &database.Error{OrigErr: err, Query: []byte(query)}
	}
	defer func() {
		if errClose := tables.Close(); errClose != nil {
			err = multierror.Append(err, errClose)
		}
	}()

	// delete one table after another
	tableNames := make([]string, 0)
	for tables.Next() {
		var tableName string
		if err := tables.Scan(&tableName); err != nil {
			return err
		}
		if len(tableName) > 0 {
			tableNames = append(tableNames, tableName)
		}
	}
	if err := tables.Err(); err != nil {
		return &database.Error{OrigErr: err, Query: []byte(query)}
	}

	if len(tableNames) > 0 {
		// delete one by one ...
		for _, t := range tableNames {
			query = `DROP TABLE IF EXISTS ` + quoteIdentifier(m.config.SchemaName) + `.` + quoteIdentifier(t)
			if _, err := m.conn.ExecContext(context.Background(), query); err != nil {
				return &database.Error{OrigErr: err, Query: []byte(query)}
			}
		}
	}

	return nil
}

func (m *DM) ensureVersionTable() (err error) {
	// if not, create the empty migration table
	query := "CREATE TABLE IF NOT EXISTS " + quoteIdentifier(m.config.SchemaName) + `.` + quoteIdentifier(m.config.MigrationsTable) + " (\"version\" BIGINT NOT NULL,\"dirty\" TINYINT NOT NULL, NOT CLUSTER PRIMARY KEY(\"version\")) STORAGE(ON " + quoteIdentifier(m.config.SchemaName) + ", CLUSTERBTR);"
	if _, err := m.conn.ExecContext(context.Background(), query); err != nil {
		return &database.Error{OrigErr: err, Query: []byte(query)}
	}
	return nil
}

func (m *DM) ensureTableSpace() (err error) {
	query := "create tablespace IF NOT EXISTS " + quoteIdentifier(m.config.SchemaName) + " datafile '" + m.config.SchemaName + "' size 128 autoextend on maxsize 67108863 CACHE = NORMAL;"
	if _, err := m.conn.ExecContext(context.Background(), query); err != nil {
		return &database.Error{OrigErr: err, Query: []byte(query)}
	}
	return nil
}

func (m *DM) ensureSchema() (err error) {
	query := "CREATE SCHEMA " + quoteIdentifier(m.config.SchemaName) + " AUTHORIZATION \"SYSDBA\";"
	if _, err := m.conn.ExecContext(context.Background(), query); err != nil {
		return &database.Error{OrigErr: err, Query: []byte(query)}
	}
	return nil
}
