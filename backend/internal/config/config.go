package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 全局配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Mall     MallConfig     `mapstructure:"mall"`
}

// ServerConfig HTTP 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug / release
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	TLS             bool   `mapstructure:"tls"`
	CAcert          string `mapstructure:"ca_cert"`
}

// DSN 返回 MySQL DSN
func (d *DatabaseConfig) DSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name)
	if d.TLS {
		dsn += "&tls=true"
	}
	return dsn
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret               string `mapstructure:"secret"`
	ExpireHours          int    `mapstructure:"expire_hours"`
	RefreshExpireHours   int    `mapstructure:"refresh_expire_hours"`
}

// LogConfig 日志配置
type LogConfig struct {
	Mode         string `mapstructure:"mode"`
	Level        string `mapstructure:"level"`
	ServiceName  string `mapstructure:"service_name"`
}

// MallConfig WSY商城接口配置
type MallConfig struct {
	BaseURL    string `mapstructure:"base_url"`
	AppID      string `mapstructure:"app_id"`
	AppSecret  string `mapstructure:"app_secret"`
	CustomerID string `mapstructure:"customer_id"`
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// 支持环境变量覆盖
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
