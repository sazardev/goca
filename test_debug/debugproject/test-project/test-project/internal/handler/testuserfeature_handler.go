package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/test/debugproject/internal/usecase"
	// TODO: Add validation imports when needed
)

type TestUserFeatureHandler struct {
	usecase usecase.TestUserFeatureUseCase
}

func NewTestUserFeatureHandler(uc usecase.TestUserFeatureUseCase) *TestUserFeatureHandler {
	return &TestUserFeatureHandler{usecase: uc}
}

func (t *TestUserFeatureHandler) CreateTestUserFeature(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateTestUserFeatureInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	output, err := t.usecase.CreateTestUserFeature(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (t *TestUserFeatureHandler) GetTestUserFeature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid testuserfeature ID", http.StatusBadRequest)
		return
	}

	testuserfeature, err := t.usecase.GetTestUserFeature(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testuserfeature)
}

func (t *TestUserFeatureHandler) UpdateTestUserFeature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid testuserfeature ID", http.StatusBadRequest)
		return
	}

	var input usecase.UpdateTestUserFeatureInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := t.usecase.UpdateTestUserFeature(id, input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *TestUserFeatureHandler) DeleteTestUserFeature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid testuserfeature ID", http.StatusBadRequest)
		return
	}

	if err := t.usecase.DeleteTestUserFeature(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *TestUserFeatureHandler) ListTestUserFeatures(w http.ResponseWriter, r *http.Request) {
	output, err := t.usecase.ListTestUserFeatures()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
