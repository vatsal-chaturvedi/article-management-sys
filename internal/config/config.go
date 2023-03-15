package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"
	"io/ioutil"
	"log"
	"reflect"
	"time"
)

type Config struct {
	ServerConfig ServerConfig `json:"server_config"`
	DataBase     DbCfg        `json:"data_source"`
	Cacher       CacheConfig  `json:"cacher"`
}

type SvcConfig struct {
	Cfg       *Config
	SvrCfg    ServerConfig
	DbSvc     DbSvc
	CacherSvc CacheSvc
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// DbSvc struct defines the database service
type DbSvc struct {
	Db *sql.DB
}
type CacheSvc struct {
	Rdb *redis.Client
}

// DbCfg struct defines the configuration for the database service
type DbCfg struct {
	Port      string `json:"dbPort"`
	Host      string `json:"dbHost"`
	Driver    string `json:"dbDriver"`
	User      string `json:"dbUser"`
	Pass      string `json:"dbPass"`
	DbName    string `json:"dbName"`
	TableName string `json:"tableName"`
}

type CacheConfig struct {
	Address           string        `json:"address"`
	Username          string        `json:"username"`
	Password          string        `json:"password"`
	DB                int           `json:"db"`
	MaxRetries        int           `json:"max_retries"`
	DialTimeout       time.Duration `json:"dial_timeout"`
	PoolSize          int           `json:"pool_size"`
	MinIdleConns      int           `json:"min_idle_conns"`
	KeyExpiry         string        `json:"key_expiry"`
	KeyExpiryDuration time.Duration
}

func LoadFromJson(filepath string, cfg interface{}) error {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(cfg)
	if v.Type().Kind() != reflect.Ptr || !v.Elem().CanSet() {
		return fmt.Errorf("unable to set into given type: must be a pointer")
	}
	err = json.Unmarshal(content, cfg)
	if err != nil {
		return err
	}
	return nil
}
func ConnectSql(cfg DbCfg, tableName string) *sql.DB {
	// Construct the database connection string from the configuration data.
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True", cfg.User, cfg.Pass, cfg.Host, cfg.Port)

	// Open the database connection.
	db, err := sql.Open(cfg.Driver, connectionString)
	if err != nil {
		panic(err.Error())
	}

	// Create the database schema if it does not already exist.
	dbString := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s ;", cfg.DbName)
	prepare, err := db.Prepare(dbString)
	if err != nil {
		panic(err.Error())
	}
	_, err = prepare.Exec()
	if err != nil {
		panic(err.Error())
	}

	// Close the database connection.
	db.Close()

	// Reopen the database connection with the specified database name.
	connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DbName)
	db, err = sql.Open(cfg.Driver, connectionString)
	if err != nil {
		panic(err.Error())
	}

	// Create the specified table if it does not already exist.
	x := fmt.Sprintf("create table if not exists %s", tableName)
	log.Print(x + model.Schema)
	_, err = db.Exec(x + model.Schema)
	if err != nil {
		panic(err.Error())
	}

	// Return the database connection object.
	return db
}
func ConnectRedis(c CacheConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         c.Address,
		Username:     c.Username,
		Password:     c.Password, // no password set
		DB:           c.DB,       // use default DB
		MaxRetries:   c.MaxRetries,
		DialTimeout:  c.DialTimeout,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
	})
}

func InitSvcConfig(cfg Config) *SvcConfig {
	dataBase := ConnectSql(cfg.DataBase, cfg.DataBase.TableName)
	cache := ConnectRedis(cfg.Cacher)
	duration, err := time.ParseDuration(cfg.Cacher.KeyExpiry)
	if err != nil {
		panic(err.Error())
	}
	cfg.Cacher.KeyExpiryDuration = duration
	return &SvcConfig{
		Cfg:       &cfg,
		SvrCfg:    cfg.ServerConfig,
		CacherSvc: CacheSvc{Rdb: cache},
		DbSvc:     DbSvc{Db: dataBase},
	}
}
