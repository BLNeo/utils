package template

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
	"time"
)

var (
	Log *Logger = initLog(true,"default")
	once *sync.Once
)

type Logger struct {
	name  string
	debug bool
	*zap.Logger
}


func initLog (debug bool,serviceName string)*Logger {
	l :=&Logger{
		debug: debug,
		name:  serviceName,
	}
	return l.Init()
}



func LogSetMode(debug bool,serviceName string)*Logger {
	l := initLog(debug,serviceName)
	if Log !=nil {
		Log = l
	}
	return l
}


func (log *Logger) Init() *Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	cores := []zapcore.Core{}

	//设置info waring和error的日志
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})

	waringLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel
	})

	//debug为ture  日志输出到终端
	if log.debug {
		//debug 直接输出到终端中
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout)),
			infoLevel))
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout)),
			waringLevel))
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout)),
			errorLevel))

	} else {
		// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
		infoWriter := getWriter(log.name + "_info.log")
		waringWriter := getWriter(log.name + "_waring.log")
		errorWriter := getWriter(log.name + "_error.log")
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(&infoWriter)),
			infoLevel))

		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(&waringWriter)),
			waringLevel))
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(&errorWriter)),
			errorLevel))
	}

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		cores...,
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	logger := zap.New(core, caller, development)
	log.Logger = logger
	return log
}

func (log *Logger) Close() error {
	return nil
}

func getWriter(filename string) lumberjack.Logger {
	today := time.Now().Format("20060102")
	filename = fmt.Sprintf("./logs/%s/%s", today, filename)
	return lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    128,      // 每个日志文件保存的最大尺寸 单位：M  128
		MaxBackups: 30,       // 日志文件最多保存多少个备份 30
		MaxAge:     7,        // 文件最多保存多少天 7
		Compress:   true,     // 是否压缩
	}
}
