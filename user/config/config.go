package config

import (
	"path"

	"common/applog"
	"common/env"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server ServerConfig     `toml:"server"`
	Redis  RedisConfig      `toml:"redis"`
	Mysql  MysqlConfig      `toml:"mysql"`
	Mongo  MongoConfig      `toml:"mongo"`
	AppLog applog.LogConfig `toml:"app_log"`
	Jaeger JaegerConfig     `toml:"jaeger"`
	Grpc   GrpcConfig       `toml:"grpc"`
	Etcd   EtcdConfig       `toml:"etcd"`
}

// ServerConfig server配置
type ServerConfig struct {
	Name string `toml:"name"`
	Env  string `toml:"env"`
}

// RedisConfig redis配置（包含client和pool子表）
type RedisConfig struct {
	Client RedisClientConfig `toml:"client"`
	Pool   RedisPoolConfig   `toml:"pool"`
}

// RedisClientConfig redis客户端配置
type RedisClientConfig struct {
	Name           string `toml:"name"`
	Host           string `toml:"host"`
	Port           int    `toml:"port"`
	DB             int    `toml:"db"`
	Password       string `toml:"password"`
	KeepAlive      int    `toml:"keep_alive"`
	ConnectTimeout int    `toml:"connect_timeout"`
	WriteTimeout   int    `toml:"write_timeout"`
	ReadTimeout    int    `toml:"read_timeout"`
}

// RedisPoolConfig redis连接池配置
type RedisPoolConfig struct {
	MaxIdle         int `toml:"max_idle"`
	MaxActive       int `toml:"max_active"`
	IdleTimeout     int `toml:"idle_timeout"`
	MaxConnLifetime int `toml:"max_conn_lifetime"`
	PoolSize        int `toml:"pool_size"`
}

// MysqlConfig mysql配置（包含master和slaves子表）
type MysqlConfig struct {
	Name            string       `toml:"name"`
	Separation      bool         `toml:"separation"`
	ConnTimeOut     int          `toml:"conn_time_out"`
	WriteTimeOut    int          `toml:"write_time_out"`
	ReadTimeOut     int          `toml:"read_time_out"`
	Master          MysqlMaster  `toml:"master"`
	Slaves          []MysqlSlave `toml:"slaves"` // 数组对应TOML的[[mysql.slaves]]
	DBDriver        string       `toml:"db_driver"`
	ConnMaxLifeTime int          `toml:"conn_max_life_time"`
	MaxIdleCount    int          `toml:"max_idle_count"`
	MaxOpenCount    int          `toml:"max_open_count"`
	IsOutLog        bool         `toml:"is_out_log"`
	SQLLogLen       int          `toml:"sql_log_len"`
	SQLArgsLogLen   int          `toml:"sql_args_log_len"`
	LogIDTransport  bool         `toml:"log_id_transport"`
	DSNParams       string       `toml:"dsn_params"`
}

// MysqlMaster mysql主库配置
type MysqlMaster struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
}

// MysqlSlave mysql从库配置（数组元素）
type MysqlSlave struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
}

// MongoConfig mongo配置
type MongoConfig struct {
	Name            string `toml:"name"`
	Host            string `toml:"host"`
	Port            int    `toml:"port"`
	Username        string `toml:"username"`
	Password        string `toml:"password"`
	Database        string `toml:"database"`
	AuthSource      string `toml:"auth_source"`
	ConnectTimeout  int    `toml:"connect_timeout"`
	MaxPoolSize     uint64 `toml:"max_pool_size"`
	MinPoolSize     uint64 `toml:"min_pool_size"`
	MaxConnIdleTime int    `toml:"max_conn_idle_time"`
}

// JaegerConfig jaeger追踪配置（修正原拼写错误）
type JaegerConfig struct {
	Endpoints              string `toml:"endpoints"`
	IsOpenOnlySamplerError bool   `toml:"is_open_only_sampler_error"`
}

// GrpcConfig grpc配置
type GrpcConfig struct {
	EtcdAddr string `toml:"etcd_addr"`
	Addr     string `toml:"addr"`
	Name     string `toml:"name"`
	Version  string `toml:"version"`
	Weight   int    `toml:"weight"`
}

// EtcdConfig etcd配置
type EtcdConfig struct {
	Addrs []string `toml:"addrs"`
}

var cfg Config

// InitConfig 加载TOML配置文件
func InitConfig(confFileName string) error {
	filePath := path.Join(env.GetEnvConfig().ConfDir, confFileName)
	if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
		return err
	}
	return nil
}

func GetConfig() Config {
	return cfg
}
