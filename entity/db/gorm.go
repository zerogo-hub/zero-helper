package db

import (
	"context"
	"errors"
	"reflect"

	"gorm.io/gorm"

	zerodatabase "github.com/zerogo-hub/zero-helper/database"
	zeroentity "github.com/zerogo-hub/zero-helper/entity"
	zeroutils "github.com/zerogo-hub/zero-helper/utils"
)

// wrapGorm 封装 gorm
type wrapGorm struct {
	db          zerodatabase.Database
	errNotFound error
}

func newGormRead(db zerodatabase.Database) zeroentity.WrapReadDB {
	return &wrapGorm{
		db:          db,
		errNotFound: gorm.ErrRecordNotFound,
	}
}

func NewGormReadF2(dbs ...zerodatabase.Database) []zeroentity.WrapReadDB {
	if len(dbs) == 0 {
		return nil
	}

	if len(dbs) == 1 {
		return []zeroentity.WrapReadDB{newGormRead(dbs[0])}
	}

	n := len(dbs)
	f2n := zeroutils.F2(n)
	results := make([]zeroentity.WrapReadDB, 0, f2n)

	for _, db := range dbs {
		results = append(results, newGormRead(db))
	}

	if f2n > n {
		for _, db := range dbs {
			results = append(results, newGormRead(db))
			if len(results) > f2n {
				break
			}
		}
	}

	return results
}

func NewGormWrite(db zerodatabase.Database) zeroentity.WrapWriteDB {
	return &wrapGorm{
		db:          db,
		errNotFound: gorm.ErrRecordNotFound,
	}
}

func (w *wrapGorm) Get(out interface{}, id uint64) error {
	return w.db.DB().First(out, id).Error
}

func (w *wrapGorm) MGet(out interface{}, ids ...uint64) ([]uint64, []interface{}, error) {
	tx := w.db.DB().Find(out, ids)
	if tx.Error != nil {
		return nil, nil, tx.Error
	}
	return w.parseQueryResults(tx, out)
}

// parseQueryResults 解析 MGet 查询结果
// return: 主键集合, 数据集合, error
func (w *wrapGorm) parseQueryResults(tx *gorm.DB, out interface{}) ([]uint64, []interface{}, error) {

	// 1
	values := make([]reflect.Value, 0)
	destValue := reflect.Indirect(reflect.ValueOf(out))
	switch destValue.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < destValue.Len(); i++ {
			elem := destValue.Index(i)
			values = append(values, elem)
		}
	case reflect.Struct:
		values = append(values, destValue)
	}

	// 2
	var valueOf func(context.Context, reflect.Value) (value interface{}, zero bool) = nil
	if tx.Statement.Schema != nil {
		for _, field := range tx.Statement.Schema.Fields {
			if field.PrimaryKey {
				valueOf = field.ValueOf
				break
			}
		}
	}

	// 3
	primaryKeys := make([]uint64, 0, len(values))
	objects := make([]interface{}, 0, len(values))
	for _, elemValue := range values {
		if valueOf != nil {
			primaryKey, isZero := valueOf(context.Background(), elemValue)
			if isZero {
				continue
			}
			u64, ok := primaryKey.(uint64)
			if !ok {
				return nil, nil, errors.New("primaryKey must be uint64")
			}

			primaryKeys = append(primaryKeys, u64)
		}
		objects = append(objects, elemValue.Interface())
	}

	return primaryKeys, objects, nil
}

func (w *wrapGorm) Update(in interface{}) error {
	return w.db.DB().Save(in).Error
}

func (w *wrapGorm) Delete(model interface{}, id uint64) error {
	return w.db.DB().Unscoped().Delete(model, id).Error
}

func (w *wrapGorm) MDelete(model interface{}, ids ...uint64) error {
	return nil
}

func (w *wrapGorm) ErrNotFound() error {
	return w.errNotFound
}
