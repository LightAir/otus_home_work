package logger

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct{}

func New(level string) *Logger {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		log.Fatalf("failed to parse the level: %v", err)
	}

	logrus.SetLevel(logLevel)

	return &Logger{}
}

func (l Logger) Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func (l Logger) Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func (l Logger) Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func (l Logger) Debug(msg ...interface{}) {
	logrus.Debug(msg...)
}

func (l Logger) Info(msg ...interface{}) {
	logrus.Info(msg...)
}

func (l Logger) Warn(msg ...interface{}) {
	logrus.Warn(msg...)
}

func (l Logger) Error(msg ...interface{}) {
	logrus.Error(msg...)
}

func (l Logger) LogHTTPRequest(r *http.Request, d time.Duration, statusCode int) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		logrus.Errorf("error split host and port: %q", r.RemoteAddr)
	}

	if net.ParseIP(ip) == nil {
		logrus.Errorf("error parse ip: %q", r.RemoteAddr)
	}

	formatString := "%s [%s] %s %s %s %d %d %s"
	t := time.Now().Format("02/Jan/2006:15:04:05 -0700")
	logrus.Infof(formatString, ip, t, r.Method, r.URL.String(), r.Proto, statusCode, d.Microseconds(), r.UserAgent())
}
