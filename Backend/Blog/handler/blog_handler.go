package handler

import (
	"context"

	"github.com/pavlovicisidora/soa-team7/Backend/Blog/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/service"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/pavlovicisidora/soa-team7/Backend/Blog/proto"
)

type BlogHandler struct {
	pb.UnimplementedBlogServiceServer
	blogService service.BlogService
}

func NewBlogHandler(blogService service.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
	}
}

func (h *BlogHandler) CreateBlog(ctx context.Context, req *pb.CreateBlogRequest) (*pb.CreateBlogResponse, error) {
	var images []model.Image
	for _, img := range req.GetImages() {
		images = append(images, model.Image{URL: img.Url})
	}

	createdBlog, err := h.blogService.CreateBlog(
		ctx,
		req.GetTitle(),
		req.GetContent(),
		images,
		req.GetUserId(),
	)
	if err != nil {
		return nil, err
	}

	protoBlog := &pb.Blog{
		Id:        createdBlog.ID.Hex(),
		Title:     createdBlog.Title,
		Content:   createdBlog.Content,
		CreatedAt: timestamppb.New(createdBlog.CreatedAt),
		UserId:    createdBlog.UserID,
	}
	for _, img := range createdBlog.Images {
		protoBlog.Images = append(protoBlog.Images, &pb.Image{Url: img.URL})
	}

	return &pb.CreateBlogResponse{
		Blog: protoBlog,
	}, nil
}

/*
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
}*/
