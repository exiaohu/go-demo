package config

import (
	"errors"

	"github.com/spf13/viper"
)

// Config 应用程序配置

type Config struct {
	AppName string `json:"app_name" yaml:"app_name"`
	Version string `json:"version"  yaml:"version"`
	Port    int    `json:"port"     yaml:"port"`
	Debug   bool   `json:"debug"    yaml:"debug"`
	// 数据库配置
	Database struct {
		Host         string `json:"host"           yaml:"host"`
		Port         int    `json:"port"           yaml:"port"`
		Name         string `json:"name"           yaml:"name"`
		User         string `json:"user"           yaml:"user"`
		Password     string `json:"password"       yaml:"password"`
		MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"`
		MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"`
		MaxLifeTime  int    `json:"max_life_time"  yaml:"max_life_time"`
	} `json:"database" yaml:"database"`
	// 限流配置
	RateLimit struct {
		Enabled bool    `json:"enabled" yaml:"enabled"`
		RPS     float64 `json:"rps"     yaml:"rps"`
		Burst   int     `json:"burst"   yaml:"burst"`
	} `json:"rate_limit" yaml:"rate_limit"`
}

// C 全局配置实例
var C *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")

	// 设置默认值
	viper.SetDefault("app_name", "playground")
	viper.SetDefault("version", "1.0.0")
	viper.SetDefault("port", 8080)
	viper.SetDefault("debug", false)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "playground.db")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.max_life_time", 3600)

	// 限流默认值
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.rps", 100.0)
	viper.SetDefault("rate_limit.burst", 20)

	// 设置环境变量前缀
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，返回默认配置
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			C = &Config{}
			if unmarshalErr := viper.Unmarshal(C); unmarshalErr != nil {
				return nil, unmarshalErr
			}
			return C, nil
		}
		return nil, err
	}

	// 解析配置到结构体
	C = &Config{}
	if err := viper.Unmarshal(C); err != nil {
		return nil, err
	}

	return C, nil
}
