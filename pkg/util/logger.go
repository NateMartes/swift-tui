package util

import (
	"github.com/NateMartes/swift-tui/pkg/errors"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

// cliHandler writes out to a specifc stream, holding a mutex to stop
// other processes from writing to the stream
type cliHandler struct {
	mu       *sync.Mutex
	out      io.Writer
	minLevel slog.Level
}

func (h *cliHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.minLevel
}
func (h *cliHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *cliHandler) WithGroup(name string) slog.Handler       { return h }

// Format logs as [LEVEL] msg
func (h *cliHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := fmt.Fprintf(h.out, "[%s] %s\n", r.Level.String(), r.Message)
	return err
}

// logStreamHandler is a type of stream logs can write to
type logStreamHandler struct {
	errHandler slog.Handler
	outHandler slog.Handler
	debugHandler slog.Handler
}

func (h *logStreamHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.errHandler.Enabled(ctx, level) || h.outHandler.Enabled(ctx, level) || h.debugHandler.Enabled(ctx, level)
}

func (h *logStreamHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= slog.LevelWarn {
		return h.errHandler.Handle(ctx, r)
	}
	if r.Level == slog.LevelDebug {
		return h.debugHandler.Handle(ctx, r)
	}
	return h.outHandler.Handle(ctx, r)
}

func (h *logStreamHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logStreamHandler{
		errHandler: h.errHandler.WithAttrs(attrs),
		outHandler: h.outHandler.WithAttrs(attrs),
		debugHandler: h.debugHandler.WithAttrs(attrs),
	}
}

func (h *logStreamHandler) WithGroup(name string) slog.Handler {
	return &logStreamHandler{
		errHandler: h.errHandler.WithGroup(name),
		outHandler: h.outHandler.WithGroup(name),
		debugHandler: h.debugHandler.WithGroup(name),
	}
}

var outHandler *cliHandler
var errHandler *cliHandler
var debugHandler *cliHandler
var debugFile *os.File = nil

// Inits the logger to be used by the CLI
func SetupLogger() {

	errHandler = &cliHandler{out: os.Stderr, mu: &sync.Mutex{}, minLevel: slog.LevelWarn}
	outHandler = &cliHandler{out: os.Stdout, mu: &sync.Mutex{}, minLevel: slog.LevelInfo}	
	logger := slog.New(&logStreamHandler{
		errHandler: errHandler,
		outHandler: outHandler,
	})

	slog.SetDefault(logger)
}

// Log writes to the default logger with INFO, should be called after SetupLogger()
func LogInfo(message string) {
	slog.Log(context.TODO(), slog.LevelInfo, message)
}

// Log writes to the default logger with DEBUG, should be called after SetupLogger()
func LogDebug(message string) {
	slog.Log(context.TODO(), slog.LevelDebug, message)
}

// Log writes to the default logger with WARN, should be called after SetupLogger()
func LogWarning(message string) {
	slog.Log(context.TODO(), slog.LevelWarn, message)
}

// Log writes to the default logger with ERROR, should be called after SetupLogger()
func LogError(message string) {
	slog.Log(context.TODO(), slog.LevelError, message)
}

// Logs to the error log stream in the default logger, exiting with a code aswell
func LogFatal(message string, exitCode int) {
	slog.Log(context.TODO(), slog.LevelError, message)
	os.Exit(exitCode)
}

// Set log level to debug, making a debug file on true
func SetDebugLogging(val bool) {

	if val {
		
		now := time.Now().Format(time.RFC3339)
		if debugFile == nil {
			var err error
			debugFile, err = os.Create(fmt.Sprintf("swift-tui-%s.log", now))
			if err != nil {
				LogFatal(fmt.Sprintf("Failed to make debug file: %s", err.Error()), errors.IO_ERROR)
			}
		}
		
		debugHandler := &cliHandler{out: debugFile, mu: &sync.Mutex{}, minLevel: slog.LevelDebug}
		logger := slog.New(&logStreamHandler{
			errHandler: errHandler,
			outHandler: outHandler,
			debugHandler: debugHandler,
		})
		slog.SetDefault(logger)
		
	} else {
		logger := slog.New(&logStreamHandler{
			errHandler: errHandler,
			outHandler: outHandler,
		})
		debugFile.Close()
		debugFile = nil
		slog.SetDefault(logger)	
	}
}

// Called when logging is finished
func CleanLogger() {
	if debugFile != nil {
		debugFile.Close()
		debugFile = nil
	}
}
