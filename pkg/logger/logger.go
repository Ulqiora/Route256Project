package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"
)

const (
	L_ERROR   string = "ERROR"
	L_DEBUG   string = "DEBUG"
	L_INFO    string = "INFO"
	L_WARNING string = "WARNING"
)

const bufferSize = 20

type message struct {
	logType    string
	ThreadName string
	Content    string
	Objects    []any
}

type logFunction func(msg string, args ...any)

type LoggerAsync struct {
	log             *slog.Logger
	closer          io.Closer
	dataLog         chan message
	mapLogFunctions map[string]logFunction
}

func New(ctx context.Context,
	w io.WriteCloser,
	options *slog.HandlerOptions) *LoggerAsync {
	var logger = LoggerAsync{
		log:             slog.New(slog.NewJSONHandler(w, options)),
		closer:          w,
		dataLog:         make(chan message, bufferSize),
		mapLogFunctions: make(map[string]logFunction),
	}
	logger.mapLogFunctions[L_INFO] = logger.log.Info
	logger.mapLogFunctions[L_DEBUG] = logger.log.Debug
	logger.mapLogFunctions[L_WARNING] = logger.log.Warn
	logger.mapLogFunctions[L_ERROR] = logger.log.Error
	go logger.start(ctx)
	return &logger
}

func (l *LoggerAsync) start(ctx context.Context) {
	for {
		select {
		case object := <-l.dataLog:
			str := strings.Join([]string{object.ThreadName, " - context:", object.Content}, "")
			l.mapLogFunctions[object.logType](str, object)
		case <-ctx.Done():
			l.mapLogFunctions[L_INFO]("", "context:", errors.New("end log"))
			fmt.Println("Logger stopped")
			time.Sleep(1 * time.Second)
			close(l.dataLog)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

type LogEngine interface {
	Info(threadName string, msg string, a ...any)
	Warn(threadName string, msg string, a ...any)
	Err(threadName string, msg string, a ...any)
	Debug(threadName string, msg string, a ...any)
	Close()
}

func (l *LoggerAsync) Close() {
	l.closer.Close()
}

func (l *LoggerAsync) Info(threadName string, msg string, a ...any) {
	l.dataLog <- message{
		logType:    L_INFO,
		ThreadName: threadName,
		Content:    msg,
		Objects:    a,
	}
}
func (l *LoggerAsync) Warn(threadName string, msg string, a ...any) {
	l.dataLog <- message{
		logType:    L_WARNING,
		ThreadName: threadName,
		Content:    msg,
		Objects:    a,
	}
}

func (l *LoggerAsync) Err(threadName string, msg string, a ...any) {
	l.dataLog <- message{
		logType:    L_ERROR,
		ThreadName: threadName,
		Content:    msg,
		Objects:    a,
	}
}

func (l *LoggerAsync) Debug(threadName string, msg string, a ...any) {
	l.dataLog <- message{
		logType:    L_DEBUG,
		ThreadName: threadName,
		Content:    msg,
		Objects:    a,
	}
}
