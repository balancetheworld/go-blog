package repo

import (
	 "fmt"
        "os"
        "path/filepath"

	        "github.com/zyj/my-blog/internal/model"
	        "github.com/zyj/my-blog/pkg/constant"
        "github.com/zyj/my-blog/pkg/utils"
        "gorm.io/driver/postgres"
        "gorm.io/driver/sqlite"
        "gorm.io/gorm"
)

type DBConfig struct {
        Driver           string
        SQLitePath       string
        PostgresHost     string
        PostgresPort     int
        PostgresUser     string
        PostgresPassword string
        PostgresDBName   string
        PostgresSSLMode  string
  }

  //使用指针，不需要拷贝整个结构体
  var db *gorm.DB

  func GetDB() *gorm.DB{
	return db
  }

  func loadDBConfig() DBConfig {
      return DBConfig{
		Driver:           utils.Get(constant.EnvKeyDBDriver, "sqlite"),
                SQLitePath:       utils.Get(constant.EnvKeyDBSQLitePath, "data/blog.db"),
                PostgresHost:     utils.Get(constant.EnvKeyDBHost, "127.0.0.1"),
                PostgresPort:     utils.GetAsInt(constant.EnvKeyDBPort, 5432),
                PostgresUser:     utils.Get(constant.EnvKeyDBUser, "postgres"),
                PostgresPassword: utils.Get(constant.EnvKeyDBPassword),
                PostgresDBName:   utils.Get(constant.EnvKeyDBName, "my_blog"),
                PostgresSSLMode:  utils.Get(constant.EnvKeyDBSSLMode, "disable"),
	  }
  }

  func InitDatabase() error {
	// 读取数据库相关配置（从环境变量/.env加载host、账号、路径等信息）
	config := loadDBConfig()
	// 初始化GORM全局配置对象，用于传入数据库初始化函数
	gormConfig := &gorm.Config{}

	// database：临时变量，保存初始化成功后的gorm数据库实例
	var database  *gorm.DB
	var err error

	// 根据配置中的驱动类型，区分初始化逻辑
	switch config.Driver {
	case "sqlite":
		// 文件数据库，传入sqlite文件路径与gorm配置，创建连接
		database, err = initSQLite(config.SQLitePath, gormConfig)
	case "postgres":
		// 网络型数据库，传入完整数据库配置结构体，创建连接
		database, err = initPostgres(config, gormConfig)
	default:
		// 驱动不支持，直接返回错误终止初始化
		return fmt.Errorf("unsupported database driver: %s", config.Driver)
	}
	// 数据库连接发生错误，向上返回
	if err != nil {
		return err
  }
  db = database
  if err := migrate(); err != nil {
	db = nil
	return fmt.Errorf("migrate database: %w", err)
  }
  return nil
  }

   func initSQLite(path string, config *gorm.Config) (*gorm.DB, error) {
        if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
                return nil, fmt.Errorf("create sqlite directory: %w", err)
        }
        return gorm.Open(sqlite.Open(path), config)
  }

  func initPostgres(config DBConfig, gormConfig *gorm.Config) (*gorm.DB, error) {
        dsn := fmt.Sprintf(
                "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
                config.PostgresHost,
                config.PostgresPort,
                config.PostgresUser,
                config.PostgresPassword,
                config.PostgresDBName,
                config.PostgresSSLMode,
        )
        return gorm.Open(postgres.Open(dsn), gormConfig)
  }

  func migrate() error {
	        return db.AutoMigrate(&model.User{})
  }
