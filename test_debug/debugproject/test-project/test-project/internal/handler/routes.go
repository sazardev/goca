package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/test/debugproject/internal/usecase"
)

func SetupTestUserFeatureRoutes(router *mux.Router, uc usecase.TestUserFeatureUseCase) {
	handler := NewTestUserFeatureHandler(uc)

	// Apply middleware
	testuserfeatureRouter := router.PathPrefix("/testuserfeatures").Subrouter()
	testuserfeatureRouter.Use(corsMiddleware)
	testuserfeatureRouter.Use(loggingMiddleware)

	testuserfeatureRouter.HandleFunc("", handler.CreateTestUserFeature).Methods("POST")
	testuserfeatureRouter.HandleFunc("/{id}", handler.GetTestUserFeature).Methods("GET")
	testuserfeatureRouter.HandleFunc("/{id}", handler.UpdateTestUserFeature).Methods("PUT")
	testuserfeatureRouter.HandleFunc("/{id}", handler.DeleteTestUserFeature).Methods("DELETE")
	testuserfeatureRouter.HandleFunc("", handler.ListTestUserFeatures).Methods("GET")
}

// Middleware functions

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
