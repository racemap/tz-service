package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func InitLogger() *logrus.Logger {
	var logger = logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	return logger
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

func (lw *loggingResponseWriter) Write(b []byte) (int, error) {
	return lw.ResponseWriter.Write(b)
}

func BuildMiddleware(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entry := logrus.NewEntry(logger)
			start := time.Now()

			if reqID := r.Header.Get("X-Request-Id"); reqID != "" {
				entry = entry.WithField("requestId", reqID)
			}

			lw := newLoggingResponseWriter(w)
			next.ServeHTTP(lw, r)

			latency := time.Since(start)

			entry.WithFields(logrus.Fields{
				"status": lw.statusCode,
				"took":   latency,
			}).Info("completed handling request")
		})
	}
}
