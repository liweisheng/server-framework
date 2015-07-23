/**
 * 向redis服务器写日志.
 */

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

/// 在日志写向服务器之前，格式化日志.
///
/// 实现该接口可以自定义格式化器.
type Formatter interface {
	Format(msg []byte) []byte
}

/// 默认的日志格式化器
type defaultFormatter struct{}

func (*defaultFormatter) Format(msg []byte) []byte {
	afterFmt := fmt.Sprintf("{%v}  {%s}", time.Now(), msg)

	return []byte(afterFmt)
}

type Logger struct {
	logLev   uint8  /// 限制日志输出级别，大于等于这个级别可以输出.
	servAddr string /// redis服务器地址
	servPort int    /// redis服务器端口
	client   redis.Client
	fmter    Formatter /// 格式化器
}

/// 创建新的日志记录器.
///
/// @param servAddr redis服务器地址
/// @param servPort redis服务器端口
/// @param logLev 限制日志输出级别
/// @return 成功则返回日志记录器，同时error为nil.
func NewLogger(servAddr string, servPort int, logLev uint8) (*Logger, error) {
	spec := redis.DefaultSpec().Host(servAddr).Port(servPort)
	cli, err := redis.NewSynchClientWithSpec(spec)

	if nil != err {
		return nil, err
	}

	return &Logger{logLev, servAddr, servPort, cli, new(defaultFormatter)}, nil
}

/// 设置格式化器.
func (log *Logger) SetFormatter(fmt Formatter) {
	log.fmter = fmt
}

/// 写如FATAL级别日志.
func (log *Logger) Fatal(msg []byte) {
	if log.logLev > LOG_FATAL {
		return
	}
	msg = log.fmter.Format(msg)
	log.client.Lpush("FATAL", msg)
}

/// ERROR级别日志.
func (log *Logger) Error(msg []byte) {
	if log.logLev > LOG_ERROR {
		return
	}
	msg = log.fmter.Format(msg)
	log.client.Lpush("ERROR", msg)
}

/// WARN级别日志.
func (log *Logger) Warn(msg []byte) {
	if log.logLev > LOG_WARN {
		return
	}
	msg = log.fmter.Format(msg)
	log.client.Lpush("WARN", msg)
}

/// DEBUG级别日志.
func (log *Logger) Debug(msg []byte) {
	if log.logLev > LOG_DEBUG {
		return
	}
	msg = log.fmter.Format(msg)
	log.client.Lpush("DEBUG", msg)
}

/// INFO级别.
func (log *Logger) Info(msg []byte) {
	if log.logLev > LOG_INFO {
		return
	}
	msg = log.fmter.Format(msg)
	log.client.Lpush("INFO", msg)
}

/// 写入由参数level指定的级别的日志.
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
