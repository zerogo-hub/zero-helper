package database

import (
	"errors"

	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database 数据库
type Database interface {
	config() *config
	Open() error
	DB() *gorm.DB
	AutoMigrate(values ...interface{}) (Database, error)
	Close()
}

type database struct {
	conf *config
	db   *gorm.DB
}

// New ..
func New(opts ...Option) Database {
	d := &database{
		conf: defaultConfig(),
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

// config ..
func (d *database) config() *config {
	return d.conf
}

// Open ..
func (d *database) Open() error {
	dsn, err := d.dsn()
	if err != nil {
		return err
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if d.conf.maxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(d.conf.maxIdleConns)
	}
	if d.conf.maxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(d.conf.maxOpenConns)
	}
	if d.conf.maxConnLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(d.conf.maxConnLifeTime) * time.Second)
	}

	d.db = db

	return nil
}

// DB ..
func (d *database) DB() *gorm.DB {
	return d.db
}

// AutoMigrate 迁移
func (d *database) AutoMigrate(values ...interface{}) (Database, error) {
	if err := d.db.AutoMigrate(values...); err != nil {
		return nil, err
	}
	return d, nil
}

// Close ..
func (d *database) Close() {
	sqlDB, _ := d.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func (d *database) dsn() (string, error) {
	c := d.conf
	s := ""

	switch d.conf.dialect {
	case "mysql":
		// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
		s = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.username, c.password, c.host, c.port, c.dbname)
	case "postgres":
		s = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			c.host, c.username, c.dbname, c.password)
	case "sqlite3":
		s = fmt.Sprintf("%s/%s", c.host, c.dbname)
	}

	if s == "" {
		return "", errors.New("invalid dialect")
	}

	return s, nil
}
