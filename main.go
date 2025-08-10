package main

import (
	"encoding/json"
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
	Dados map[string]Filme
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

	return &RepositorioFilmes{Dados: filmes}
}

func (repo *RepositorioFilmes) Listar(w http.ResponseWriter, r *http.Request) {
	filmes := make([]Filme, 0, len(repo.Dados))
	for _, filme := range repo.Dados {
		filmes = append(filmes, filme)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filmes)
}

func main() {
	repo := criarRepo()

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
