package logger

import (
	"fmt"
	"github.com/Myriad-Dreamin/minimum-lib/logger/bufferpool"
	kitzaplog "github.com/go-kit/kit/log/zap"
	colorful "github.com/gookit/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	// 	lightSkyBlue      = "87CEFA"
	lightlightSkyBlue = "B0E2FF"
	skyBlue           = "7EC0EE"
	pathColorPrefix   = "\x1b[38;5;230m "
	pathColorSuffix   = " \x1b[0m"
)

var (
	// colorLightSkyBlue      = colorful.HEX(lightSkyBlue)
	colorLightLightSkyBlue = colorful.HEX(lightlightSkyBlue)
	colorSkyBlue           = colorful.HEX(skyBlue)

	colorInfo              = colorLightLightSkyBlue.Sprintf("Info")
	colorDebug             = colorSkyBlue.Sprintf("Debug")

	colorWarn              = colorful.Yellow.Sprintf("Warn")
	colorPanic             = colorful.Red.Sprintf("Panic")
	colorError             = colorful.Red.Sprintf("Error")
	colorFatal             = colorful.Red.Sprintf("Fatal")
	colorDPanic            = colorful.Red.Sprintf("DPanic")
)

func zapColorfulLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString(colorDebug)
	case zapcore.InfoLevel:
		enc.AppendString(colorInfo)
	case zapcore.WarnLevel:
		enc.AppendString(colorWarn)
	case zapcore.ErrorLevel:
		enc.AppendString(colorError)
	case zapcore.DPanicLevel:
		enc.AppendString(colorDPanic)
	case zapcore.PanicLevel:
		enc.AppendString(colorPanic)
	case zapcore.FatalLevel:
		enc.AppendString(colorFatal)
	default:
		enc.AppendString(fmt.Sprintf("LEVEL(%d)", level))
	}
}

// ShortColorfulCallerEncoder serializes a caller in package/file:line format, trimming
// all but the final directory from the full path.
func ShortColorfulCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// TODO: consider using a byte-oriented API to save an allocation.
	if !caller.Defined {
		enc.AppendString("undefined")
		return
	}
	// nb. To make sure we trim the path correctly on Windows too, we
	// counter-intuitively need to use '/' and *not* os.PathSeparator here,
	// because the path given originates from Go stdlib, specifically
	// runtime.Caller() which (as of Mar/17) returns forward slashes even on
	// Windows.
	//
	// See https://github.com/golang/go/issues/3335
	// and https://github.com/golang/go/issues/18151
	//
	// for discussion on the issue on Go side.
	//
	// Find the last separator.
	//
	idx := strings.LastIndexByte(caller.File, '/')
	if idx == -1 {
		enc.AppendString(caller.FullPath())
	}
	// Find the penultimate separator.
	idx = strings.LastIndexByte(caller.File[:idx], '/')
	if idx == -1 {
		enc.AppendString(caller.FullPath())
	}
	buf := bufferpool.Get()
	buf.AppendString(pathColorPrefix)
	// Keep everything after the penultimate separator.
	buf.AppendString(caller.File[idx+1:])
	buf.AppendByte(':')
	buf.AppendInt(int64(caller.Line))
	buf.AppendString(pathColorSuffix)
	enc.AppendString(buf.String())
	buf.Free()
}

// FullColorfulCallerEncoder serializes a caller in package/file:line format, trimming
// all but the final directory from the full path.
func FullColorfulCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// TODO: consider using a byte-oriented API to save an allocation.
	if !caller.Defined {
		enc.AppendString("undefined")
		return
	}
	// nb. To make sure we trim the path correctly on Windows too, we
	// counter-intuitively need to use '/' and *not* os.PathSeparator here,
	// because the path given originates from Go stdlib, specifically
	// runtime.Caller() which (as of Mar/17) returns forward slashes even on
	// Windows.
	//
	// See https://github.com/golang/go/issues/3335
	// and https://github.com/golang/go/issues/18151
	//
	// for discussion on the issue on Go side.
	//
	// Find the last separator.
	//
	buf := bufferpool.Get()
	buf.AppendString(pathColorPrefix)
	// Keep everything after the penultimate separator.
	buf.AppendString(caller.File)
	buf.AppendByte(':')
	buf.AppendInt(int64(caller.Line))
	buf.AppendString(pathColorSuffix)
	enc.AppendString(buf.String())
	buf.Free()
}

func NewZapDevelopmentSugarOption() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapColorfulLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   ShortColorfulCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func NewZapLogger(cfg zap.Config, level zapcore.Level, options ...zap.Option) (Logger, error) {
	logger, err := cfg.Build(append(options, zap.AddCallerSkip(1))...)
	// zap.NewDevelopment(options...)
	if err != nil {
		return nil, err
	}
	return NewKitLogger(kitzaplog.NewZapSugarLogger(logger, level)), nil
}

func NewZapColorfulDevelopmentSugarLogger(options ...zap.Option) (Logger, error) {
	return NewZapLogger(NewZapDevelopmentSugarOption(), zapcore.DebugLevel, options...)
}
