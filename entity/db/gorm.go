package db

import (
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

func (w *wrapGorm) Get(id uint64, out interface{}) error {
	return w.db.DB().First(out, id).Error
}

func (w *wrapGorm) MGet(out interface{}, ids ...uint64) error {
	return w.db.DB().Find(out, ids).Error
}

func (w *wrapGorm) Update(in interface{}) error {
	return w.db.DB().Save(in).Error
}

func (w *wrapGorm) Delete(id uint64, model interface{}) error {
	return w.db.DB().Unscoped().Delete(model, id).Error
}

func (w *wrapGorm) ErrNotFound() error {
	return w.errNotFound
}
