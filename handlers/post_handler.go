package handlers

import (
	"fmt"
	"forum/models"
	"forum/services"
	"github.com/gofrs/uuid"
	"net/http"
	"time"
)

func PostsHandler(w http.ResponseWriter, r *http.Request, postService *services.PostsService, sessionService *services.SessionService) {
	// Vérifier que la méthode est POST
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	// Parse du formulaire
	if err := r.ParseForm(); err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, "Unable to parse form")
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

	// Récupérer les valeurs du formulaire pour le post
	title := r.FormValue("title")
	content := r.FormValue("content")
	category := r.FormValue("category")

	// Vérifier que les champs obligatoires sont présents
	if title == "" || content == "" {
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
	_, err = postService.InsertPost(userID, title, content, category, postUUID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to insert post: %v", err))
		return
	}

	// Optionnel : vous pouvez créer un objet post pour un log ou traitement ultérieur
	_ = &models.Post{
		PostId:    postUUID,
		UserId:    userID,
		Title:     title,
		Content:   content,
		Category:  category,
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
