package dm

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

type DMsqlConfig struct {
	User       string // Username
	Passwd     string // Password (requires User)
	Addr       string // Network address (requires Net)
	SchemaName string
	ParseTime  bool              // Parse time values to time.Time
	Timeout    time.Duration     // Dial timeout
	Params     map[string]string // Connection parameters
}

func parseDSNParams(cfg *DMsqlConfig, params string) (err error) {
	for _, v := range strings.Split(params, "&") {
		param := strings.SplitN(v, "=", 2)
		if len(param) != 2 {
			continue
		}
		// cfg params
		switch value := param[1]; param[0] {
		// Disable INFILE allowlist / enable all files

		case "parseTime":
			var isBool bool
			cfg.ParseTime, isBool = readBool(value)
			if !isBool {
				return errors.New("invalid bool value: " + value)
			}
		// Dial Timeout
		case "timeout":
			cfg.Timeout, err = time.ParseDuration(value)
			if err != nil {
				return
			}
		case "schema":
			cfg.SchemaName = value
		default:
			// lazy init
			if cfg.Params == nil {
				cfg.Params = make(map[string]string)
			}

			if cfg.Params[param[0]], err = url.QueryUnescape(value); err != nil {
				return
			}
		}
	}
	return
}

// Returns the bool value of the input.
func readBool(input string) (value bool, valid bool) {
	switch input {
	case "1", "true", "TRUE", "True":
		return true, true
	case "0", "false", "FALSE", "False":
		return false, true
	}

	// Not a valid bool value
	return
}

func ParseDSN(dsn string) (cfg *DMsqlConfig, err error) {
	cfg = &DMsqlConfig{}
	j := 0
	for ; j < len(dsn)-1; j++ {
		if dsn[j] == '@' {
			// username[:password]
			// Find the first ':' in dsn[:j]
			for k := 0; k < j; k++ {
				if dsn[k] == ':' {
					cfg.User = dsn[:k]
					cfg.Passwd = dsn[k+1 : j]
					break
				}
			}
			break
		}
	}

	for i := j + 1; i < len(dsn); i++ {
		if dsn[i] == '?' {
			if err = parseDSNParams(cfg, dsn[i+1:]); err != nil {
				return
			}
			cfg.Addr = dsn[j+1 : i]
			break
		}
	}
	return
}

func quoteIdentifier(name string) string {
	end := strings.IndexRune(name, 0)
	if end > -1 {
		name = name[:end]
	}
	return `"` + strings.Replace(name, `"`, `""`, -1) + `"`
}

func generateAdvisoryLockId(databaseName string) int {
	sum := 0
	for _, v := range databaseName {
		sum += int(v)
	}
	return sum
}
