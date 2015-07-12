package logger

import (
	"fmt"
	"redis"
	"time"
)

///log level
const (
	LOG_INFO  = iota
	LOG_DEBUG = iota
	LOG_WARN  = iota
	LOG_ERROR = iota
	LOG_FATAL = iota
)

type Formatter interface {
	Format(msg []byte) []byte
}

type defaultFormatter struct{}

func (*defaultFormatter) Format(msg []byte) []byte {
	afterFmt := fmt.Sprintf("{%v}  {%s}", time.Now(), msg)

	return []byte(afterFmt)
}

type Logger struct {
	logLev   uint8
	servAddr string
	servPort int
	client   redis.Client
	fmter    Formatter
}

func NewLogger(servAddr string, servPort int, logLev uint8) (*Logger, error) {
	spec := redis.DefaultSpec().Host(servAddr).Port(servPort)
	cli, err := redis.NewSynchClientWithSpec(spec)

	if nil != err {
		return nil, err
	}

	return &Logger{logLev, servAddr, servPort, cli, new(defaultFormatter)}, nil
}

func (log *Logger) SetFormatter(fmt Formatter) {
	log.fmter = fmt
}

func (log *Logger) Fatal(msg []byte) {
	if log.logLev > LOG_FATAL {
		return
	}
	msg = log.fmter.Format(msg)
	log.client.Lpush("FATAL", msg)
}

func (log *Logger) Error(msg []byte) {
	if log.logLev > LOG_ERROR {
		return
	}
	msg = log.fmter.Format(msg)
	log.client.Lpush("ERROR", msg)
}

func (log *Logger) Warn(msg []byte) {
	if log.logLev > LOG_WARN {
		return
	}
	msg = log.fmter.Format(msg)
	log.client.Lpush("WARN", msg)
}

func (log *Logger) Debug(msg []byte) {
	if log.logLev > LOG_DEBUG {
		return
	}
	msg = log.fmter.Format(msg)
	log.client.Lpush("DEBUG", msg)
}

func (log *Logger) Info(msg []byte) {
	if log.logLev > LOG_INFO {
		return
	}
	msg = log.fmter.Format(msg)
	log.client.Lpush("INFO", msg)
}

func (log *Logger) Log(level uint8, msg []byte) {
	if log.logLev > level {
		return
	}
	switch level {
	case LOG_INFO:
		log.Info(msg)
	case LOG_WARN:
		log.Warn(msg)
	case LOG_ERROR:
		log.Error(msg)
	case LOG_FATAL:
		log.Fatal(msg)
	case LOG_DEBUG:
		log.Debug(msg)
	default:
		return
	}
}
