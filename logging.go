package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type ctxKey string

const (
	loggerKey ctxKey = "logger"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *wrappedWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func create() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	slog.SetDefault(logger)
}

func Logging(next http.Handler) http.Handler {
	create()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := uuid.New().String()

		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anon"
		}

		logger := slog.With(
			slog.String("trace_id", traceID),
			slog.String("user_id", userID),
		)

		ctx := r.Context()
		ctx = context.WithValue(ctx, loggerKey, logger)

		logger.Info("request_in",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("host", r.Host),
			slog.String("ip", r.RemoteAddr),
			slog.String("ua", r.UserAgent()),
		)

		wrapped := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r.WithContext(ctx))

		logger.Info("request_out",
			slog.Int("status_code", wrapped.statusCode),
			slog.Int("size", wrapped.size),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

func logFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
