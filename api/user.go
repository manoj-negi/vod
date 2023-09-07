package api

import (
	"encoding/json"
	"net/http"

	db "github.com/vod/db/sqlc"
	util "github.com/vod/utils"
)

type User struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JsonResponse struct {
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"status_code"`
}

type SuceessResponse struct {
	Status    bool        `json:"status"`
	Message   string      `json:"message"`
	UserToken interface{} `json:"userToken"`
}

func (server *Server) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}
	ctx := r.Context()

	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {

		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "invalid JSON request",
			StatusCode: http.StatusTeapot,
		}
		util.WriteJSONResponse(w, http.StatusTeapot, jsonResponse)
		return
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {

		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "invalid JSON request",
			StatusCode: http.StatusTeapot,
		}
		util.WriteJSONResponse(w, http.StatusTeapot, jsonResponse)
		return
	}

	arg := db.CreateUserParams{
		Username:       user.Username,
		HashedPassword: hashedPassword,
		FullName:       user.FullName,
		Email:          user.Email,
	}

	userInfo, err := server.store.CreateUser(ctx, arg)
	if err != nil {

		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "invalid JSON request",
			StatusCode: http.StatusTeapot,
		}
		util.WriteJSONResponse(w, http.StatusTeapot, jsonResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

func (server *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}

	// dummy JSON response
	userInfo := struct {
		Message string `json:"message"`
		UserID  int    `json:"user_id"`
	}{
		Message: "Login successful",
		UserID:  123, // Replace with an actual user ID if needed
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

func (server *Server) refreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}

	// dummy JSON response
	userInfo := struct {
		Message string `json:"message"`
		UserID  int    `json:"user_id"`
	}{
		Message: "refresh token",
		UserID:  123, // Replace with an actual user ID if needed
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}
