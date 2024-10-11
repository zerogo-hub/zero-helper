// Package 分库分表 (未完成)
package sharding

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	zerodatabase "github.com/zerogo-hub/zero-helper/database"
)

var (
	// ErrNoTables 未配置数据表
	ErrNoTables = errors.New("no tables provided in sharding config")
	// ErrNoDatabases 未配置数据库
	ErrNoDatabases = errors.New("no databases provided in sharding config")
	// ErrTableDuplicates 数据表重复配置，每张表只能配置一个规则
	ErrTableDuplicates = errors.New("table duplicates provided in shrolling")
)

// Register 注册 gorm 插件
func Register(configs ...*Config) *Sharding {
	s := &Sharding{
		groups: make(map[string]*TableGroup),
	}
	return s.Register(configs...)
}

// Sharding 分库分表
type Sharding struct {
	// groups 每一张表的分组信息
	groups  map[string]*TableGroup
	configs []*Config
}

// Name gorm 插件接口
func (s *Sharding) Name() string {
	return "zerodatabase:sharding"
}

// Initialize gorm 插件接口
func (s *Sharding) Initialize(_ *gorm.DB) error {
	if err := s.parseTables(); err != nil {
		return err
	}

	return s.buildTableGroups()
}

// parsedTables 解析表名
func (s *Sharding) parseTables() error {
	for _, config := range s.configs {
		if len(config.Tables) == 0 {
			return ErrNoTables
		}

		if len(config.Databases) == 0 {
			return ErrNoDatabases
		}

		var err error
		for _, table := range config.Tables {
			group := &TableGroup{
				TableModle:  table,
				tableShards: make([]string, 0, config.NumberOfShards),
				config:      config,
			}

			group.tableName, err = s.parseTableName(table)
			if err != nil {
				return err
			}
			s.groups[group.tableName] = group
		}
	}

	return nil
}

func (s *Sharding) parseTableName(table any) (string, error) {
	if tableName, ok := table.(string); ok {
		return tableName, nil
	}

	stmt := &gorm.Statement{DB: s.configs[0].Databases[0].DB()}
	if err := stmt.Parse(table); err != nil {
		return "", err
	}

	return stmt.Table, nil
}

func (s *Sharding) buildTableGroups() error {

	for _, group := range s.groups {
		if err := group.build(); err != nil {
			return err
		}
	}

	return nil
}

// Register 注册当前插件
func (s *Sharding) Register(configs ...*Config) *Sharding {
	s.configs = append(s.configs, configs...)
	return s
}

// Scope 提供表名，并根据传入的 id，计算出对应的数据库和表名
func (s *Sharding) Scope(table interface{}, id uint64) (zerodatabase.Database, string, error) {
	tableName, err := s.parseTableName(table)
	if err != nil {
		return nil, "", err
	}

	group, ok := s.groups[tableName]
	if !ok {
		return nil, "", fmt.Errorf("table %s not found", tableName)
	}

	var idxDB, idxTable uint64
	if group.config.PolicyFn != nil {
		idxDB, idxTable = group.config.PolicyFn(id)
	} else {
		idxDB = id
		idxTable = id
	}
	if idxDB > uint64(len(group.config.Databases)) {
		idxDB = id % uint64(len(group.config.Databases))
	}
	if idxTable > uint64(len(group.tableShards)) {
		idxTable = id % uint64(len(group.tableShards))
	}

	return group.config.Databases[idxDB], group.tableShards[idxTable], nil
}

// TableGroup 表配置
// tables 按照 policyFn 策略分布在 databases 中
// 每一个表模型对应一个 TableGroup
type TableGroup struct {
	// TableModle 表对应的模型
	TableModle any
	// tableName 表名，例如 accounts
	tableName string
	// tableShards 分表名，例如 accounts_01, accounts_02
	tableShards []string
	config      *Config
}

func (tg *TableGroup) tableNameFormat() string {
	n := tg.config.NumberOfShards
	if n <= 0 {
		n = 1
	} else if n > 9999 {
		n = 9999
	}

	tableNameFormat := "%s_%04d"
	if n < 10 {
		tableNameFormat = "%s_%01d"
	} else if n < 100 {
		tableNameFormat = "%s_%02d"
	} else if n < 1000 {
		tableNameFormat = "%s_%03d"
	}

	return tableNameFormat
}

func (tg *TableGroup) build() error {
	nameFormat := tg.tableNameFormat()

	// 拼接表名
	tableNum := tg.config.NumberOfShards
	for i := 0; i < tableNum; i++ {
		tableName := fmt.Sprintf(nameFormat, tg.tableName, i)
		tg.tableShards = append(tg.tableShards, tableName)
	}

	// 在对应的数据库创建表
	for idx, tableUniqueName := range tg.tableShards {
		db := tg.config.Databases[idx%len(tg.config.Databases)]
		DB := db.DB()
		st := &gorm.Statement{DB: DB}

		tx := st.DB.Session(&gorm.Session{}).Table(tableUniqueName)
		if err := DB.Dialector.Migrator(tx).AutoMigrate(tg.TableModle); err != nil {
			return err
		}
	}

	return nil
}

// Config 分库分表配置
type Config struct {
	// NumberOfShards 分表数量
	NumberOfShards int

	// Tables 哪些表使用同一个分表方式
	// 可以是对象，例如 Account{}
	Tables []any

	// Databases 分布到哪些数据库
	Databases []zerodatabase.Database

	// PolicyFn 自定义策略
	// return 数据库索引，表索引
	PolicyFn func(id uint64) (uint64, uint64)
}
