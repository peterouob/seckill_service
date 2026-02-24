package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logg *zap.Logger

func InitLogger(serviceName ...string) {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(cfg)

	name := "app"
	if len(serviceName) > 0 && serviceName[0] != "" {
		name = serviceName[0]
	}

	logFile := getLogFile(name)
	writer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logFile))

	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	logg = zap.New(core, zap.AddCaller())
}

func getLogFile(serviceName string) *os.File {
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		if _, err := os.Stat("../../../go.mod"); err == nil {
			logDir = "../../../app_logs"
		} else {
			logDir = "app_logs"
		}
	}

	if err := os.MkdirAll(logDir, 0777); err != nil {
		panic(fmt.Sprintf("failed to create log directory: %v", err))
	}

	today := time.Now().Format("2006-01-02")
	logFileName := filepath.Join(logDir, fmt.Sprintf("log-%s-%s.log", serviceName, today))
	file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file: %v", err))
	}
	return file
}

func getLogger() *zap.Logger {
	if logg == nil {
		InitLogger()
	}
	return logg
}

func Log(msg any) {
	getLogger().Info(fmt.Sprint(msg))
}

func Logf(format string, args ...any) {
	getLogger().Info(fmt.Sprintf(format, args...))
}

func Warn(msg any) {
	getLogger().Warn(msg.(string))
}

func Error(msg string, err error, fields ...zap.Field) {
	getLogger().Error(msg, append(fields, zap.Error(err))...)
}

func HandelError(msg string, err error, f ...func(args ...any)) {
	if err != nil {
		if len(f) > 0 && f[0] != nil {
			f[0](msg, err)
		}

		Error(msg, err)
	}
}

func ErrorMsg(msg string) {
	getLogger().Error(msg)
}

func ErrorMsgF(format string, args ...any) {
	getLogger().Error(fmt.Sprintf(format, args...))
}
