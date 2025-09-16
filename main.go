package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

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
		Addr:              fmt.Sprintf(":%d", config.Server.Port),
		Handler:           Logging(v1),
		ReadTimeout:       config.Server.ReadTimeout,
		ReadHeaderTimeout: config.Server.ReadHeaderTimeout,
		WriteTimeout:      config.Server.WriteTimeout,
		IdleTimeout:       config.Server.IdleTimeout,
	}

	slog.InfoContext(context.Background(), "starting server", "app", config.App.Name, "version", config.App.Version)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
