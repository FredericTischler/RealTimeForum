package handlers

import (
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
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
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
		http.Error(w, "Invalid age", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "Failed to generate UUID", http.StatusInternalServerError)
		return
	}

	db, err := services.ConnectDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to database: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	userService := services.UserMethod{DB: db}

	_, err = userService.InsertInUser(userName, email, password, firstName, lastName, gender, age, userUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)
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

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

}
