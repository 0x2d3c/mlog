### Mlog

- use std lib impl logger with go, help you to collect your log file with configurable
- mlog have dev|release mode, dev mode don't save log to file but to std, release mode will save to file
- support trace request in log

### How to use

```go
ctx := context.WithValue(context.Background(), "userID", "9527")
mlog := New(&Option{
Release:  true,
Lvl:      Info,
TraceKey: "userID",
Writer:   assertFile("m.log"),
})

mlog.Info(ctx, "this is New info test")
mlog.Info(ctx, "this is New %s test", "infof")
mlog.Warn(ctx, "this is New warn test")
mlog.Warn(ctx, "this is New %s test", "warnf")
mlog.Error(ctx, "this is New error test")
mlog.Error(ctx, "this is New %s test", "errorf")
```

```shell script
go test -benchmem -run=^$mlog -bench '^(BenchmarkNewWithOutF)$'
go test -benchmem -run=^$mlog -bench '^(BenchmarkNewWithF)$'
go test -benchmem -run=^$mlog -bench '^(BenchmarkNewParallelWithOutF)$'
go test -benchmem -run=^$mlog -bench '^(BenchmarkNewParallelWithF)$'
```

### Test with std close and no trace tag

```
goos: windows
goarch: amd64
pkg: mlog
cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
BenchmarkNewWithOutF:
 - info-12              321913              3652 ns/op              96 B/op          2 allocs/op
 - debug-12             312220              3661 ns/op              96 B/op          2 allocs/op
 - error-12             212070              5508 ns/op             536 B/op          6 allocs/op
 - warn-12              343502              3619 ns/op              96 B/op          2 allocs/op
BenchmarkNewWithF:
 - infof-12             300884              3909 ns/op             144 B/op          3 allocs/op
 - debugf-12            308331              3886 ns/op             144 B/op          3 allocs/op
 - errorf-12            190932              6331 ns/op             632 B/op          7 allocs/op
 - warnf-12             308222              3837 ns/op             144 B/op          3 allocs/op
BenchmarkNewParallelWithOutF:
 - info-12              410865              2964 ns/op              96 B/op          2 allocs/op
 - debug-12             407143              3007 ns/op              96 B/op          2 allocs/op
 - error-12             353847              3386 ns/op             544 B/op          7 allocs/op
 - warn-12              414385              2936 ns/op              96 B/op          2 allocs/op
BenchmarkNewParallelWithF:
 - infof-12             400852              3014 ns/op             144 B/op          3 allocs/op
 - debugf-12            414873              3009 ns/op             144 B/op          3 allocs/op
 - errorf-12            353948              3486 ns/op             640 B/op          8 allocs/op
 - warnf-12             387709              3056 ns/op             144 B/op          3 allocs/op
```

### Example with Gin

```go
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/0x2d3c/mlog"
	"github.com/gin-gonic/gin"
)

var ml *mlog.Mlog

func init() {
	ml = mlog.New(&mlog.Option{
		Lvl:      mlog.Info,
		Writer:   file(),
		Release:  false,
		TraceKey: "m.key",
	})
}
func file() io.StringWriter {
	name := "m.log"
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if os.IsNotExist(err) {
		f, _ = os.Create(name)
	}
	return f
}

func RecordRequest(ctx *gin.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)

	uri := ctx.Request.RequestURI
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

	cc := ctx.Copy()

	ml.Info(cc, fmt.Sprintf("[uri:%s][body: %s]", uri, body))
	ml.Debug(cc, fmt.Sprintf("[uri:%s][body: %s]", uri, body))
	ml.Error(cc, fmt.Sprintf("[uri:%s][body: %s]", uri, body))

	ctx.Next()
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(RecordRequest)

	r.POST("/ping", func(ctx *gin.Context) {
		var data interface{}
		ctx.Bind(&data)
		ctx.JSON(200, gin.H{"message": data})
	})

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()
}
```