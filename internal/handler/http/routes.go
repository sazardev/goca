package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sazardev/goca/internal/usecase"
)

func SetupOrderRoutes(router *mux.Router, uc usecase.OrderUseCase) {
	handler := NewOrderHandler(uc)

	// Apply middleware
	orderRouter := router.PathPrefix("/orders").Subrouter()
	orderRouter.Use(corsMiddleware)
	orderRouter.Use(loggingMiddleware)

	orderRouter.HandleFunc("", handler.CreateOrder).Methods("POST")
	orderRouter.HandleFunc("/{id}", handler.GetOrder).Methods("GET")
	orderRouter.HandleFunc("/{id}", handler.UpdateOrder).Methods("PUT")
	orderRouter.HandleFunc("/{id}", handler.DeleteOrder).Methods("DELETE")
	orderRouter.HandleFunc("", handler.ListOrders).Methods("GET")
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
