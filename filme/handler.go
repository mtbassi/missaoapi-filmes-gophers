package filme

import (
	"encoding/json"
	"log/slog"
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
	logger := r.Context().Value("logger").(*slog.Logger)

	id := r.PathValue("id")
	filme, exists := h.Repo.Dados[id]

	w.Header().Set("Content-Type", "application/json")
	if !exists {
		logger.With(slog.String("evento", "listar_pelo_id")).
			Info("Nenhum filme encontrado com o id " + id)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	logger.With(slog.String("evento", "listar_pelo_id")).
		Info("Filme " + filme.Filme + " encontrado pelo id " + filme.ID.String())
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filme)
}

func (h *Handler) Cadastrar(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*slog.Logger)

	var novoFilme FilmePost
	errJson := json.NewDecoder(r.Body).Decode(&novoFilme)
	if errJson != nil {
		logger.Error("Falha na desserialização de JSON")
		retornarErroCampo(w, "json", "JSON inválido")
		return
	}

	validate := validator.New()
	err := validate.Struct(novoFilme)
	if err != nil {
		logger.Error("Falha na validação de campos")
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

	logger.With(slog.String("evento", "cadastro_filme")).
		Info("Filme cadastrado com sucesso")

	w.Header().Set("Location", "/filmes/"+filme.ID.String())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(filme)
}

func (h *Handler) Deletar(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*slog.Logger)

	delete(h.Repo.Dados, r.PathValue("id"))

	logger.With(slog.String("evento", "deletar")).
		Info("Operação realizada com sucesso")
	w.WriteHeader(http.StatusNoContent)
}
