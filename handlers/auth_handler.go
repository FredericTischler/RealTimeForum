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
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, "Unable to parse multipart form")
		return
	}

	// Ensure the field sent corresponds to "identifier"
	identifier := r.FormValue("identifier")
	fmt.Println("Identifiant : " + identifier)
	password := r.FormValue("password")

	userId, token, err := authService.Login(identifier, password)

	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, "Invalid identifier or password")
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
		ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logout logic here
}
