package zapwrapper

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

var callerSkip int = 1

func SetCallerSkip(v int) {
	callerSkip = v
}

func InitZapLogger(showConsole bool, serviceName string, logFilePath string) (*zap.Logger, error) {

	zCfg := zap.NewProductionConfig()
	zCfg.EncoderConfig.TimeKey = "t"
	zCfg.EncoderConfig.LevelKey = "l"
	zCfg.EncoderConfig.NameKey = "logger"
	zCfg.EncoderConfig.CallerKey = "c"
	zCfg.EncoderConfig.MessageKey = "msg"
	zCfg.EncoderConfig.StacktraceKey = "st"
	zCfg.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
	zCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	zCfg.EncoderConfig.EncodeTime = timeEncoder
	zCfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	zCfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	hook := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100, // megabytes
		MaxBackups: 10,
		MaxAge:     10, // days
		LocalTime:  true,
		Compress:   true,
	}

	fileWriter := zapcore.AddSync(hook)

	jsonEncoder := zapcore.NewJSONEncoder(zCfg.EncoderConfig)
	if len(serviceName) > 0 {
		jsonEncoder.AddString("svr", serviceName)
	}

	var core zapcore.Core
	if showConsole {
		consoleDebugging := zapcore.Lock(os.Stdout)
		consoleEncoder := zapcore.NewConsoleEncoder(zCfg.EncoderConfig)
		if len(serviceName) > 0 {
			consoleEncoder.AddString("svr", serviceName)
		}

		core = zapcore.NewTee(
			zapcore.NewCore(jsonEncoder, fileWriter, zap.DebugLevel),
			zapcore.NewCore(consoleEncoder, consoleDebugging, zap.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(jsonEncoder, fileWriter, zap.DebugLevel)
	}

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(callerSkip), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(logger)

	return logger, nil
}
