package api

import (
	response "ledger/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func InitialiseRoutes() http.Handler {
	route := chi.NewRouter()
	route.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-API-KEY", "X-Api-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	route.Get("/", HealthCheck)

	route.Get("/balance", GetBalanceHandler)
	route.Post("/balance", CreateAccount)
	route.Post("/balance/add", AddAmountHandler)
	route.Post("/balance/deduct", DeductAmountHandler)
	return route
}

func HealthCheck(res http.ResponseWriter, req *http.Request) {
	response.RespondWithJSON(res, http.StatusOK, "I am working fine :)")
}
