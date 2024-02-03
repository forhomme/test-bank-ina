package baselogger

import (
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger() *Logger {
	loggr := &logrus.Logger{
		Out:          io.MultiWriter(os.Stdout),
		ReportCaller: true,
		Level:        logrus.DebugLevel,
		Formatter: &logrus.JSONFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				s := strings.Split(frame.Function, ".")
				funcname := s[len(s)-1]
				fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
				return funcname, fileName
			},
			TimestampFormat: "2006-01-02 15:04:05.000",
		},
		Hooks: make(map[logrus.Level][]logrus.Hook),
	}

	// Instrument logrus.
	loggr.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	)))

	return &Logger{
		Logger: loggr,
	}
}
