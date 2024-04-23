package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social/database"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type pamameters struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := pamameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Failed to decode %v", err)
		return
	}

	user, err := apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Email:     params.Email,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldnt create user::%v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, databaseUserToUser(user))

}

func (apiConfig *apiConfig) handleGetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := apiConfig.DB.GetUsers(r.Context())
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldnt create user::%v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, databaseUserToUsers(users))

}

func (apiConfig *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)

	if err != nil {
		log.Printf("Failed to parse UUID: %v", err)
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	users, err := apiConfig.DB.GetUser(r.Context(), id)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldnt create user::%v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, databaseUserToUser(users))

}

func (apiConfig *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Failed to parse UUID: %v", err)
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	type parameter struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Failed to decode updating users: %v", err)
		http.Error(w, "Invaid Updating User", http.StatusBadRequest)
		return
	}

	user, err := apiConfig.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:        id,
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Email:     params.Email,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldnt create user::%v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, databaseUserToUser(user))

}

func (apiConfig *apiConfig) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Failed to parse UUID: %v", err)
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	err = apiConfig.DB.DeleteUser(r.Context(), id)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldnt create user::%v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, struct{}{})

}
