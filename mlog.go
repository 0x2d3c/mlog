package mlog

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
)

const (
	Trace = iota // Trace log level, must have a trace key with value in context
	Info         // Info level
	Warn         // Warn level
	Debug        // Debug level
	Error        // Error level
	Panic        // Panic level
	Fatal        // Fatal level

	red      = "\x1B[31m"
	green    = "\x1B[32m"
	yellow   = "\x1B[33m"
	blue     = "\x1B[34m"
	magenta  = "\x1B[35m"
	cyan     = "\x1B[36m"
	colorEnd = "\x1B[0m"
	newL     = "\n"
	blank    = " "
	end      = "]"
	caller   = "[caller:"
	trace    = "[trace:"
	INFO     = "[INFO] "
	WARN     = "[WARN] "
	DEBUG    = "[DEBUG]"
	ERROR    = "[ERROR]"
	PANIC    = "[PANIC]"
	FATAL    = "[FATAL]"

	errInfo  = green + INFO + colorEnd
	errWarn  = yellow + WARN + colorEnd
	errDebug = blue + DEBUG + colorEnd
	errError = red + ERROR + colorEnd
	errPanic = magenta + PANIC + colorEnd
	errFatal = cyan + FATAL + colorEnd
)

// Mlog logger data struct define
type Mlog struct {
	release         bool        // when value is true which error, fatal, panic will be write to log file
	lvl             int         // log level
	traceKey        string      // trace customer in log with key
	send            chan string // send err to chan
	io.StringWriter             // log file writer
}

// Option minlog cfg define
type Option struct {
	Release  bool            // when value is true which error, fatal, panic will be write to log file
	Lvl      int             // log level
	TraceKey string          // trace customer in log with key
	Writer   io.StringWriter // log writer
}

// New create minlog
func New(op *Option) *Mlog {
	m := &Mlog{
		lvl:          op.Lvl,
		release:      op.Release,
		traceKey:     op.TraceKey,
		send:         make(chan string),
		StringWriter: op.Writer,
	}
	go m.log()
	return m
}

func (m *Mlog) caller() string {
	_, file, line, _ := runtime.Caller(3)
	return file + ":" + strconv.Itoa(line)
}

func (m *Mlog) log() {
	for e := range m.send {
		m.WriteString(e)
	}
}

// Exit exit minlog
func (m *Mlog) Exit() {
	close(m.send)
}

func (m *Mlog) wrap(ctx context.Context, lvl int, tag, msg string, args ...interface{}) {

	// filter log level
	if m.lvl > lvl {
		return
	}

	if m.lvl == Trace {
		tag += trace + fmt.Sprint(ctx.Value(m.traceKey)) + end
	}

	switch lvl {
	case Error, Panic, Fatal:
		tag += caller + m.caller() + end + blank
	}

	msg = tag + msg + newL
	if len(args) != 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	if m.release {
		m.send <- msg
	} else {
		fmt.Fprint(os.Stderr, msg)
	}
}

// Info info log method
func (m *Mlog) Info(ctx context.Context, msg string, args ...interface{}) {
	m.wrap(ctx, Info, errInfo, msg, args...)
}

// Warn warn log method
func (m *Mlog) Warn(ctx context.Context, msg string, args ...interface{}) {
	m.wrap(ctx, Warn, errWarn, msg, args...)
}

// Debug debug log method
func (m *Mlog) Debug(ctx context.Context, msg string, args ...interface{}) {
	m.wrap(ctx, Debug, errDebug, msg, args...)
}

// Error err log method
func (m *Mlog) Error(ctx context.Context, msg string, args ...interface{}) {
	m.wrap(ctx, Error, errError, msg, args...)
}

// Panic panic log method
func (m *Mlog) Panic(ctx context.Context, msg string, args ...interface{}) {
	m.wrap(ctx, Panic, errPanic, msg, args...)
	panic(fmt.Errorf(msg, args...))
}

// Fatal fatal log method
func (m *Mlog) Fatal(ctx context.Context, msg string, args ...interface{}) {
	m.wrap(ctx, Fatal, errFatal, msg, args...)
	os.Exit(0)
}
