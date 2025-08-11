package main

import (
	"net/http"

	"github.com/mtbassi/missaoapi-filmes-gophers/filme"
	"github.com/mtbassi/missaoapi-filmes-gophers/middleware"
)

func main() {
	handler := filme.NewHandler()

	router := http.NewServeMux()
	router.HandleFunc("GET /filmes", handler.Listar)
	router.HandleFunc("GET /filmes/{id}", handler.ListarPeloId)
	router.HandleFunc("POST /filmes", handler.Cadastrar)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	server := http.Server{
		Addr:    ":8081",
		Handler: middleware.Logging(v1),
	}

	server.ListenAndServe()
}
