package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	if err := validate.RegisterValidation("yearValid", yearValid); err != nil {
		panic(fmt.Sprintf("failed to register yearValid validation: %v", err))
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	var in MovieCreate
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := validate.Struct(in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	m, err := store.create(r.Context(), in)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, m)
}

func handleListMovies(w http.ResponseWriter, r *http.Request) {
	ms := store.list()
	writeJSON(w, http.StatusOK, ms)
}

func handleGetMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	m, ok := store.get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func handlePatchMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	var in MoviePatch
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := validate.Struct(in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	m, err := store.patch(r.Context(), id, in)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	ok := store.delete(r.Context(), id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
