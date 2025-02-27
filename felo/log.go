package felo

import (
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

//https://juejin.cn/post/7080869429805842469

var (
	logrusPackage      string
	minimumCallerDepth int
	callerInitOnce     sync.Once
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4

	LogrusPrefix = "github.com/sirupsen/logrus"

	FieldKeyCallerFunc = "caller_func"
	FieldKeyCallerFile = "caller_file"
)

type CallerHook struct{}

func (h *CallerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}

func (h *CallerHook) Fire(entry *logrus.Entry) error {
	caller := getCaller()
	funcVal, fileVal := prettierCaller(caller)

	// 暂时按照text的formatter格式填充
	if funcVal != "" {
		entry.Data[FieldKeyCallerFunc] = funcVal
	}
	if fileVal != "" {
		entry.Data[FieldKeyCallerFile] = fileVal
	}
	return nil
}

func prettierCaller(f *runtime.Frame) (string, string) {
	_, file, _, _ := runtime.Caller(0)
	prefix := filepath.Dir(file) + "/"
	function := strings.TrimPrefix(f.Function, prefix) + "()"
	fileLine := strings.TrimPrefix(f.File, prefix) + ":" + strconv.Itoa(f.Line)
	return function, fileLine
}

// getCaller 参考logrus.getCaller()改造
func getCaller() *runtime.Frame {
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		minimumCallerDepth = knownLogrusFrames
		_ = runtime.Callers(0, pcs)

		var alreadyFindLogurs bool
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()

			// 找到logurs包上一级调用方
			if strings.HasPrefix(funcName, LogrusPrefix) {
				alreadyFindLogurs = true
				continue
			}
			if alreadyFindLogurs {
				logrusPackage = getPackageName(funcName)
				return
			}
		}
	})

	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		if strings.HasPrefix(f.Function, LogrusPrefix) {
			continue
		}

		pkg := getPackageName(f.Function)
		if pkg != logrusPackage {
			return &f
		}
	}
	return nil
}

func getPackageName(f string) string {
	lastSlash := strings.LastIndex(f, "/")
	if lastSlash == -1 {
		lastSlash = 0
	} else {
		lastSlash++
	}
	lastPeriod := strings.Index(f[lastSlash:], ".")
	if lastPeriod == -1 {
		return f
	}
	return f[:lastSlash+lastPeriod]
}
