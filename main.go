package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Filme struct {
	ID         uuid.UUID `json:"id"`
	Filme      string    `json:"filme"`
	Duracao    int       `json:"duracao"`
	Categorias []string  `json:"categorias"`
	CriadoEm   time.Time `json:"criadoEm"`
}

type RepositorioFilmes struct {
	dados map[string]Filme
}

func criarRepo() *RepositorioFilmes {
	filmes := map[string]Filme{}

	addFilme := func(nome string, duracao int, categorias []string) {
		id := uuid.New()
		filmes[id.String()] = Filme{
			ID:         id,
			Filme:      nome,
			Duracao:    duracao,
			Categorias: categorias,
			CriadoEm:   time.Now(),
		}
	}

	addFilme("Mad Max: Estrada da Fúria", 120, []string{"Ação", "Ficção científica"})
	addFilme("Django Livre", 165, []string{"Ação", "Faroeste"})

	return &RepositorioFilmes{dados: filmes}
}

func main() {
	router := http.NewServeMux()

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	server.ListenAndServe()
}
