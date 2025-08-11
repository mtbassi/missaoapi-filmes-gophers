package filme

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
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

func (h *Handler) Cadastrar(w http.ResponseWriter, r *http.Request) {
	var novoFilme FilmePost
	errJson := json.NewDecoder(r.Body).Decode(&novoFilme)
	if errJson != nil {
		retornarErroCampo(w, "json", "JSON inválido")
		return
	}

	validate := validator.New()
	err := validate.Struct(novoFilme)
	if err != nil {
		mapearCamposErros(w, err)
		return
	}

	filme := Filme{
		ID:         uuid.New(),
		Filme:      novoFilme.Filme,
		Duracao:    novoFilme.Duracao,
		Categorias: novoFilme.Categorias,
		CriadoEm:   time.Now(),
	}
	h.Repo.Dados[filme.ID.String()] = filme

	w.Header().Set("Location", "/filmes/"+filme.ID.String())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(filme)
}
