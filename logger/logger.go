package logger

import (
	"os"
	"sync"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

var once sync.Once

func init() {
	once.Do(func() {
		// -------------------------
		// File Encoder
		// -------------------------

		fileConfig := zap.NewProductionEncoderConfig()
		fileConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		fileEncoder := zapcore.NewJSONEncoder(fileConfig)

		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "./log/logger.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     28,
		})

		// -------------------------
		// Console Encoder
		// -------------------------

		consoleConfig := zap.NewDevelopmentEncoderConfig()
		consoleConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		consoleConfig.EncodeLevel = coloredLevelEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)

		consoleWriter := zapcore.AddSync(os.Stdout)

		// -------------------------
		// Core
		// -------------------------

		core := zapcore.NewTee(
			// File
			zapcore.NewCore(
				fileEncoder,
				fileWriter,
				zapcore.InfoLevel,
			),

			// Console
			zapcore.NewCore(
				consoleEncoder,
				consoleWriter,
				zapcore.DebugLevel,
			),
		)

		Logger = zap.New(
			core,
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	})
}

func coloredLevelEncoder(
	level zapcore.Level,
	enc zapcore.PrimitiveArrayEncoder,
) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString(color.MagentaString(level.CapitalString()))

	case zapcore.InfoLevel:
		enc.AppendString(color.GreenString(level.CapitalString()))

	case zapcore.WarnLevel:
		enc.AppendString(color.YellowString(level.CapitalString()))

	case zapcore.ErrorLevel:
		enc.AppendString(color.RedString(level.CapitalString()))

	case zapcore.DPanicLevel:
		enc.AppendString(color.RedString(level.CapitalString()))

	case zapcore.PanicLevel:
		enc.AppendString(color.RedString(level.CapitalString()))

	case zapcore.FatalLevel:
		enc.AppendString(color.RedString(level.CapitalString()))

	default:
		enc.AppendString(level.CapitalString())
	}
}
