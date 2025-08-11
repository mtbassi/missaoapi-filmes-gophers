package filme

import (
	"encoding/json"
	"net/http"
)

func (rf *RepositorioFilmes) Listar(w http.ResponseWriter, r *http.Request) {
	filmes := make([]Filme, 0, len(rf.Dados))
	for _, filme := range rf.Dados {
		filmes = append(filmes, filme)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filmes)
}
