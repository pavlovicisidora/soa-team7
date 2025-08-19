package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/service"
)

type BlogHandler struct {
	blogService service.BlogService
}

func NewBlogHandler(blogService service.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
	}
}

func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody struct {
		Title   string        `json:"title"`
		Content string        `json:"content"`
		UserID  string        `json:"user_id"`
		Images  []model.Image `json:"images"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	blog, err := h.blogService.CreateBlog(
		r.Context(),
		requestBody.Title,
		requestBody.Content,
		requestBody.Images,
		requestBody.UserID,
	)
	if err != nil {
		http.Error(w, "Failed to create blog post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	blogs, err := h.blogService.GetAllBlogs(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve blog posts", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) GetBlogByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	blog, err := h.blogService.GetBlogByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Blog post not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(blog)
}
