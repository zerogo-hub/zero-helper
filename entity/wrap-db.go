package entity

import (
	zerodatabase "github.com/zerogo-hub/zero-helper/database"
	"gorm.io/gorm"
)

// wrapGorm 封装 gorm
type wrapGorm struct {
	db          zerodatabase.Database
	errNotFound error
}

func NewWrapDB(db zerodatabase.Database) WrapDB {
	return &wrapGorm{
		db:          db,
		errNotFound: gorm.ErrRecordNotFound,
	}
}

func (w *wrapGorm) Get(id uint64, out interface{}) error {
	return w.db.DB().First(out, id).Error
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
