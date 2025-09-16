package main

import (
	"context"
	"log/slog"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /health", healthCheck)

	router.HandleFunc("GET /movies", handleListMovies)
	router.HandleFunc("GET /movies/{id}", handleGetMovie)
	router.HandleFunc("POST /movies", handleCreateMovie)
	router.HandleFunc("DELETE /movies/{id}", handleDeleteMovie)
	router.HandleFunc("PATCH /movies/{id}", handlePatchMovie)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	server := http.Server{
		Addr:    ":8081",
		Handler: Logging(v1),
	}

	slog.Log(context.Background(), slog.LevelInfo, "starting server")
	server.ListenAndServe()
}
