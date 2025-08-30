package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"strings"
)

type Config struct {
	Level         string
	Filename      string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Compress      bool
	ConsoleOutput bool
}

var sugar *zap.SugaredLogger

// main.go 에서 initialize 함, 설정 파일에서 로그레벨 불러와서 설정
func LoggerInit(loggerConfig Config) {
	strLevel := strings.ToLower(loggerConfig.Level)
	var level zap.AtomicLevel
	switch strLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel) // 기본값
	}
	// lumberjack 설정
	rotatingLogger := &lumberjack.Logger{
		Filename:   loggerConfig.Filename,   // 로그 파일 경로
		MaxSize:    loggerConfig.MaxSize,    // 최대 크기(MB)
		MaxBackups: loggerConfig.MaxBackups, // 보관할 백업 파일 수
		MaxAge:     loggerConfig.MaxAge,     // 보관할 최대 일수
		Compress:   loggerConfig.Compress,   // 압축 여부
	}

	// 파일 출력과 콘솔 출력을 동시에 처리
	fileSyncer := zapcore.AddSync(rotatingLogger)
	consoleSyncer := zapcore.AddSync(log.Writer())

	// 인코더 설정
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Core 구성
	var core zapcore.Core
	if loggerConfig.ConsoleOutput == true {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),                  // JSON 형식 로그
			zapcore.NewMultiWriteSyncer(fileSyncer, consoleSyncer), // 파일과 콘솔 동시 출력
			//zapcore.NewMultiWriteSyncer(fileSyncer),
			level, // 로그 레벨
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), // JSON 형식 로그
			zapcore.NewMultiWriteSyncer(fileSyncer),
			level, // 로그 레벨
		)
	}

	// Logger 생성
	logger := zap.New(core)
	defer logger.Sync()

	// SugarLogger로 로깅
	sugar = logger.Sugar()
	sugar.Infow("Logger initialized",
		"level", loggerConfig.Level,
		"output", loggerConfig.Filename,
	)

	// 샘플 로그
	sugar.Infow("Hello world")
}

func Debug(msg string, keysAndValues ...interface{}) {
	sugar.Debugw(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	sugar.Infow(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	sugar.Warnw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	sugar.Errorw(msg, keysAndValues...)
}
