package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/Backend/APIGateway/middleware"
	comment_proto "github.com/pavlovicisidora/soa-team7/Backend/Blog/proto"
)

type CommentHandler struct {
	client comment_proto.CommentServiceClient
}

func NewCommentHandler(client comment_proto.CommentServiceClient) *CommentHandler {
	return &CommentHandler{client: client}
}

// ServeHTTP registruje rute za comment
func (h *CommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/comments/{blogID}", h.GetCommentsHandler).Methods("GET")
	router.HandleFunc("/comments", h.AddCommentHandler).Methods("POST")
	router.HandleFunc("/comments/{id}", h.UpdateCommentHandler).Methods("PUT")
	router.HandleFunc("/comments/{id}", h.DeleteCommentHandler).Methods("DELETE")

	router.ServeHTTP(w, r)
}

type AddCommentRequest struct {
	BlogID string `json:"blog_id"` // ID bloga na koji se odnosi komentar
	Text   string `json:"text"`
}

// AddCommentHandler kreira novi komentar
func (h *CommentHandler) AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	authorID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || authorID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	var reqBody AddCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reqBody.BlogID == "" {
		http.Error(w, "Blog ID is required", http.StatusBadRequest)
		return
	}

	grpcReq := &comment_proto.AddCommentRequest{
		BlogId:   reqBody.BlogID,
		AuthorId: authorID,
		Text:     reqBody.Text,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.AddComment(ctx, grpcReq)
	if err != nil {
		log.Printf("Failed to add comment via gRPC: %v", err)
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp.GetComment())
}

// GetCommentsHandler vraća sve komentare za dati blog
func (h *CommentHandler) GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["blogID"] // iz path-a

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetComments(ctx, &comment_proto.GetCommentsRequest{BlogId: blogID})
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetComments())
}

// UpdateCommentHandler menja tekst komentara
func (h *CommentHandler) UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]

	var reqBody AddCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.UpdateComment(ctx, &comment_proto.UpdateCommentRequest{
		CommentId: commentID,
		Text:      reqBody.Text,
	})
	if err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetComment())
}

// DeleteCommentHandler briše komentar
func (h *CommentHandler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.DeleteComment(ctx, &comment_proto.DeleteCommentRequest{
		CommentId: commentID,
	})
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
