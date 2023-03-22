package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"

	e "github.com/Walker088/gorealestate/error"
)

const (
	currentPackage = "github.com/Walker088/gorealestate/config"

	ConfigFileNotFoundError = "C000001"
	ConfigUnmarshalError    = "C000002"
)

type AppConfig struct {
	pgConfig     *PgConfig
	loggerConfig *LoggerConfig
}

type PgConfig struct {
	DbHost   string `mapstructure:"DB_HOST"`
	DbPort   int    `mapstructure:"DB_PORT"`
	DbSchema string `mapstructure:"DB_SCHEMA"`
	DbName   string `mapstructure:"DB_NAME"`
	DbUser   string `mapstructure:"DB_USER"`
	DbPass   string `mapstructure:"DB_PW"`

	MinConns              int32         `mapstructure:"DB_MIN_CONN"`
	MaxConns              int32         `mapstructure:"DB_MAX_CONN"`
	MaxConnIdleTime       time.Duration `mapstructure:"DB_MAX_CONN_IDLE"`
	MaxConnLifetime       time.Duration `mapstructure:"DB_MAX_CONN_LIFETIME"`
	MaxConnLifetimeJitter time.Duration `mapstructure:"DB_MAX_CONN_LIFETIME_JITTER"`
	HealthCheckPeriod     time.Duration `mapstructure:"DB_HEALTH_CHECK_PERIOD"`
}

type LoggerConfig struct {
	ConsoleLogLevel string `mapstructure:"CONSOLE_LOG_LEVEL"`
	FileLogLevel    string `mapstructure:"FILE_LOG_LEVEL"`
	EncoderConfig   zapcore.EncoderConfig
}

func New(envFile string) (*AppConfig, *e.ErrorData) {
	encoder := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(
				fmt.Sprintf(
					"%d-%02d-%02d %02d:%02d:%02d",
					t.Year(),
					t.Month(),
					t.Day(),
					t.Hour(),
					t.Minute(),
					t.Second(),
				),
			)
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	pgconf := &PgConfig{
		DbHost:                "localhost",
		DbPort:                5432,
		DbSchema:              "public",
		DbName:                "gorealestate",
		DbUser:                "postgres",
		MinConns:              10,
		MaxConns:              100,
		MaxConnIdleTime:       10 * time.Minute,
		MaxConnLifetime:       30 * time.Minute,
		MaxConnLifetimeJitter: 1 * time.Minute,
		HealthCheckPeriod:     1 * time.Minute,
	}
	zapconf := &LoggerConfig{
		ConsoleLogLevel: "info",
		FileLogLevel:    "info",
	}
	viper.SetConfigFile(envFile)
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, e.NewErrorData(
				ConfigFileNotFoundError,
				fmt.Sprintf("failed to load file, %s not found", envFile),
				fmt.Sprintf("%s.New", currentPackage),
				nil,
				nil,
			)
		}
	}

	viper.AutomaticEnv()
	if err := viper.Unmarshal(zapconf); err != nil {
		return nil, e.NewErrorData(
			ConfigUnmarshalError,
			err.Error(),
			fmt.Sprintf("%s.New", currentPackage),
			nil,
			nil,
		)
	}
	if err := viper.Unmarshal(pgconf); err != nil {
		return nil, e.NewErrorData(
			ConfigUnmarshalError,
			err.Error(),
			fmt.Sprintf("%s.New", currentPackage),
			nil,
			nil,
		)
	}
	c := &AppConfig{
		pgConfig:     pgconf,
		loggerConfig: zapconf,
	}
	c.loggerConfig.EncoderConfig = encoder
	return c, nil
}

func (c *AppConfig) WatchConfig() {
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()
}

func (c *AppConfig) GetPgConfig() *PgConfig {
	return c.pgConfig
}

func (c *AppConfig) GetLoggerConfig() *LoggerConfig {
	return c.loggerConfig
}

func (p *PgConfig) ToConnString() string {
	//urlExample := postgres://username:password@localhost:5432/database_name
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		p.DbUser,
		p.DbPass,
		p.DbHost,
		p.DbPort,
		p.DbName,
	)
}

func (l *LoggerConfig) GetConsoleLogLvl() zapcore.Level {
	if l.ConsoleLogLevel == "debug" {
		return -1
	}
	if l.ConsoleLogLevel == "info" {
		return 0
	}
	if l.ConsoleLogLevel == "warn" {
		return 1
	}
	if l.ConsoleLogLevel == "error" {
		return 2
	}
	if l.ConsoleLogLevel == "dpanic" {
		return 3
	}
	if l.ConsoleLogLevel == "panic" {
		return 4
	}
	if l.ConsoleLogLevel == "fatal" {
		return 5
	}
	return 0
}

func (l *LoggerConfig) GetFileLogLvl() zapcore.Level {
	if l.FileLogLevel == "debug" {
		return -1
	}
	if l.FileLogLevel == "info" {
		return 0
	}
	if l.FileLogLevel == "warn" {
		return 1
	}
	if l.FileLogLevel == "error" {
		return 2
	}
	if l.FileLogLevel == "dpanic" {
		return 3
	}
	if l.FileLogLevel == "panic" {
		return 4
	}
	if l.FileLogLevel == "fatal" {
		return 5
	}
	return 0
}
