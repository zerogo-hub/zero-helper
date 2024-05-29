package main

import (
	"errors"
	"time"

	zerocmsgpack "github.com/zerogo-hub/zero-helper/codec/msgpack"
	zerodatabase "github.com/zerogo-hub/zero-helper/database"
	zeroentity "github.com/zerogo-hub/zero-helper/entity"

	zeroentitycache "github.com/zerogo-hub/zero-helper/entity/cache"
	zeroentitydb "github.com/zerogo-hub/zero-helper/entity/db"
	zerologger "github.com/zerogo-hub/zero-helper/logger"
	zerorandom "github.com/zerogo-hub/zero-helper/random"
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

	// 数据库
	db := zerodatabase.New(
		zerodatabase.WithUsername("root"),
		zerodatabase.WithPassword("123456"),
		zerodatabase.WithHost("127.0.0.1"),
		zerodatabase.WithPort(3306),
		zerodatabase.WithDBName("test"),
	)
	if err := db.Open(); err != nil {
		logger.Fatalf("open db failed, err: %s", err.Error())
		return
	}
	testReady(db)

	// 数据统计
	st := zeroentity.NewStat("example", func(localHit, localMiss, remoteHit, remoteMiss, dbFails, customFails uint64) {
		logger.Infof("cache: %s, localHit: %d, localMiss: %d, remoteHit: %d, remoteMiss: %d, db_fails: %d, custom_fails: %d",
			"example", localHit, localMiss, remoteHit, remoteMiss, dbFails, customFails)
	})

	// 创建实例管理器
	e := zeroentity.New(
		st,
		logger,
		zerocmsgpack.New(),
	)
	e.WithReadDB(zeroentitydb.NewGormReadF2(db)...).WithWriteDB(zeroentitydb.NewGormWrite(db))
	e.WithLocalCache(zeroentitycache.NewBigCache(10 * time.Minute))
	e.WithTimeout(30 * time.Minute)
	e.Build()

	// 测试
	testQueryNotFound(e)
	account, err := testQueryExist(e, logger)
	if err != nil {
		return
	}
	if err := testUpdate(e, logger, account); err != nil {
		return
	}

	// 测试批量操作
	testMulti(db, e, logger)
}

func testReady(database zerodatabase.Database) {
	db := database.DB()
	if err := db.AutoMigrate(&Account{}); err != nil {
		return
	}

	account := Account{UUID: 1}
	db.First(&account)
	if account.CreatedAt > 0 {
		return
	}

	db.Create(&Account{
		UUID:     1,
		Username: "zero1",
		Password: "e10adc3949ba59abbe56e057f20f883e",
		Age:      12,
	})
	db.Create(&Account{
		UUID:     2,
		Username: "zero2",
		Password: "e10adc3949ba59abbe56e057f20f883e",
		Age:      18,
	})
}

func testQueryNotFound(e zeroentity.Entity) {
	var account Account
	// 查找一个不存在的数据，第一次查找不存在，设置短期缓存
	_ = e.Get(&account, 3)
	// 查找一个不存在的数据，第二次从缓存中获取
	_ = e.Get(&account, 3)
}

func testQueryExist(e zeroentity.Entity, logger zerologger.Logger) (*Account, error) {
	// 查找一个已存在的数据，第一次从数据库中获取
	var account Account

	if err := e.Get(&account, 1); err != nil {
		logger.Errorf("get account failed, err: %s", err.Error())
		return nil, err
	}
	// 查找一个已存在的数据，第二次从缓存中获取
	if err := e.Get(&account, 1); err != nil {
		logger.Errorf("get account failed, err: %s", err.Error())
		return nil, err
	}

	return &account, nil
}

func testUpdate(e zeroentity.Entity, logger zerologger.Logger, account *Account) error {
	newAge := int16(zerorandom.Int(18, 100))
	account.Age = newAge
	if err := e.Update(account, 1); err != nil {
		logger.Errorf("update account failed, err: %s", err.Error())
		return err
	}

	var account2 Account
	if err := e.Get(&account2, 1); err != nil {
		logger.Errorf("get account2 failed, err: %s", err.Error())
		return err
	}
	if account2.Age != newAge {
		logger.Errorf("get account2.age failed, age: %d", account2.Age)
		return errors.New("invalid age after update")
	}

	return nil
}

func testMulti(database zerodatabase.Database, e zeroentity.Entity, logger zerologger.Logger) {
	testMultiReady(database)
	testMultiQuery(e, logger)
}

func testMultiReady(database zerodatabase.Database) {
	db := database.DB()
	if err := db.AutoMigrate(&Account{}); err != nil {
		return
	}

	db.Create(&Account{
		UUID:     11,
		Username: "zero11",
		Password: "e10adc3949ba59abbe56e057f20f883e",
		Age:      12,
	})
	db.Create(&Account{
		UUID:     12,
		Username: "zero12",
		Password: "e10adc3949ba59abbe56e057f20f883e",
		Age:      18,
	})
}

func testMultiQuery(e zeroentity.Entity, logger zerologger.Logger) {
	var accounts []Account
	results, err := e.MGet(&accounts, 11, 12)
	if err != nil {
		logger.Errorf("[testMultiQuery] MGet failed, err: %s", err.Error())
		return
	}

	for _, bs := range results.Vals {
		var account Account
		if err := e.Unmarshal(bs, &account); err != nil {
			logger.Errorf("[testMultiQuery] Unmarshal failed: %s", err.Error())
		} else {
			logger.Infof("[testMultiQuery] from DB, UUID: %d, Username: %s", account.UUID, account.Username)
		}
	}

	results, err = e.MGet(&accounts, 11, 12)
	if err != nil {
		logger.Errorf("[testMultiQuery] MGet failed, err: %s", err.Error())
		return
	}

	for _, bs := range results.Vals {
		var account Account
		if err := e.Unmarshal(bs, &account); err != nil {
			logger.Errorf("[testMultiQuery] Unmarshal failed: %s", err.Error())
		} else {
			logger.Infof("[testMultiQuery] from cache, UUID: %d, Username: %s", account.UUID, account.Username)
		}
	}
}
