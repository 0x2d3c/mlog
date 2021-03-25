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
	_  = 1 << (iota * 10)
	KB // 1KB
	MB // 1MB
	GB // 1GB

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
	caller   = "[caller: "
	trace    = "[trace: "
	INFO     = "[INFO]"
	WARN     = "[WARN]"
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

// Mlog minlog define
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
	_, file, line, _ := runtime.Caller(4)
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

func (m *Mlog) wrapCommon(ctx context.Context, lvl int, tag, msg string) string {
	// trace flag
	var traceTag string
	if m.lvl == Trace {
		traceTag = trace + fmt.Sprint(ctx.Value(m.traceKey)) + end
	}

	var callerTag string
	switch lvl {
	default:
		callerTag += blank
	case Debug:
	case Error, Panic, Fatal:
		callerTag += caller + m.caller() + end
	}
	return traceTag + tag + callerTag + msg
}

func (m *Mlog) wrap(ctx context.Context, lvl int, tag, msg string) {
	// filter log level
	if m.lvl > lvl {
		return
	}

	msg = m.wrapCommon(ctx, lvl, tag, msg) + newL
	if m.release {
		m.send <- msg
	} else {
		fmt.Fprint(os.Stderr, msg)
	}
}

func (m *Mlog) wrapf(ctx context.Context, lvl int, tag, msg string, args ...interface{}) {

	// filter log level
	if m.lvl > lvl {
		return
	}

	msg = fmt.Sprintf(m.wrapCommon(ctx, lvl, tag, msg), args...) + newL
	if m.release {
		m.send <- msg
	} else {
		fmt.Fprint(os.Stderr, msg)
	}
}

// Info info log method
func (m *Mlog) Info(ctx context.Context, msg string) {
	m.wrap(ctx, Info, errInfo, msg)
}

func (m *Mlog) Infof(ctx context.Context, msg string, args ...interface{}) {
	m.wrapf(ctx, Info, errInfo, msg, args...)
}

// Warn warn log method
func (m *Mlog) Warn(ctx context.Context, msg string) {
	m.wrap(ctx, Warn, errWarn, msg)
}

func (m *Mlog) Warnf(ctx context.Context, msg string, args ...interface{}) {
	m.wrapf(ctx, Warn, errWarn, msg, args...)
}

// Debug debug log method
func (m *Mlog) Debug(ctx context.Context, msg string) {
	m.wrap(ctx, Debug, errDebug, msg)
}

func (m *Mlog) Debugf(ctx context.Context, msg string, args ...interface{}) {
	m.wrapf(ctx, Debug, errDebug, msg, args...)
}

// Error err log method
func (m *Mlog) Error(ctx context.Context, msg string) {
	m.wrap(ctx, Error, errError, msg)
}

func (m *Mlog) Errorf(ctx context.Context, msg string, args ...interface{}) {
	m.wrapf(ctx, Error, errError, msg, args...)
}

// Panic panic log method
func (m *Mlog) Panic(ctx context.Context, msg string) {
	m.wrap(ctx, Panic, errPanic, msg)
	panic(fmt.Errorf(msg))
}

func (m *Mlog) Panicf(ctx context.Context, msg string, args ...interface{}) {
	m.wrapf(ctx, Panic, errPanic, msg, args...)
	panic(fmt.Errorf(msg, args...))
}

// Fatal fatal log method
func (m *Mlog) Fatal(ctx context.Context, msg string) {
	m.wrap(ctx, Fatal, errFatal, msg)
	os.Exit(0)
}

func (m *Mlog) Fatalf(ctx context.Context, msg string, args ...interface{}) {
	m.wrapf(ctx, Fatal, errFatal, msg, args...)
	os.Exit(0)
}
