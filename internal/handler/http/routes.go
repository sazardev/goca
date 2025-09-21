package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sazardev/goca/internal/usecase"
)

func SetupTestPerfRoutes(router *mux.Router, uc usecase.TestPerfUseCase) {
	handler := NewTestPerfHandler(uc)

	// Apply middleware
	testperfRouter := router.PathPrefix("/testperfs").Subrouter()
	testperfRouter.Use(corsMiddleware)
	testperfRouter.Use(loggingMiddleware)

	testperfRouter.HandleFunc("", handler.CreateTestPerf).Methods("POST")
	testperfRouter.HandleFunc("/{id}", handler.GetTestPerf).Methods("GET")
	testperfRouter.HandleFunc("/{id}", handler.UpdateTestPerf).Methods("PUT")
	testperfRouter.HandleFunc("/{id}", handler.DeleteTestPerf).Methods("DELETE")
	testperfRouter.HandleFunc("", handler.ListTestPerfs).Methods("GET")
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
