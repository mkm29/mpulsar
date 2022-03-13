package log

import (
	"net"
	"net/http"

	logrus "github.com/sirupsen/logrus"
)

var logger = logrus.New()

func configure() {
	// set up logging output
	logger.SetFormatter(&logrus.JSONFormatter{})
	// set logging level (change to higher level when in production)
	logger.SetLevel(logrus.DebugLevel)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func WithConn(conn net.Conn) *logrus.Entry {
	var addr string = "unknown"
	if conn != nil {
		addr = conn.RemoteAddr().String()
	}
	return logger.WithField("addr", addr)
}

// extract fields from a request
func RequestFields(req *http.Request) logrus.Fields {
	return logrus.Fields{
		"method": req.Method,
		"path":   req.URL.Path,
		"proto":  req.Proto,
		"host":   req.Host,
		// IP address of the client
		"remote": req.RemoteAddr,
	}
}

func WithRequest(req *http.Request) *logrus.Entry {
	return logger.WithFields(RequestFields(req))
}
