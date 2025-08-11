package filme

import (
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

func CriarRepo() *RepositorioFilmes {
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
