package mlog

import (
	"context"
	"io"
	"os"
	"testing"
	"time"
)

func assertFile(name string) io.StringWriter {
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if os.IsNotExist(err) {
		file, _ = os.Create(name)
	}
	return file
}

func assertNew() *Mlog {
	return New(&Option{
		Release:  true,
		Lvl:      Info,
		TraceKey: "userID",
		Writer:   assertFile("min.log"),
	})
}

func TestNew(t *testing.T) {
	minLog := assertNew()
	ctx := context.WithValue(context.Background(), "userID", "9527")
	minLog.Info(ctx, "this is New info test")
	minLog.Infof(ctx, "this is New %s test", "infof")
	minLog.Debug(ctx, "this is New debug test")
	minLog.Debugf(ctx, "this is New %s test", "debugf")
	minLog.Error(ctx, "this is New error test")
	minLog.Errorf(ctx, "this is New %s test", "errorf")
	minLog.Warn(ctx, "this is New warn test")
	minLog.Warnf(ctx, "this is New %s test", "warnf")
	minLog.Exit()
	time.Sleep(time.Second)
}

func BenchmarkNewWithOutF(b *testing.B) {
	minLog := assertNew()
	ctx := context.WithValue(context.Background(), "userID", "9527")
	b.Run("info", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {
			minLog.Info(ctx, "this is New info test")
		}
	})
	b.Run("debug", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {
			minLog.Debug(ctx, "this is New debug test")
		}
	})
	b.Run("error", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {
			minLog.Error(ctx, "this is New error test")
		}
	})
	b.Run("warn", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {
			minLog.Warn(ctx, "this is New warn test")
		}
	})
	minLog.Exit()
}

func BenchmarkNewWithF(b *testing.B) {
	minLog := assertNew()
	ctx := context.WithValue(context.Background(), "userID", "9527")
	b.Run("infof", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {
			minLog.Infof(ctx, "this is New %s test", "infof")
		}
	})
	b.Run("debugf", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {
			minLog.Debugf(ctx, "this is New %s test", "errorf")
		}
	})
	b.Run("errorf", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {
			minLog.Errorf(ctx, "this is New %s test", "errorf")
		}
	})
	b.Run("warnf", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {
			minLog.Warnf(ctx, "this is New %s test", "warnf")
		}
	})
	minLog.Exit()
}

func BenchmarkNewParallelWithOutF(b *testing.B) {
	minLog := assertNew()
	ctx := context.WithValue(context.Background(), "userID", "9527")
	b.Run("info", func(bb *testing.B) {
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				minLog.Info(ctx, "this is New info test")
			}
		})
	})
	b.Run("debug", func(bb *testing.B) {
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				minLog.Debug(ctx, "this is New debug test")
			}
		})
	})
	b.Run("error", func(bb *testing.B) {
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				minLog.Error(ctx, "this is New error test")
			}
		})
	})
	b.Run("warn", func(bb *testing.B) {
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				minLog.Warn(ctx, "this is New warn test")
			}
		})
	})
	minLog.Exit()
}

func BenchmarkNewParallelWithF(b *testing.B) {
	minLog := assertNew()
	ctx := context.WithValue(context.Background(), "userID", "9527")
	b.Run("infof", func(bb *testing.B) {
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				minLog.Infof(ctx, "this is New %s test", "infof")
			}
		})
	})
	b.Run("debugf", func(bb *testing.B) {
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				minLog.Debugf(ctx, "this is New %s test", "debugf")
			}
		})
	})
	b.Run("errorf", func(bb *testing.B) {
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				minLog.Errorf(ctx, "this is New %s test", "errorf")
			}
		})
	})
	b.Run("warnf", func(bb *testing.B) {
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				minLog.Warnf(ctx, "this is New %s test", "warnf")
			}
		})
	})
	minLog.Exit()
}
