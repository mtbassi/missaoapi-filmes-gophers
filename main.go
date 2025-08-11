package main

import (
	"net/http"

	"github.com/mtbassi/missaoapi-filmes-gophers/filme"
)

func main() {
	repo := filme.CriarRepo()

	router := http.NewServeMux()
	router.HandleFunc("GET /filmes", repo.Listar)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	server := http.Server{
		Addr:    ":8081",
		Handler: v1,
	}

	server.ListenAndServe()
}
