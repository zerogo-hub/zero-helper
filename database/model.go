package database

import (
	"gorm.io/gorm"
)

// Model 替代 gorm.Model，无 id，需要自定义主键
// 时间为 秒，
type Model struct {
	CreatedAt int64          `gorm:"comment:'记录创建时间'"`
	UpdatedAt int64          `gorm:"comment:'最后一次更新时间'"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:'软删除时间'"`
}

// ModelID 替代 gorm.Model，带自增 id
type ModelID struct {
	ID uint `gorm:"primaryKey;AUTO_INCREMENT"`
	Model
}
