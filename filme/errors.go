package filme

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type ErroCampo struct {
	Campo string `json:"campo"`
	Erro  string `json:"erro"`
}

func retornarErroCampo(w http.ResponseWriter, campo, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErroCampo{
		Campo: campo,
		Erro:  msg,
	})
}

func mapearCamposErros(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var erros []ErroCampo

	for _, e := range err.(validator.ValidationErrors) {
		erros = append(erros, ErroCampo{
			Campo: e.StructField(),
			Erro:  "campo obrigatório ou inválido",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(erros)
}
