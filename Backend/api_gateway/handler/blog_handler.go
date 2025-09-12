package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/Backend/APIGateway/middleware"
	blog_proto "github.com/pavlovicisidora/soa-team7/Backend/Blog/proto"
	follower_proto "github.com/pavlovicisidora/soa-team7/Backend/Follower/proto"
)

func (h *BlogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/blogs", h.CreateBlogHandler).Methods("POST")
	router.HandleFunc("/blogs", h.GetAllBlogsHandler).Methods("GET")
	router.HandleFunc("/blogs/{id}", h.GetBlogByIDHandler).Methods("GET")
	router.HandleFunc("/blogs/{id}/like", h.LikeBlogHandler).Methods("POST")
	router.HandleFunc("/blogs/{id}/like", h.UnlikeBlogHandler).Methods("DELETE")

	router.ServeHTTP(w, r)
}

type CreateBlogRequest struct {
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Images  []Image `json:"images"`
}
type Image struct {
	URL string `json:"url"`
}

type BlogHandler struct {
	client         blog_proto.BlogServiceClient
	followerClient follower_proto.FollowerServiceClient
}

func NewBlogHandler(client blog_proto.BlogServiceClient, followerClient follower_proto.FollowerServiceClient) *BlogHandler {
	return &BlogHandler{
		client:         client,
		followerClient: followerClient,
	}
}

func (h *BlogHandler) CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var reqBody CreateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var protoImages []*blog_proto.Image
	for _, img := range reqBody.Images {
		protoImages = append(protoImages, &blog_proto.Image{Url: img.URL})
	}

	grpcRequest := &blog_proto.CreateBlogRequest{
		Title:   reqBody.Title,
		Content: reqBody.Content,
		UserId:  userID,
		Images:  protoImages,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.CreateBlog(ctx, grpcRequest)
	if err != nil {
		log.Printf("Failed to create blog via gRPC: %v", err)
		http.Error(w, "Failed to create blog", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp.GetBlog()); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *BlogHandler) GetAllBlogsHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)

	defer cancel()

	followingResp, err := h.followerClient.GetFollowing(ctx, &follower_proto.GetFollowingRequest{UserId: userID})
	if err != nil {
		log.Printf("Failed to get following list from Follower service: %v", err)
		http.Error(w, "Failed to get user's following list", http.StatusInternalServerError)
		return
	}

	var followingUserIDs []string
	for _, user := range followingResp.Users {
		followingUserIDs = append(followingUserIDs, user.UserId)
	}

	followingUserIDs = append(followingUserIDs, userID)

	grpcRequest := &blog_proto.GetAllBlogsRequest{
		FollowedUserIds: followingUserIDs,
	}

	resp, err := h.client.GetAllBlogs(ctx, grpcRequest)
	if err != nil {
		log.Printf("Failed to get blogs from Blog service: %v", err)
		http.Error(w, "Failed to get blogs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetBlogs())
}

func (h *BlogHandler) GetBlogByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetBlog(ctx, &blog_proto.GetBlogRequest{Id: id})
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetBlog())
}

func (h *BlogHandler) LikeBlogHandler(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]
	userID := r.Context().Value(middleware.UserKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.LikeBlog(ctx, &blog_proto.LikeBlogRequest{BlogId: blogID, UserId: userID})
	if err != nil {
		http.Error(w, "Failed to like blog", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetBlog())
}

func (h *BlogHandler) UnlikeBlogHandler(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]
	userID := r.Context().Value(middleware.UserKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.UnlikeBlog(ctx, &blog_proto.UnlikeBlogRequest{BlogId: blogID, UserId: userID})
	if err != nil {
		http.Error(w, "Failed to unlike blog", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.GetBlog())
}
