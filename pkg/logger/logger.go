package logger

import (
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger
var filepath = "./logs"
var defaultLevel = zapcore.InfoLevel

//InitZapLogger 初始化zap日志
func InitZapLogger(logdir string, zapLevel zapcore.Level) {
	filepath = logdir
	//默认等级 >= 的输出
	defaultLevel = zapLevel
	// 设置一些基本日志格式
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		NameKey:       "logger",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		TimeKey:       "timestamp",
		FunctionKey:   "function",
		StacktraceKey: "stack",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder, //短路径
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	//配置哪些level 写log
	levelMap := map[zapcore.Level]string{
		zapcore.DebugLevel: filepath + "/log-debug.log",
		zapcore.InfoLevel:  filepath + "/log-info.log",
		zapcore.WarnLevel:  filepath + "/log-warn.log",
		zapcore.ErrorLevel: filepath + "/log-error.log",
		zapcore.PanicLevel: filepath + "/log-panic.log",
	}

	//筛选level
	coreArr := []zapcore.Core{}
	for level, filename := range levelMap {
		if level >= defaultLevel {
			coreArr = append(coreArr, levelZapCore(encoder, level, filename))
		}
	}
	// 最后创建具体的Logger
	core := zapcore.NewTee(coreArr...)

	log := zap.New(core, zap.AddCaller()) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数
	logger = log.Sugar()
}

//文件 和 控制台 一起写入
func levelZapCore(encoder zapcore.Encoder, level zapcore.Level, filename string) zapcore.Core {
	levelFunc := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == level
	})
	writer := getWriter(filename)
	ws := zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer), zapcore.AddSync(os.Stdout))

	return zapcore.NewCore(encoder, ws, levelFunc)
}

//只是写入到文件
func levelZapCoreFile(encoder zapcore.Encoder, level zapcore.Level, filename string) zapcore.Core {
	levelFunc := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == level
	})
	writer := getWriter(filename)
	return zapcore.NewCore(encoder, zapcore.AddSync(writer), levelFunc)
}

//获取对于level日志文件的io.Writer 抽象
func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"-%Y%m%d.%H.log", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),                          //liunx 软链接 指向最新的log
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*1),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func ToLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	logger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	logger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}
