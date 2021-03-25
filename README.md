### mlog

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
Writer:   assertFile("min.log"),
})

minlog.Info(ctx, "this is New info test")
minlog.Infof(ctx, "this is New %s test", "infof")
minlog.Warn(ctx, "this is New warn test")
minlog.Warnf(ctx, "this is New %s test", "warnf")
minlog.Error(ctx, "this is New error test")
minlog.Errorf(ctx, "this is New %s test", "errorf")
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
 - info-12              315081              5452 ns/op             144 B/op          3 allocs/op
 - debug-12             236989              4344 ns/op             144 B/op          3 allocs/op
 - error-12             187606              6633 ns/op             704 B/op          8 allocs/op
 - warn-12              298220              4073 ns/op             144 B/op          3 allocs/op
BenchmarkNewWithF:
 - infof-12             304359              4655 ns/op             192 B/op          4 allocs/op
 - debugf-12            297622              4515 ns/op             192 B/op          4 allocs/op
 - errorf-12            169093             10873 ns/op             816 B/op          9 allocs/op
 - warnf-12             297420              5000 ns/op             192 B/op          4 allocs/op
BenchmarkNewParallelWithOutF:
 - info-12              383181              3008 ns/op             144 B/op          3 allocs/op
 - debug-12             507474              2995 ns/op             144 B/op          3 allocs/op
 - error-12             391755              3687 ns/op             704 B/op          8 allocs/op
 - warn-12              403269              3005 ns/op             144 B/op          3 allocs/op
BenchmarkNewParallelWithF:
 - infof-12             452118              3080 ns/op             192 B/op          4 allocs/op
 - debugf-12            386314              3089 ns/op             192 B/op          4 allocs/op
 - errorf-12            330423              3664 ns/op             816 B/op          9 allocs/op
 - warnf-12             378901              3099 ns/op             192 B/op          4 allocs/op
```

### example with gin

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

var minlog *mlog.Mlog

func init() {
	minlog = mlog.New(&mlog.Option{
		Lvl:      mlog.Info,
		Writer:   file(),
		Release:  false,
		TraceKey: "min-example",
	})
}
func file() io.StringWriter {
	name := "min.log"
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

	minlog.Info(cc, fmt.Sprintf("[uri:%s][body: %s]", uri, body))
	minlog.Debug(cc, fmt.Sprintf("[uri:%s][body: %s]", uri, body))
	minlog.Error(cc, fmt.Sprintf("[uri:%s][body: %s]", uri, body))

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