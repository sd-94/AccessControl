package main

import (
	"context"
	"log"
	"net/http"
	"prac/config"
	"prac/databases"
	"prac/handlers"
	"prac/repositories"
	"prac/types"
	"prac/utils"
	"slices"

	"github.com/gorilla/mux"
	"github.com/joomcode/errorx"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	database, err := databases.GetDB(cfg.DBConfig)
	if err != nil {
		if errorx.IsOfType(err, types.ConnectionError) {
			log.Println("Error: Failed to connect to the database:", err)
			return
		}

		if errorx.IsOfType(err, types.TableCreationError) {
			log.Println("Error: Failed to create a table:", err)
			return
		}

		log.Println("Unknown error:", err)
	}
	defer database.Close()

	repository := repositories.NewRepository(database)
	handler := handlers.NewHandler(cfg, repository)

	router := mux.NewRouter()
	router.Use(AuthMiddleWareWithConfig(cfg))
	router.HandleFunc("/auth/token", handlers.HandleError(handler.HandleGetToken)).Methods("POST")

	router.HandleFunc("/accounts/{acc_id}", handlers.HandleError(handler.HandleGetAccount)).Methods("GET")
	router.HandleFunc("/accounts", handlers.HandleError(handler.HandleGetAccounts)).Methods("GET")
	router.HandleFunc("/accounts", handlers.HandleError(handler.HandleCreateAccount)).Methods("POST")
	router.HandleFunc("/accounts/{acc_id}", handlers.HandleError(handler.HandleUpdateAccount)).Methods("PUT")
	router.HandleFunc("/accounts/{acc_id}", handlers.HandleError(handler.HandleDeleteAccount)).Methods("DELETE")

	log.Println("server is up and running")
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}

func AuthMiddleWareWithConfig(config *config.Config) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if slices.Contains(config.Auth.Ignore, path) {
				h.ServeHTTP(w, r)
				return
			}
			token := r.Header.Get(config.Auth.Header)
			if token == "" {
				log.Println("permission denied. Authorization token is absent.")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Permission denied"))
				return
			}
			email, err := utils.ValidateToken(config.Auth.Secret, token)
			if err != nil {
				log.Println("permission denied. Token is invalid.")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Permission denied"))
				return
			}

			context := context.WithValue(r.Context(), "email", email)
			h.ServeHTTP(w, r.WithContext(context))
		})
	}
}
