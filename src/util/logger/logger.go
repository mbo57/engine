package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	ERROR = iota + 1
	WARNING
	INFO
	DEBUG
)

type Log interface {
	Debug(msg string)
	Info(msg string)
	Warnig(msg string)
	Error(msg string)
}

type customLogger struct {
	logger *log.Logger
	level  int
}

func New() Log {
	logger := log.Default()
	logger.SetFlags(log.Ldate | log.Lmicroseconds)

	return &customLogger{
		logger: logger,
		level: func() int {
			switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
			case "ERROR":
				return ERROR
			case "WARNING":
				return WARNING
			case "INFO":
				return INFO
			case "DEBUG":
				return DEBUG
			default:
				return INFO
			}
		}(),
	}
}

func (l *customLogger) logOutput(prefix string, msg string) {
	l.logger.SetPrefix(fmt.Sprintf("[%s]", prefix))
	_, file, line, ok := runtime.Caller(1)
	execFileInfo := ""
	if ok {
		execFileInfo = fmt.Sprintf("%s:%d: ", filepath.Base(file), line)
	}
	logMsg := fmt.Sprintf(execFileInfo + msg)

	l.logger.SetOutput(io.MultiWriter(os.Stdout))
	l.logger.Printf(logMsg)
}

func (l *customLogger) Error(msg string) {
	if l.level < ERROR {
		return
	}
	l.logOutput("ERROR", msg)
}

func (l *customLogger) Warnig(msg string) {
	if l.level < WARNING {
		return
	}
	l.logOutput("WARNING", msg)
}

func (l *customLogger) Info(msg string) {
	if l.level < INFO {
		return
	}
	l.logOutput("INFO", msg)
}

func (l *customLogger) Debug(msg string) {
	if l.level < DEBUG {
		return
	}
	l.logOutput("DEBUG", msg)
}
