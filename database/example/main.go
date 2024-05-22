package main

import (
	"strconv"

	zerodatabase "github.com/zerogo-hub/zero-helper/database"
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

	// 如果不存在，创建表
	if err := db.DB().AutoMigrate(&Account{}); err != nil {
		return
	}

	// 新增一条数据
	snowflake, _ := zerorandom.NewBit46Snowflake(1)
	uuid, _ := snowflake.UnsafeUUID()

	db.DB().Create(&Account{
		UUID:     uuid,
		Username: "zero_" + strconv.FormatUint(uuid, 10),
		Password: "e10adc3949ba59abbe56e057f20f883e",
		Age:      12,
	})

	// 查询数据
	var account Account
	if err := db.DB().First(&account, uuid).Error; err != nil {
		logger.Fatalf("find uuid: %d failed, err: %s", uuid, err.Error())
	}
	logger.Infof("account.UUID: %d", account.UUID)

	// 修改数据
	account.Age = 13
	if err := db.DB().Save(&account).Error; err != nil {
		logger.Fatalf("update uuid: %d failed, err: %s", uuid, err.Error())
	}
	logger.Infof("account.Age: %d", account.Age)

	// 删除数据
	if err := db.DB().Delete(&account).Error; err != nil {
		logger.Fatalf("delete uuid: %d failed, err: %s", uuid, err.Error())
	}
	logger.Infof("delete uuid: %d success", uuid)

	db.Close()
}
