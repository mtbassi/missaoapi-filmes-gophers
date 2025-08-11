package filme

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Repo *RepositorioFilmes
}

func NewHandler() *Handler {
	repo := criarRepo()
	return &Handler{
		Repo: repo,
	}
}

func (h *Handler) Listar(w http.ResponseWriter, r *http.Request) {
	filmes := make([]Filme, 0, len(h.Repo.Dados))
	for _, filme := range h.Repo.Dados {
		filmes = append(filmes, filme)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filmes)
}

func (h *Handler) ListarPeloId(w http.ResponseWriter, r *http.Request) {
	filme, exists := h.Repo.Dados[r.PathValue("id")]

	w.Header().Set("Content-Type", "application/json")
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filme)
}
