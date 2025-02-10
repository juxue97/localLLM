package user

import (
	"fmt"
	"net/http"
	"time"

	"chatbot/cmd/service/auth"
	"chatbot/config"
	"chatbot/types"
	"chatbot/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Define Payload Structure
	var payload types.LoginUserPayload
	// Get JSON payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// Validate Request JSON Body
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors)
		return
	}

	// Check if the email is registered
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	// Check if the password is correct
	if !auth.ComparePassword(u.Credential.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("not found, invalid email or password"))
		return
	}

	// Provide jwt access token
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID.Hex())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	responseFormat := types.APIResponse{
		Success: true,
		Message: "Login Successful",
		Data:    map[string]string{"accessToken": token},
	}
	utils.WriteJSON(w, http.StatusOK, responseFormat)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload

	// Get JSON payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// Validate Request JSON Body
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid paylod %v", errors))
		return
	}

	// check if user already exist
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exist", payload.Email))
		return
	}

	// hash the password before insert payload into the database
	hashPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Get IP address of the user
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}

	userID, err := h.store.CreateUser(types.User{
		Credential: types.Credential{
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Email:     payload.Email,
			Password:  hashPassword,
		},
		Roles:         "user",
		CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
		LastLoginAt:   primitive.NewDateTimeFromTime(time.Time{}), // Set LastLoginAt to null
		LastSessionIP: ip,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	responseFormat := types.APIResponse{
		Success: true,
		Message: "User Successfully Register",
		Data: types.User{
			ID: userID,
			Credential: types.Credential{
				FirstName: payload.FirstName,
				LastName:  payload.LastName,
				Email:     payload.Email,
				Password:  hashPassword,
			},
			Roles:     "user",
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	utils.WriteJSON(w, http.StatusCreated, responseFormat)
}
