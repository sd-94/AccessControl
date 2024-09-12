package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"prac/config"
	"prac/repositories"
	"prac/types"
	"prac/utils"

	"github.com/joomcode/errorx"
)

type Handler struct {
	Config     *config.Config
	Repository *repositories.Repository
}

type HandlerWithErrorFunc func(w http.ResponseWriter, r *http.Request) error

func NewHandler(cfg *config.Config, repository *repositories.Repository) *Handler {
	return &Handler{
		Config:     cfg,
		Repository: repository,
	}
}

func HandleError(handler HandlerWithErrorFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			if errorx.IsOfType(err, types.ValidationError) {
				log.Println("Validation failed:", err.Error())
				writeErrorResponse(w, http.StatusBadRequest, utils.WriteJson(map[string]string{"error": err.Error()}))
				return
			}

			if errorx.IsOfType(err, types.NotFoundError) {
				log.Println("Resource not found:", err.Error())
				writeErrorResponse(w, http.StatusNotFound, utils.WriteJson(map[string]string{"error": err.Error()}))
				return
			}

			if errorx.IsOfType(err, types.AuthorizationError) {
				log.Println("Authorization failed:", err.Error())
				writeErrorResponse(w, http.StatusUnauthorized, utils.WriteJson(map[string]string{"error": err.Error()}))
				return
			}

			log.Println("Internal server error:", err.Error())
			writeErrorResponse(w, http.StatusInternalServerError, utils.WriteJson(map[string]string{"error": "Internal server error"}))
		}
	}
}

func writeErrorResponse(w http.ResponseWriter, httpCode int, msg []byte) {
	w.WriteHeader(httpCode)
	w.Write(msg)
}

func (handler *Handler) HandleGetToken(w http.ResponseWriter, r *http.Request) error {
	var body types.SignIn

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return types.JSONDecodingError.Wrap(err, "failed to decode JSON body")
	}

	if body.Email == "" || body.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return types.ValidationError.New("email or password is missing")
	}

	exists, err := handler.Repository.SignIn(&body)
	if err != nil {
		log.Println("Error finding account:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return types.SQLExecutionError.Wrap(err, "failed to find account")
	}

	if exists {
		token, err := utils.GenerateToken(handler.Config.Auth.Secret, body.Email)
		if err != nil {
			log.Println("Error generating token:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return types.TokenGenerationError.Wrap(err, "failed to generate token")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return types.AuthorizationError.New("invalid email or password")
	}

	return nil
}
