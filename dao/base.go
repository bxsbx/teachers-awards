package dao

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/common/errorz"
)

// 对于简短且不常用的sql可使用以下方法，如果sql过长或者复杂或者常用的尽可能地使用具体的方法，这么做也是为了对err进行包装

type BaseMysql struct {
	db     *gorm.DB
	appCtx context.Context
}

func (t *BaseMysql) Pluck(where interface{}, list interface{}, fields string) (err error) {
	db := t.db.WithContext(t.appCtx)
	err = db.Where(where).Pluck(fields, list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *BaseMysql) Find(where interface{}, list interface{}, fields ...string) (err error) {
	db := t.db.WithContext(t.appCtx)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	err = db.Where(where).Find(list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *BaseMysql) Count(where interface{}) (count int64, err error) {
	err = t.db.WithContext(t.appCtx).Where(where).Count(&count).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *BaseMysql) FindWithPageOrder(where interface{}, list interface{}, page, limit int, order string, fields ...string) (err error) {
	db := t.db.WithContext(t.appCtx)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	if page > 0 && limit > 0 {
		db = db.Offset((page - 1) * limit).Limit(limit)
	}
	err = db.Where(where).Order(order).Find(list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *BaseMysql) FindAndCountWithPageOrder(where interface{}, list interface{}, page, limit int, order string, fields ...string) (total int64, err error) {
	db := t.db.WithContext(t.appCtx)
	db = db.Where(where)
	err = db.Count(&total).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	if order != "" {
		db = db.Order(order)
	}
	if page > 0 && limit > 0 {
		db = db.Offset((page - 1) * limit).Limit(limit)
	}
	err = db.Find(list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *BaseMysql) First(where interface{}, row interface{}) (err error) {
	err = t.db.WithContext(t.appCtx).Where(where).First(row).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *BaseMysql) Create(obj interface{}) (err error) {
	err = t.db.WithContext(t.appCtx).Create(obj).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_CREATE_ERR, err)
	}
	return
}

func (t *BaseMysql) UpdateByWhere(where interface{}, update interface{}) (err error) {
	err = t.db.WithContext(t.appCtx).Where(where).Updates(update).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_UPDATE_ERR, err)
	}
	return
}

func (t *BaseMysql) DeleteByWhere(where interface{}, args ...interface{}) (err error) {
	err = t.db.WithContext(t.appCtx).Where(where, args...).Delete(nil).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_DELETE_ERR, err)
	}
	return
}

func (t *BaseMysql) BatchInsert(list interface{}) (err error) {
	err = t.db.WithContext(t.appCtx).CreateInBatches(list, 1000).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_CREATE_ERR, err)
	}
	return
}
