package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	zerocodecjson "github.com/zerogo-hub/zero-helper/codec/json"
	zerodatabase "github.com/zerogo-hub/zero-helper/database"
	zeroentity "github.com/zerogo-hub/zero-helper/entity"
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

	// 数据库
	db := zerodatabase.NewDatabase(
		zerodatabase.WithUsername("root"),
		zerodatabase.WithPassword("sNoJO6iqbBudAxUt"),
		zerodatabase.WithHost("127.0.0.1"),
		zerodatabase.WithPort(3306),
		zerodatabase.WithDBName("test"),
	)
	if err := db.Open(); err != nil {
		logger.Fatalf("open db failed, err: %s", err.Error())
		return
	}

	testReady(db)

	em := zeroentity.NewManager(
		zeroentity.NewWrapDB(db),
		zeroentity.NewWrapCache(10*time.Minute),
		zeroentity.NewStat("example", logger),
		logger,
		zerocodecjson.NewJSONCodec(),
	)

	var account Account

	// 查找一个不存在的数据
	_ = em.Get(3, &account)
	_ = em.Get(3, &account)

	if err := em.Get(1, &account); err != nil {
		logger.Errorf("get account failed, err: %s", err.Error())
		return
	}
	if err := em.Get(1, &account); err != nil {
		logger.Errorf("get account failed, err: %s", err.Error())
		return
	}

	account.Age = 18
	if err := em.Update(1, account); err != nil {
		logger.Errorf("update account failed, err: %s", err.Error())
		return
	}

	var account2 Account
	if err := em.Get(1, &account2); err != nil {
		logger.Errorf("get account2 failed, err: %s", err.Error())
		return
	}
	if account2.Age != 18 {
		logger.Errorf("get account2.age failed, age: %d", account2.Age)
		return
	}

	waitSignal()
}

func testReady(database zerodatabase.Database) {
	db := database.DB()
	if err := db.AutoMigrate(&Account{}); err != nil {
		return
	}

	db.Create(&Account{
		UUID:     1,
		Username: "zero",
		Password: "e10adc3949ba59abbe56e057f20f883e",
		Age:      15,
	})
}

// waitSignal 监听信号
func waitSignal() {
	// ctrl + c 或者 kill
	sigs := []os.Signal{syscall.SIGINT, syscall.SIGTERM}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, sigs...)

	sig := <-ch

	signal.Stop(ch)

	log.Println(sig)
}
