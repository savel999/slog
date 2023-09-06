package main

import (
	"context"
	"fmt"
	"log/slog"
)

type Example struct {
	A int
	B string
	C []string
	D bool
	E float64
	F Example2
}

type Example2 struct {
	MaxInt uint64
}

func main() {
	ctx := context.WithValue(context.TODO(), "foo", "bar")
	ctx = context.WithValue(ctx, "bar", "djigurda")

	example := Example{A: 123, B: "str", C: []string{"23d", "sdfs"}, E: 1.34, F: Example2{
		MaxInt: 1<<64 - 1,
	}}

	recordhandlers := []RecordHandlerFn{
		logCtx(),
		log2Ctx(),
	}

	attrHandlers := []AttrHandlerFn{
		handleTimeAttr(),
	}

	fmt.Println("JSON логгер\n")

	o := Options{Mode: LoggerModeJSON, Level: slog.LevelDebug, RecordHandlers: recordhandlers, AttrHandlers: attrHandlers}
	logger := NewLogger(o)
	logger.Debug("Debug message", slog.String("sdf", "sdfsdf"), slog.Any("config", example))
	logger.DebugContext(ctx, "Debug message", slog.String("sdf", "sdfsdf"))
	logger.Info("Info message", slog.String("sdf", "sdfsdf"))
	logger.Warn("Warning message")
	logger.Error("Error message")

	fmt.Println("\nPretty логгер\n")

	o2 := Options{Mode: LoggerModePretty, Level: slog.LevelDebug, RecordHandlers: recordhandlers}
	logger2 := NewLogger(o2)
	logger2.Debug("Debug message", slog.String("sdf", "sdfsdf"), slog.Any("config", example))
	logger2.DebugContext(ctx, "Debug message", slog.String("sdf", "sdfsdf"))
	logger2.Info("Info message", slog.String("sdf", "sdfsdf"))
	logger2.Warn("Warning message")
	logger2.Error("Error message")

	fmt.Println("\nТекстовый логгер\n")

	o3 := Options{Mode: LoggerModeText, Level: slog.LevelInfo}
	logger3 := NewLogger(o3)
	logger3.Debug("Debug message", slog.String("sdf", "sdfsdf"), slog.Any("config", example))
	logger3.DebugContext(ctx, "Debug message", slog.String("sdf", "sdfsdf"))
	logger3.Info("Info message", slog.String("sdf", "sdfsdf"))
	logger3.Warn("Warning message")
	logger3.Error("Error message")
}

func logCtx() RecordHandlerFn {
	return func(ctx context.Context, r slog.Record) slog.Record {
		if val, ok := ctx.Value("foo").(string); ok {
			r.AddAttrs(slog.Attr{Key: "foo", Value: slog.StringValue(val)})
		}

		return r
	}
}

func log2Ctx() RecordHandlerFn {
	return func(ctx context.Context, r slog.Record) slog.Record {
		if val, ok := ctx.Value("bar").(string); ok {
			r.AddAttrs(slog.Attr{Key: "bar", Value: slog.StringValue(val)})
		}

		return r
	}
}

func handleTimeAttr() AttrHandlerFn {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}

		return a
	}
}
