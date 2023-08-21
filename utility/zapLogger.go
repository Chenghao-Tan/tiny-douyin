package utility

import (
	"douyin/conf"

	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _logger *zap.SugaredLogger

func Logger() *zap.SugaredLogger {
	return _logger
}

func InitLogger() {
	core := zapcore.NewCore(getEncoder(), // 输出格式
		zapcore.NewMultiWriteSyncer(getWriterSyncer(), zapcore.AddSync(os.Stdout)), // 输出位置+同时输出到控制台
		getLevelPriority()) // 输出级别

	logger := zap.New(core).Sugar()

	if conf.Cfg().Log.ShowLine {
		logger = logger.WithOptions(zap.AddCaller()) // 显示[调用的文件 函数名称:行号] 如[main.go main:42]
	}

	_logger = logger
}

// 输出格式
func getEncoder() zapcore.Encoder {
	// 参考 zap.NewProductionEncoderConfig()
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 输出时间格式
func customTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(conf.Cfg().Log.Prefix + t.Format("2006/01/02 15:04:05"))
}

// 输出位置
func getWriterSyncer() zapcore.WriteSyncer {
	// 确保输出路径存在
	err := os.MkdirAll(filepath.Join(conf.Cfg().Log.Path, ""), 0755)
	if err != nil {
		panic(err)
	}

	// 日志路径
	stLogFilePath := filepath.Join(conf.Cfg().Log.Path, (time.Now().Format("2006-01-02") + ".log")) // 拼接路径
	stLogFilePath = filepath.Clean(stLogFilePath)                                                   // 清理路径语法
	stLogFilePath, err = filepath.Abs(stLogFilePath)                                                // 转换为绝对路径
	if err != nil {
		panic(err)
	}

	// 日志分割
	hook := lumberjack.Logger{
		Filename:   stLogFilePath,
		MaxSize:    conf.Cfg().Log.MaxSize,
		MaxBackups: conf.Cfg().Log.MaxBackups,
		MaxAge:     conf.Cfg().Log.MaxAge,
		Compress:   conf.Cfg().Log.Compress,
	}
	return zapcore.AddSync(&hook)
}

// 输出级别
func getLevelPriority() zapcore.LevelEnabler {
	switch strings.ToLower(conf.Cfg().Log.Level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	}
	return zap.InfoLevel
}
