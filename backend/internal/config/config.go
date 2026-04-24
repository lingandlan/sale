package config

import (
	"fmt"
	"log"
	"os"
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
	CORS     CORSConfig     `mapstructure:"cors"`
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
	Secret             string `mapstructure:"secret"`
	ExpireHours        int    `mapstructure:"expire_hours"`
	RefreshExpireHours int    `mapstructure:"refresh_expire_hours"`
}

// LogConfig 日志配置
type LogConfig struct {
	Mode        string `mapstructure:"mode"`
	Level       string `mapstructure:"level"`
	ServiceName string `mapstructure:"service_name"`
}

// MallConfig WSY商城接口配置
type MallConfig struct {
	BaseURL    string `mapstructure:"base_url"`
	AppID      string `mapstructure:"app_id"`
	AppSecret  string `mapstructure:"app_secret"`
	CustomerID string `mapstructure:"customer_id"`
}

// CORSConfig CORS 跨域配置
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

// envOverrides 环境变量覆盖映射（显式环境变量名 → viper key）
var envOverrides = map[string]string{
	"JWT_SECRET":          "jwt.secret",
	"DB_PASSWORD":         "database.password",
	"REDIS_PASSWORD":      "redis.password",
	"MALL_APP_ID":         "mall.app_id",
	"MALL_APP_SECRET":     "mall.app_secret",
	"MALL_CUSTOMER_ID":    "mall.customer_id",
	"MALL_BASE_URL":       "mall.base_url",
	"CORS_ALLOWED_ORIGINS": "cors.allowed_origins",
}

// placeholderValues 需要警告的占位符值
var placeholderValues = map[string]string{
	"jwt.secret": "your-super-secret-key-change-in-production",
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// 支持环境变量覆盖（APP_ 前缀）
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// 显式环境变量覆盖（优先级最高）
	for envKey, viperKey := range envOverrides {
		if val := os.Getenv(envKey); val != "" {
			viper.Set(viperKey, val)
		}
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// CORS_ALLOWED_ORIGINS: 逗号分隔字符串 → []string
	if raw := os.Getenv("CORS_ALLOWED_ORIGINS"); raw != "" {
		origins := strings.Split(raw, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
		cfg.CORS.AllowedOrigins = origins
	}

	// 检查占位符值
	for viperKey, placeholder := range placeholderValues {
		if viper.GetString(viperKey) == placeholder {
			log.Printf("[WARN] 配置 %s 使用的是默认占位符值，请在生产环境中通过环境变量 %s 设置真实值",
				viperKey, getKeyByValue(envOverrides, viperKey))
		}
	}

	return cfg, nil
}

func getKeyByValue(m map[string]string, value string) string {
	for k, v := range m {
		if v == value {
			return k
		}
	}
	return ""
}
