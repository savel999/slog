package main

import (
	"context"
	"fmt"
	"github.com/dusted-go/logging/prettylog"
	"log/slog"
	"os"
)

type ILogger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)

	DebugContext(ctx context.Context, msg string, fields ...any)
	InfoContext(ctx context.Context, msg string, fields ...any)
	WarnContext(ctx context.Context, msg string, fields ...any)
	ErrorContext(ctx context.Context, msg string, fields ...any)
}

type loggerMode string
type AttrHandlerFn func(groups []string, a slog.Attr) slog.Attr
type RecordHandlerFn func(ctx context.Context, r slog.Record) slog.Record

const (
	LoggerModeJSON   loggerMode = "json"
	LoggerModeText   loggerMode = "text"
	LoggerModePretty loggerMode = "pretty"
)

type Options struct {
	Mode           loggerMode
	Level          slog.Level
	RecordHandlers []RecordHandlerFn
	AttrHandlers   []AttrHandlerFn
}

func NewLogger(o Options) ILogger {
	attrHandlers := append(o.AttrHandlers, handleSourceAttr())

	opts := &slog.HandlerOptions{
		Level:       o.Level,
		AddSource:   true,
		ReplaceAttr: replaceAttr(attrHandlers...),
	}

	w := os.Stdout

	switch o.Mode {
	case LoggerModeJSON:
		return slog.New(handler{slog.NewJSONHandler(w, opts), o.RecordHandlers})
	case LoggerModePretty:
		return slog.New(handler{prettylog.NewHandler(opts), o.RecordHandlers})
	default:
		return slog.New(handler{slog.NewTextHandler(w, opts), o.RecordHandlers})
	}
}

func replaceAttr(fns ...AttrHandlerFn) AttrHandlerFn {
	return func(groups []string, a slog.Attr) slog.Attr {
		for _, fn := range fns {
			a = fn(groups, a)
		}

		return a
	}
}

func handleSourceAttr() AttrHandlerFn {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			if source, ok := a.Value.Any().(*slog.Source); ok {
				return slog.Attr{
					Key:   "caller",
					Value: slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line)),
				}
			}
		}

		return a
	}
}

type handler struct {
	slog.Handler
	fns []RecordHandlerFn
}

func (h handler) Handle(ctx context.Context, r slog.Record) error {
	for _, fn := range h.fns {
		r = fn(ctx, r)
	}

	return h.Handler.Handle(ctx, r)
}
