package config

import (
	"common/env"
	"github.com/BurntSushi/toml"
	"path"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Name string `toml:"name"`
	Env  string `toml:"env"`
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

// JaegerConfig Jaeger 配置
type JaegerConfig struct {
	Endpoints              string `toml:"endpoints"`
	IsOpenOnlySamplerError bool   `toml:"is_open_only_sampler_error"`
}

// EtcdConfig Etcd 配置
type EtcdConfig struct {
	Addrs []string `toml:"addrs"`
}

// Config 总配置
type Config struct {
	Server ServerConfig `toml:"server"`
	Jaeger JaegerConfig `toml:"jaeger"`
	Etcd   EtcdConfig   `toml:"etcd"`
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
