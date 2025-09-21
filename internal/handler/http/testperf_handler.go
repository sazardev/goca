package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sazardev/goca/internal/messages"
	"github.com/sazardev/goca/internal/usecase"
)

type TestPerfHandler struct {
	usecase usecase.TestPerfUseCase
}

func NewTestPerfHandler(uc usecase.TestPerfUseCase) *TestPerfHandler {
	return &TestPerfHandler{usecase: uc}
}

func (t *TestPerfHandler) CreateTestPerf(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateTestPerfInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	output, err := t.usecase.CreateTestPerf(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (t *TestPerfHandler) GetTestPerf(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid testperf ID", http.StatusBadRequest)
		return
	}

	testperf, err := t.usecase.GetTestPerf(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testperf)
}

func (t *TestPerfHandler) UpdateTestPerf(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid testperf ID", http.StatusBadRequest)
		return
	}

	var input usecase.UpdateTestPerfInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := t.usecase.UpdateTestPerf(id, input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *TestPerfHandler) DeleteTestPerf(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid testperf ID", http.StatusBadRequest)
		return
	}

	if err := t.usecase.DeleteTestPerf(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *TestPerfHandler) ListTestPerfs(w http.ResponseWriter, r *http.Request) {
	output, err := t.usecase.ListTestPerfs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
