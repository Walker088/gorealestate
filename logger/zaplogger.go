package logger

import (
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/Walker088/gorealestate/config"
)

const (
	fileOutputPath = "./logs"
	logFileName    = "gorealestate.log"
	errFileName    = "gorealestate_err.log"
)

type WriteSyncer struct {
	io.Writer
}

func (ws WriteSyncer) Sync() error {
	return nil
}

func getWriteSyncer(filename string) zapcore.WriteSyncer {
	var ioWriter = &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10, // MB
		MaxBackups: 3,  // number of backups
		MaxAge:     14, //days
		LocalTime:  true,
		Compress:   false, // disabled by default
	}
	var sw = WriteSyncer{
		ioWriter,
	}
	return sw
}

func New(cfg *config.LoggerConfig) *zap.SugaredLogger {
	if _, err := os.Stat(fileOutputPath); os.IsNotExist(err) {
		os.MkdirAll(fileOutputPath, 0700)
	}

	var logger *zap.Logger

	fn := fmt.Sprintf("%s/%s", fileOutputPath, logFileName)
	consoleEnc := zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	fileEnc := zapcore.NewJSONEncoder(cfg.EncoderConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEnc, zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(cfg.GetConsoleLogLvl())),
		zapcore.NewCore(fileEnc, zapcore.AddSync(getWriteSyncer(fn)), zap.NewAtomicLevelAt(cfg.GetFileLogLvl())),
	)

	logger = zap.New(core)
	return logger.Sugar()
}
