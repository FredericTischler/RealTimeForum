package handlers

import (
	"encoding/json"
	"fmt"
	"forum/models"
	"forum/services"
	"github.com/gofrs/uuid"
	"net/http"
	"time"
)

func PostsHandler(w http.ResponseWriter, r *http.Request, postService *services.PostsService, sessionService *services.SessionService) {
	fmt.Println("post handler")
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	// Décodage du JSON depuis r.Body
	var payload models.PostPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, "Unable to decode JSON payload")
		return
	}

	// Récupérer le cookie "session_token"
	cookie, err := r.Cookie("session_token")
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, "Missing session token")
		return
	}
	token := cookie.Value

	// Récupérer la session associée à ce token
	session, err := sessionService.GetSessionByToken(token)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, "Invalid session token")
		return
	}

	// Vérifier si la session est expirée
	if session.ExpireAt.Before(time.Now()) {
		ErrorHandler(w, r, http.StatusUnauthorized, "Session expired")
		return
	}

	// L'UUID de l'utilisateur est récupéré depuis la session
	userID := session.UserId

	// Vérifier que les champs obligatoires sont présents
	if payload.Title == "" || payload.Content == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "Missing required fields: title and content")
		return
	}

	// Générer un nouvel UUID pour le post
	postUUID, err := uuid.NewV4()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Failed to generate post UUID")
		return
	}

	// Insertion du post via le service
	_, err = postService.InsertPost(userID, payload.Title, payload.Content, payload.Category, postUUID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to insert post: %v", err))
		return
	}

	// Optionnel : création de l'objet post pour journalisation
	_ = &models.Post{
		PostId:    postUUID,
		UserId:    userID,
		Title:     payload.Title,
		Content:   payload.Content,
		Category:  payload.Category,
		CreatedAt: time.Now(),
	}

	// Redirection vers la page principale après la création du post
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func GetPostsHandler(w http.ResponseWriter, r *http.Request, postService *services.PostsService) {
	// Implement the logic to handle GET requests for posts
}

func GetPostsByIdHandler(w http.ResponseWriter, r *http.Request) {

}
