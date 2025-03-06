package handlers

import (
	"encoding/json"
	"fmt"
	"forum/models"
	"forum/services"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

// CreateCommentHandler gère l'ajout d'un commentaire
func PostCommentHandler(w http.ResponseWriter, r *http.Request, commentService *services.CommentsService, sessionService *services.SessionService) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	// Décodage du JSON
	var payload models.CommentPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, "Unable to decode JSON payload")
		return
	}

	// Vérifier si le cookie session_token existe
	cookie, err := r.Cookie("session_token")
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, "Missing session token")
		return
	}
	token := cookie.Value

	// Récupérer la session
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

	userID := session.UserId

	if payload.Content == "" || payload.PostId == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "Missing required fields: content or post_id")
		return
	}

	// Générer un UUID pour le commentaire
	commentUUID, err := uuid.NewV4()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Failed to generate comment UUID")
		return
	}

	// Insérer le commentaire
	_, err = commentService.InsertComment(userID, payload.PostId, payload.Content, commentUUID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to insert comment: %v", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment added successfully"})
}

// GetCommentsByPostHandler récupère tous les commentaires d'un post
// func GetCommentHandler(w http.ResponseWriter, r *http.Request, commentService *services.CommentsService) {
// 	postID := r.URL.Query().Get("post_id")
// 	if postID == "" {
// 		ErrorHandler(w, r, http.StatusBadRequest, "Missing post_id parameter")
// 		return
// 	}

// 	comments, err := commentsRepositories.GetCommentsByPost(postID)
// 	if err != nil {
// 		ErrorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve comments: %v", err))
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(comments)
// }

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {

}
