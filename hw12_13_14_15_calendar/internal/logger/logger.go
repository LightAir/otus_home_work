package logger

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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

func (l Logger) Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func (l Logger) Info(msg ...interface{}) {
	logrus.Info(msg...)
}

func (l Logger) Error(msg ...interface{}) {
	logrus.Error(msg...)
}

func (l Logger) LogHTTPRequest(r *http.Request, d time.Duration, statusCode int) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		logrus.Errorf("error split host and port: %q", r.RemoteAddr)
	}

	formatString := "%s [%s] %s %s %s %d %d %s"
	t := time.Now().Format("02/Jan/2006:15:04:05 -0700")
	logrus.Infof(formatString, ip, t, r.Method, r.URL.String(), r.Proto, statusCode, d.Microseconds(), r.UserAgent())
}

func (l Logger) LogGRPCRequest(ctx context.Context, info *grpc.UnaryServerInfo, d time.Duration, statusCode string) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		logrus.Errorf("error get of peer information: %s", p.Addr.String())
	}

	userAgent := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userAgent = strings.Join(md["user-agent"], " ")
	}

	ip, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		logrus.Errorf("error split host and port: %q", p.Addr.String())
	}

	formatString := "%s [%s] %s %s %s %s %d %s"
	t := time.Now().Format("02/Jan/2006:15:04:05 -0700")
	logrus.Infof(formatString, ip, t, "rpc", info.FullMethod, "HTTP/2", statusCode, d.Microseconds(), userAgent)
}
