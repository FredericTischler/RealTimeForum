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
	fmt.Println("Received request to add comment")

	if r.Method != http.MethodPost {
		fmt.Println("Invalid request method")
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var payload models.CommentPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println("Failed to decode JSON payload:", err)
		ErrorHandler(w, r, http.StatusBadRequest, "Unable to decode JSON payload")
		return
	}

	fmt.Println("Payload decoded:", payload)

	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Missing session token:", err)
		ErrorHandler(w, r, http.StatusUnauthorized, "Missing session token")
		return
	}
	token := cookie.Value

	session, err := sessionService.GetSessionByToken(token)
	if err != nil {
		fmt.Println("Invalid session token:", err)
		ErrorHandler(w, r, http.StatusUnauthorized, "Invalid session token")
		return
	}

	if session.ExpireAt.Before(time.Now()) {
		fmt.Println("Session expired")
		ErrorHandler(w, r, http.StatusUnauthorized, "Session expired")
		return
	}

	userID := session.UserId
	fmt.Println("User ID:", userID)

	if payload.Content == "" || payload.PostId == "" {
		fmt.Println("Missing required fields")
		ErrorHandler(w, r, http.StatusBadRequest, "Missing required fields: content or post_id")
		return
	}

	commentUUID, err := uuid.NewV4()
	if err != nil {
		fmt.Println("Failed to generate UUID:", err)
		ErrorHandler(w, r, http.StatusInternalServerError, "Failed to generate comment UUID")
		return
	}

	fmt.Println("Inserting comment with:", userID, payload.PostId, payload.Content, commentUUID)
	_, err = commentService.InsertComment(userID, payload.PostId, payload.Content, commentUUID)
	if err != nil {
		fmt.Println("Failed to insert comment:", err)
		ErrorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to insert comment: %v", err))
		return
	}

	fmt.Println("Comment inserted successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment added successfully"})
}

func GetCommentHandler(w http.ResponseWriter, r *http.Request, commentService *services.CommentsService) {
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "Missing post_id parameter")
		return
	}

	comments, err := commentService.GetCommentsByPost(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve comments: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {

}
