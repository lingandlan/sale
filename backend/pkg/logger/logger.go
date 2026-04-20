package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config 日志配置
type Config struct {
	Mode        string `mapstructure:"mode"`
	Level       string `mapstructure:"level"`
	ServiceName string `mapstructure:"service_name"`
}

var log *zap.Logger

// Init 初始化日志
func Init(cfg *Config) error {
	var zapCfg zap.Config

	if cfg.Mode == "production" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// 设置日志级别
	level := zapcore.InfoLevel
	if cfg.Level != "" {
		var err error
		level, err = zapcore.ParseLevel(cfg.Level)
		if err != nil {
			level = zapcore.InfoLevel
		}
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	// 输出到 stdout
	zapCfg.OutputPaths = []string{"stdout"}
	zapCfg.ErrorOutputPaths = []string{"stderr"}

	// 添加文件轮转
	fileWriter := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100, // MB
		MaxBackups: 24,
		MaxAge:     7, // 天
		Compress:   true,
		LocalTime:  true,
	}

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(fileWriter),
			level,
		),
	)

	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return nil
}

// GetLogger 获取 logger 实例
func GetLogger() *zap.Logger {
	if log == nil {
		return zap.NewNop()
	}
	return log
}

// Info 记录 Info 日志
func Info(msg string, fields ...zap.Field) {
	if log != nil {
		log.Info(msg, fields...)
	}
}

// Error 记录 Error 日志
func Error(msg string, fields ...zap.Field) {
	if log != nil {
		log.Error(msg, fields...)
	}
}

// Warn 记录 Warn 日志
func Warn(msg string, fields ...zap.Field) {
	if log != nil {
		log.Warn(msg, fields...)
	}
}

// Debug 记录 Debug 日志
func Debug(msg string, fields ...zap.Field) {
	if log != nil {
		log.Debug(msg, fields...)
	}
}

// Fatal 记录 Fatal 日志
func Fatal(msg string, fields ...zap.Field) {
	if log != nil {
		log.Fatal(msg, fields...)
	}
}
