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
	return (&Sharding{}).Register(configs...)
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

	if err := s.checkUniqueTable(); err != nil {
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

		config.tableNames = make([]string, 0, len(config.Tables))

		for _, table := range config.Tables {
			if tableName, ok := table.(string); ok {
				config.tableNames = append(config.tableNames, tableName)
			} else {
				stmt := &gorm.Statement{DB: config.Databases[0].DB()}
				if err := stmt.Parse(table); err == nil {
					config.tableNames = append(config.tableNames, stmt.Table)
				} else {
					return err
				}
			}
		}
	}

	return nil
}

// checkUniqueTable 一张表不得重复配置
func (s *Sharding) checkUniqueTable() error {

	uniqueTableNames := make(map[string]struct{})
	for _, config := range s.configs {
		for _, tableName := range config.tableNames {
			if _, exist := uniqueTableNames[tableName]; exist {
				return ErrTableDuplicates
			}

			uniqueTableNames[tableName] = struct{}{}
		}
	}

	return nil
}

func (s *Sharding) buildTableGroups() error {
	s.groups = map[string]*TableGroup{}

	for _, config := range s.configs {
		n := config.NumberOfShards
		if n <= 0 {
			n = 1
		} else if n > 9999 {
			return errors.New("number of shards should be less than or equal to 9999")
		}

		tableNameFormat := "%s_%04d"
		if n < 10 {
			tableNameFormat = "%s_%01d"
		} else if n < 100 {
			tableNameFormat = "%s_%02d"
		} else if n < 1000 {
			tableNameFormat = "%s_%03d"
		}

		// accounts, players
		for _, tableName := range config.tableNames {
			group, err := buildTableGroup(tableName, tableNameFormat, n, config)
			if err != nil {
				return err
			}

			s.groups[tableName] = group
		}
	}

	return nil
}

func buildTableGroup(tableName, tableNameFormat string, tableNum int, config *Config) (*TableGroup, error) {
	// 拼接表名
	tableShards := make([]string, 0, tableNum)
	for i := 0; i < tableNum; i++ {
		tableName := fmt.Sprintf(tableNameFormat, tableName, i)
		tableShards = append(tableShards, tableName)
	}

	// TODO 在对应的数据库创建表
	// for _, db := range config.Databases {
	// 	gDB := db.DB()
	// 	st := &gorm.Statement{DB: gDB}
	// 	for _, tableName := range tableShards {
	// 		tx := st.DB.Session(&gorm.Session{}).Table(tableName)
	// 		if err := gDB.Dialector.Migrator(tx).AutoMigrate()
	// 	}
	// }

	return nil, nil
}

// Register 注册当前插件
func (s *Sharding) Register(configs ...*Config) *Sharding {
	s.configs = append(s.configs, configs...)
	return nil
}

// Scope 提供表名，并根据传入的 id，计算出对应的数据库和表名
func (s *Sharding) Scope(table interface{}, id uint64) (zerodatabase.Database, string) {
	// TODO 提供表名，并根据传入的 id，计算出对应的数据库和表名
	return nil, ""
}

// TableGroup 表配置
// tables 按照 policyFn 策略分布在 databases 中
// 每一个配置对应一个 TableGroup
type TableGroup struct {
	// table 表名，例如 accounts
	table string
	// tableShards 分表名，例如 accounts_01, accounts_02
	tableShards []string
	databases   []zerodatabase.Database
	policyFn    func(id uint64) int
}

// Config 分库分表配置
type Config struct {
	// NumberOfShards 分表数量
	NumberOfShards int

	// Tables 哪些表使用同一个分表方式
	// 可以是表名，例如 accounts
	// 可以是对象，例如 Account{}
	Tables []any

	// tableName 解析出的表名
	tableNames []string

	// Databases 分布到哪些数据库
	Databases []zerodatabase.Database

	// PolicyFn 自定义策略
	PolicyFn func(id uint64) int
}
