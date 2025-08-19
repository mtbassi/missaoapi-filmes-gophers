package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type movieStore struct {
	mu     sync.RWMutex
	autoID string
	items  map[string]Movie
}

func newMovieStore() *movieStore {
	return &movieStore{items: make(map[string]Movie)}
}

func (s *movieStore) create(ctx context.Context, in MovieCreate) (Movie, error) {
	logger := logFromContext(ctx).With(slog.String("event", "create_movie"))
	if strings.TrimSpace(in.Name) == "" || in.Duration <= 0 {
		m := "invalid movie: name required and duration > 0"
		logger.Error(
			m,
			slog.String("status", "error"),
		)
		return Movie{}, errors.New(m)
	}
	m := Movie{
		ID:         uuid.New().String(),
		Name:       strings.TrimSpace(in.Name),
		Duration:   in.Duration,
		Categories: in.Categories,
		CreatedAt:  time.Now()}
	s.items[m.ID] = m
	logger.Info(
		"movie created",
		slog.String("status", "created"),
		slog.String("id", m.ID),
	)
	return m, nil
}

func (s *movieStore) get(id string) (Movie, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.items[id]
	return m, ok
}

func (s *movieStore) list() []Movie {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]Movie, 0, len(s.items))
	for _, m := range s.items {
		res = append(res, m)
	}
	return res
}

func (s *movieStore) patch(ctx context.Context, id string, in MoviePatch) (Movie, error) {
	logger := logFromContext(ctx).With(slog.String("event", "update_movie"))
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.items[id]
	if !ok {
		m := fmt.Sprintf("movie %s not found", id)
		logger.Error(
			m,
			slog.String("status", "error"),
		)
		return Movie{}, fmt.Errorf(m)
	}
	if in.Name != nil {
		m.Name = strings.TrimSpace(*in.Name)
	}
	if in.Duration != nil {
		m.Duration = *in.Duration
	}
	if in.Categories != nil {
		m.Categories = *in.Categories
	}
	s.items[id] = m
	logger.Info(
		"movie updated",
		slog.String("status", "updated"),
		slog.String("id", m.ID),
	)
	return m, nil
}

func (s *movieStore) delete(ctx context.Context, id string) bool {
	logger := logFromContext(ctx).With(slog.String("event", "delete_movie"))
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		m := fmt.Sprintf("movie %s not found", id)
		logger.Warn(
			m,
			slog.String("status", "warning"),
		)
		return true
	}
	delete(s.items, id)
	logger.Info(
		"movie deleted",
		slog.String("status", "deleted"),
	)
	return true
}

var store = newMovieStore()
