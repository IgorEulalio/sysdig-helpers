package logging

import (
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/config"
	"bytes"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

var Log *logrus.Entry

// ResponseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type ResponseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

// NewResponseWriter creates a new responseWriter.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}
}

// WriteHeader saves the status code and writes it to the underlying
// http.ResponseWriter.
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write writes the data to the body and underlying http.ResponseWriter.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func init() {

	if config.Config == nil {
		config.LoadConfig()
	}

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	serviceName := config.Config.ServiceName

	logLevel := config.Config.LogLevel

	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		logger.Warningf("Invalid or empty log level provided, defaulting to 'info'. Error: %s", err.Error())
		level = logrus.InfoLevel
	}

	logger.SetLevel(level)

	Log = logger.WithField("service", serviceName)
}

func InitLogger(config *config.Configuration) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	serviceName := config.ServiceName

	logLevel := config.LogLevel

	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		logger.Warningf("Invalid or empty log level provided, defaulting to 'info'. Error: %s", err.Error())
		level = logrus.InfoLevel
	}

	logger.SetLevel(level)

	Log = logger.WithField("service", serviceName)
}
