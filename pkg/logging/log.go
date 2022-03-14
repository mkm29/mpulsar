package logging

import (
	"net"
	"net/http"

	logrus "github.com/sirupsen/logrus"

	"github.com/mkm29/mpulsar/pkg/utils"
)

var logger = logrus.New()

func Configure() {
	// set up logging output
	logger.SetFormatter(&logrus.JSONFormatter{})
	// set logging level (change to higher level when in production)
	logger.SetLevel(logrus.TraceLevel)
}

func Log(l string, args ...interface{}) {
	levels := map[string]logrus.Level{
		"DEBUG": logrus.TraceLevel,
		"INFO":  logrus.InfoLevel,
		"WARN":  logrus.WarnLevel,
		"ERROR": logrus.ErrorLevel,
		"FATAL": logrus.FatalLevel,
		"TRACE": logrus.TraceLevel,
	}
	// Check the log level from the environment variable
	LOGLEVEL := utils.GetEnv("LOGLEVEL", "ERROR")
	// if LOGLEVEL from environment is less than the level return and do not log
	if levels[l] < levels[LOGLEVEL] {
		return
	}
	logger.Log(levels[l], args...)
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
