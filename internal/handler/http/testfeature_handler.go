package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sazardev/goca/internal/usecase"
)

type TestFeatureHandler struct {
	usecase usecase.TestFeatureUseCase
}

func NewTestFeatureHandler(uc usecase.TestFeatureUseCase) *TestFeatureHandler {
	return &TestFeatureHandler{usecase: uc}
}

func (t *TestFeatureHandler) CreateTestFeature(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateTestFeatureInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	output, err := t.usecase.Create(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (t *TestFeatureHandler) GetTestFeature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid testfeature ID", http.StatusBadRequest)
		return
	}

	testfeature, err := t.usecase.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testfeature)
}

func (t *TestFeatureHandler) UpdateTestFeature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid testfeature ID", http.StatusBadRequest)
		return
	}

	var input usecase.UpdateTestFeatureInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if _, err := t.usecase.Update(uint(id), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *TestFeatureHandler) DeleteTestFeature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid testfeature ID", http.StatusBadRequest)
		return
	}

	if err := t.usecase.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *TestFeatureHandler) ListTestFeatures(w http.ResponseWriter, r *http.Request) {
	output, err := t.usecase.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
