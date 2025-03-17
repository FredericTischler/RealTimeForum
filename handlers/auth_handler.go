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

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, "Unable to parse multipart form")
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

	_ = &models.User{
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

	// Récupération de l'identifiant et du mot de passe depuis le formulaire
	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	// Authentification : récupération de l'UUID de l'utilisateur et du token de session
	userId, token, err := authService.Login(identifier, password)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, "Invalid identifier or password")
		return
	}

	// Création du cookie contenant le token de session
	cookie := &http.Cookie{
		Name:     "session_token", // nom du cookie
		Value:    token,           // le token retourné par le service d'authentification
		Path:     "/",             // le cookie sera accessible sur l'ensemble du domaine
		HttpOnly: true,            // améliore la sécurité en interdisant l'accès via JavaScript
		// Secure: true,         // à activer en production sous HTTPS
		// Optionnel : définir une date d'expiration si nécessaire
	}
	http.SetCookie(w, cookie)

	// Création de la réponse JSON
	response := &models.LoginResponse{
		Response: &models.ResponseBody{
			Message: "Login successful",
			Status:  http.StatusOK,
		},
		UserId: userId,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
	}
}

func AuthStatusHandler(w http.ResponseWriter, r *http.Request, sessionService *services.SessionService, userService *services.UserService) {
	// Récupérer le cookie "session_token"
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Récupérer la session associée au token
	session, err := sessionService.GetSessionByToken(cookie.Value)
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Vérifier l'expiration de la session si un expires_at est défini
	if session.ExpireAt.Before(time.Now()) {
		http.Error(w, "Session expired", http.StatusUnauthorized)
		return
	}

	// Si la session est valide, renvoyer un JSON indiquant que l'utilisateur est authentifié
	w.Header().Set("Content-Type", "application/json")

	name, _ := userService.GetUserByUUID(session.UserId.String())
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"userId":        session.UserId, // Assurez-vous que le champ s'appelle bien UserID dans votre modèle de Session
		"username":      name,
	})
	if err != nil {
		return
	}
}

// LogoutHandler supprime la session associée au cookie "session_token" et efface le cookie.
func LogoutHandler(w http.ResponseWriter, r *http.Request, sessionService *services.SessionService) {
	// Récupérer le cookie "session_token"
	cookie, err := r.Cookie("session_token")
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, "Missing session token")
		return
	}
	token := cookie.Value

	// Supprimer la session de la base de données via le service
	err = sessionService.DeleteSessionByToken(token)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Failed to delete session")
		return
	}

	// Supprimer le cookie en le remplaçant par un cookie expiré
	expiredCookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, expiredCookie)

	// Rediriger l'utilisateur vers la page principale
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
