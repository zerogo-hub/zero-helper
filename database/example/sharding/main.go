package main

import (
	zerodatabase "github.com/zerogo-hub/zero-helper/database"
	zerosharding "github.com/zerogo-hub/zero-helper/database/sharding"
	zerologger "github.com/zerogo-hub/zero-helper/logger"
)

type Account struct {
	UUID     uint64 `gorm:"primaryKey;autoIncrement:false;comment:'账号唯一标识符'" json:"uuid"`
	Username string `gorm:"size:256;not null;unique;comment:'用户名'" json:"username"`
	Password string `gorm:"size:32;not null;comment:'加密后的密码'" json:"-"`
	Age      int16  `gorm:"default:0;comment:'岁数'"  json:"age"`

	zerodatabase.Model
}

func main() {
	logger := zerologger.NewSampleLogger()

	db := connectDB("test1", logger)
	dbs := []zerodatabase.Database{
		db,
		connectDB("test2", logger),
		connectDB("test3", logger),
	}

	if err := db.DB().Use(zerosharding.Register(
		&zerosharding.Config{
			// 创建 10 个表
			NumberOfShards: 10,
			Tables: []any{
				// 有哪些表需要分表
				Account{},
			},
			// 分布在这些数据库中
			Databases: dbs,
		},
	)); err != nil {
		logger.Fatalf("register sharding failed, err: %s", err.Error())
		return
	}

	// 添加数据
	add(db, &Account{UUID: 1, Username: "zero1", Password: "123456", Age: 11})
	add(db, &Account{UUID: 2, Username: "zero2", Password: "123456", Age: 12})
	add(db, &Account{UUID: 3, Username: "zero3", Password: "123456", Age: 13})
	add(db, &Account{UUID: 4, Username: "zero4", Password: "123456", Age: 14})

	for _, db := range dbs {
		db.Close()
	}
}

func add(db zerodatabase.Database, account *Account) {

}

func connectDB(name string, logger zerologger.Logger) zerodatabase.Database {
	db := zerodatabase.New(
		zerodatabase.WithUsername("root"),
		zerodatabase.WithPassword("123456"),
		zerodatabase.WithHost("127.0.0.1"),
		zerodatabase.WithPort(3306),
		zerodatabase.WithDBName(name),
	)

	if err := db.Open(); err != nil {
		logger.Fatalf("open db failed, err: %s", err.Error())
		return nil
	}

	return db
}
