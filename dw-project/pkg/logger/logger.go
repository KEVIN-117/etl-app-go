package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// ANSI color codes
const (
	reset = "\033[0m"
	bold  = "\033[1m"
	dim   = "\033[2m"

	fgBlack   = "\033[30m"
	fgRed     = "\033[31m"
	fgGreen   = "\033[32m"
	fgYellow  = "\033[33m"
	fgBlue    = "\033[34m"
	fgMagenta = "\033[35m"
	fgCyan    = "\033[36m"
	fgWhite   = "\033[37m"

	bgRed     = "\033[41m"
	bgGreen   = "\033[42m"
	bgYellow  = "\033[43m"
	bgBlue    = "\033[44m"
	bgMagenta = "\033[45m"
	bgCyan    = "\033[46m"
)

// Level represents the severity of a log message.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelMeta = map[Level]struct {
	label string
	color string
}{
	LevelDebug: {label: "DEBUG", color: fgCyan},
	LevelInfo:  {label: "INFO ", color: fgGreen},
	LevelWarn:  {label: "WARN ", color: fgYellow},
	LevelError: {label: "ERROR", color: fgRed},
	LevelFatal: {label: "FATAL", color: bold + fgMagenta},
}

// Field is a key-value pair for structured logging.
type Field struct {
	Key   string
	Value any
}

// F creates a named field for structured log output.
func F(key string, value any) Field {
	return Field{Key: key, Value: value}
}

// Options configures a Logger instance.
type Options struct {
	// Minimum log level to emit (default: LevelDebug)
	Level Level
	// Writer to send output to (default: os.Stdout)
	Output io.Writer
	// Include caller file:line in output (default: true)
	ShowCaller bool
	// Include timestamp in output (default: true)
	ShowTime bool
	// Disable ANSI colors (default: false)
	NoColor bool
	// Timestamp layout (default: "2006-01-02 15:04:05")
	TimeFormat string
	// How many stack frames to skip when resolving caller (default: 2)
	CallerSkip int
}

// Logger is a thread-safe, colorful, leveled CLI logger.
type Logger struct {
	mu   sync.Mutex
	opts Options
}

// New creates a Logger with the given options.
// Fields left at their zero value get sensible defaults.
func New(opts Options) *Logger {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}
	if opts.TimeFormat == "" {
		opts.TimeFormat = "2006-01-02 15:04:05"
	}
	if opts.CallerSkip == 0 {
		opts.CallerSkip = 2
	}
	if !opts.ShowCaller {
		// default on
		opts.ShowCaller = true
	}
	if !opts.ShowTime {
		opts.ShowTime = true
	}
	return &Logger{opts: opts}
}

// Default is a ready-to-use logger with sensible defaults.
var Default = New(Options{
	ShowCaller: true,
	ShowTime:   true,
})

// colorize wraps text with an ANSI code if colors are enabled.
func (l *Logger) colorize(code, text string) string {
	if l.opts.NoColor {
		return text
	}
	return code + text + reset
}

// log is the internal write method.
func (l *Logger) log(level Level, msg string, fields []Field) {
	if level < l.opts.Level {
		return
	}

	meta, ok := levelMeta[level]
	if !ok {
		return
	}

	var sb strings.Builder

	// Timestamp
	if l.opts.ShowTime {
		ts := time.Now().Format(l.opts.TimeFormat)
		sb.WriteString(l.colorize(dim, ts))
		sb.WriteByte(' ')
	}

	// Level badge
	badge := fmt.Sprintf(" %s ", meta.label)
	sb.WriteString(l.colorize(bold+meta.color, badge))
	sb.WriteByte(' ')

	// Caller
	if l.opts.ShowCaller {
		_, file, line, ok := runtime.Caller(l.opts.CallerSkip)
		if ok {
			caller := fmt.Sprintf("%s:%d", filepath.Base(file), line)
			sb.WriteString(l.colorize(dim, caller))
			sb.WriteByte(' ')
		}
	}

	// Message — FATAL gets extra emphasis
	if level == LevelFatal {
		sb.WriteString(l.colorize(bold+fgRed, msg))
	} else {
		sb.WriteString(msg)
	}

	// Structured fields
	for _, f := range fields {
		sb.WriteByte(' ')
		sb.WriteString(l.colorize(fgCyan, f.Key))
		sb.WriteByte('=')
		sb.WriteString(l.colorize(fgYellow, fmt.Sprintf("%v", f.Value)))
	}

	sb.WriteByte('\n')

	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprint(l.opts.Output, sb.String())

	if level == LevelFatal {
		os.Exit(1)
	}
}

// Debug logs a debug-level message with optional fields.
func (l *Logger) Debug(msg string, fields ...Field) { l.log(LevelDebug, msg, fields) }

// Info logs an info-level message with optional fields.
func (l *Logger) Info(msg string, fields ...Field) { l.log(LevelInfo, msg, fields) }

// Warn logs a warning-level message with optional fields.
func (l *Logger) Warn(msg string, fields ...Field) { l.log(LevelWarn, msg, fields) }

// Error logs an error-level message with optional fields.
func (l *Logger) Error(msg string, fields ...Field) { l.log(LevelError, msg, fields) }

// Fatal logs a fatal message and calls os.Exit(1).
func (l *Logger) Fatal(msg string, fields ...Field) { l.log(LevelFatal, msg, fields) }

// Debugf formats and logs a debug message (no structured fields).
func (l *Logger) Debugf(format string, args ...any) {
	l.log(LevelDebug, fmt.Sprintf(format, args...), nil)
}

// Infof formats and logs an info message.
func (l *Logger) Infof(format string, args ...any) {
	l.log(LevelInfo, fmt.Sprintf(format, args...), nil)
}

// Warnf formats and logs a warning message.
func (l *Logger) Warnf(format string, args ...any) {
	l.log(LevelWarn, fmt.Sprintf(format, args...), nil)
}

// Errorf formats and logs an error message.
func (l *Logger) Errorf(format string, args ...any) {
	l.log(LevelError, fmt.Sprintf(format, args...), nil)
}

// SetLevel changes the minimum log level at runtime.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.opts.Level = level
}

// SetOutput redirects log output to a new writer.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.opts.Output = w
}

// WithFields returns a child logger that always appends the given fields.
func (l *Logger) WithFields(fields ...Field) *ChildLogger {
	return &ChildLogger{parent: l, fields: fields}
}

// ChildLogger wraps a Logger with pre-attached fields (e.g. per-request context).
type ChildLogger struct {
	parent *Logger
	fields []Field
}

func (c *ChildLogger) merge(extra []Field) []Field {
	all := make([]Field, 0, len(c.fields)+len(extra))
	all = append(all, c.fields...)
	all = append(all, extra...)
	return all
}

func (c *ChildLogger) Debug(msg string, fields ...Field) {
	c.parent.log(LevelDebug, msg, c.merge(fields))
}
func (c *ChildLogger) Info(msg string, fields ...Field) {
	c.parent.log(LevelInfo, msg, c.merge(fields))
}
func (c *ChildLogger) Warn(msg string, fields ...Field) {
	c.parent.log(LevelWarn, msg, c.merge(fields))
}
func (c *ChildLogger) Error(msg string, fields ...Field) {
	c.parent.log(LevelError, msg, c.merge(fields))
}
func (c *ChildLogger) Fatal(msg string, fields ...Field) {
	c.parent.log(LevelFatal, msg, c.merge(fields))
}

// Package-level helpers using the Default logger.
func Debug(msg string, fields ...Field) { Default.log(LevelDebug, msg, fields) }
func Info(msg string, fields ...Field)  { Default.log(LevelInfo, msg, fields) }
func Warn(msg string, fields ...Field)  { Default.log(LevelWarn, msg, fields) }
func Error(msg string, fields ...Field) { Default.log(LevelError, msg, fields) }
func Debugf(format string, args ...any) { Default.log(LevelDebug, fmt.Sprintf(format, args...), nil) }
func Infof(format string, args ...any)  { Default.log(LevelInfo, fmt.Sprintf(format, args...), nil) }
func Warnf(format string, args ...any)  { Default.log(LevelWarn, fmt.Sprintf(format, args...), nil) }
func Errorf(format string, args ...any) { Default.log(LevelError, fmt.Sprintf(format, args...), nil) }
