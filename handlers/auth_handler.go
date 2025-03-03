package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	db, err := services.ConnectDB()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to connect to database: %v", err))
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			// TODO LOG BDD ERRORS
		}
	}(db)

	userService := services.UserMethod{DB: db}

	_, err = userService.InsertInUser(userName, email, password, firstName, lastName, gender, age, userUUID)
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Implement login logic here
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logout logic here
}
