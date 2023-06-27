package main

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

func newDefaultProductionLogEncoder(colorize bool) zapcore.Encoder {
	encCfg := zap.NewProductionEncoderConfig()
	// if interactive terminal, make output more human-readable by default
	encCfg.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format("2006/01/02 15:04:05.000"))
	}
	if colorize {
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	return zapcore.NewConsoleEncoder(encCfg)
}

func setupLogging() {
	consoleStream := zapcore.Lock(os.Stderr)
	consoleEncoder := newDefaultProductionLogEncoder(true)
	noPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})
	consoleCore := zapcore.NewCore(consoleEncoder, consoleStream, noPriority)

	logger := zap.New(consoleCore)
	defer logger.Sync() // flushes buffer, if any
	log = logger.Sugar().Named("InnkeeperBot")
}
