package handlers

import (
	"encoding/json"
	"fmt"
	"forum/models"
	"forum/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, userService *services.UserService) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	err := r.ParseForm()
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, "Unable to parse form")
		return
	}

	userName := r.FormValue("user_name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	gender := r.FormValue("gender")
	ageStr := r.FormValue("age")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, "Invalid age")
		return
	}

	userUUID, err := uuid.NewV4()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Failed to generate UUID")
		return
	}

	_, err = userService.InsertUser(userName, email, password, firstName, lastName, gender, age, userUUID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to register user: %v", err))
		return
	}

	user := &models.User{
		UserId:    userUUID,
		Username:  userName,
		Email:     email,
		Password:  password,
		Firstname: firstName,
		Lastname:  lastName,
		Age:       age,
		Gender:    gender,
		CreatedAt: time.Now(),
	}

	fmt.Println(user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, authService *services.AuthService) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid request method",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Unable to parse form",
			"status":  http.StatusBadRequest,
		})
		return
	}

	// Ensure the field sent corresponds to "identifier"
	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	userId, token, err := authService.Login(identifier, password)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid identifier or password",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	response := &models.LoginResponse{
		Response: &models.ResponseBody{
			Message: "Login successful",
			Status:  http.StatusOK,
		},
		UserId: userId,
		Token:  token,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logout logic here
}
