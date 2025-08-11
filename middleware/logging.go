package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func CriarLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		traceID := uuid.New().String()
		logger := slog.With(
			slog.String("trace_id", traceID),
		)

		ctx := r.Context()
		ctx = context.WithValue(ctx, "logger", logger)
		next.ServeHTTP(wrapped, r.WithContext(ctx))

		logger.LogAttrs(
			context.Background(),
			slog.LevelInfo,
			"Info response",
			slog.Group(
				"request",
				slog.Int("status_code", wrapped.statusCode),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("host", r.Host),
				slog.String("ip", r.RemoteAddr),
				slog.String("ua", r.UserAgent()),
				slog.Duration("duration", time.Since(start)),
			),
		)
	})
}
