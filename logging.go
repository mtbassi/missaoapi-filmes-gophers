package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// ==== Context Keys Tipados ====
type ctxKey string

const (
	loggerKey  ctxKey = "logger"
	traceIDKey ctxKey = "trace_id"
	userIDKey  ctxKey = "user_id"
)

// ==== Writer wrapper ====
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

// ==== Configura logger global ====
func create() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	slog.SetDefault(logger)
}

// ==== Middleware ====
func Logging(next http.Handler) http.Handler {
	create()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := uuid.New().String()

		// Recupera user-id do header ou atribui "anon"
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anon"
		}

		// Cria logger contextualizado
		logger := slog.With(
			slog.String("trace_id", traceID),
			slog.String("user_id", userID),
		)

		// Cria contexto com logger e metadados
		ctx := r.Context()
		ctx = context.WithValue(ctx, loggerKey, logger)
		ctx = context.WithValue(ctx, traceIDKey, traceID)
		ctx = context.WithValue(ctx, userIDKey, userID)

		// Loga a entrada da requisição
		logger.Info("request_in",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("host", r.Host),
			slog.String("ip", r.RemoteAddr),
			slog.String("ua", r.UserAgent()),
		)

		// Wrap da resposta
		wrapped := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Processa request
		next.ServeHTTP(wrapped, r.WithContext(ctx))

		// Loga a saída da requisição
		logger.Info("request_out",
			slog.Int("status_code", wrapped.statusCode),
			slog.Int("size", wrapped.size),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

// ==== Helpers para pegar logger no contexto ====
func logFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

func TraceIDFrom(ctx context.Context) string {
	if t, ok := ctx.Value(traceIDKey).(string); ok {
		return t
	}
	return ""
}

func UserIDFrom(ctx context.Context) string {
	if u, ok := ctx.Value(userIDKey).(string); ok {
		return u
	}
	return ""
}
