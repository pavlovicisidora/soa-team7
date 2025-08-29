package handler

import (
	"context"

	"github.com/pavlovicisidora/soa-team7/Backend/Blog/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/service"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/pavlovicisidora/soa-team7/Backend/Blog/proto"
)

func toProtoBlog(blog *model.Blog) *pb.Blog {
	return &pb.Blog{
		Id:        blog.ID.Hex(),
		Title:     blog.Title,
		Content:   blog.Content,
		CreatedAt: timestamppb.New(blog.CreatedAt),
		UserId:    blog.UserID,
		Images:    toProtoImages(blog.Images),
		LikeCount: int32(len(blog.LikedBy)),
	}
}

func toProtoImages(images []model.Image) []*pb.Image {
	var protoImages []*pb.Image
	for _, img := range images {
		protoImages = append(protoImages, &pb.Image{Url: img.URL})
	}
	return protoImages
}

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

func (h *BlogHandler) GetBlog(ctx context.Context, req *pb.GetBlogRequest) (*pb.GetBlogResponse, error) {
	blog, err := h.blogService.GetBlogByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.GetBlogResponse{Blog: toProtoBlog(blog)}, nil
}

func (h *BlogHandler) GetAllBlogs(ctx context.Context, req *pb.GetAllBlogsRequest) (*pb.GetAllBlogsResponse, error) {
	blogs, err := h.blogService.GetAllBlogs(ctx)
	if err != nil {
		return nil, err
	}
	var protoBlogs []*pb.Blog
	for _, blog := range blogs {
		protoBlogs = append(protoBlogs, toProtoBlog(&blog))
	}
	return &pb.GetAllBlogsResponse{Blogs: protoBlogs}, nil
}

func (h *BlogHandler) LikeBlog(ctx context.Context, req *pb.LikeBlogRequest) (*pb.LikeBlogResponse, error) {
	blog, err := h.blogService.LikeBlog(ctx, req.GetBlogId(), req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &pb.LikeBlogResponse{Blog: toProtoBlog(blog)}, nil
}

func (h *BlogHandler) UnlikeBlog(ctx context.Context, req *pb.UnlikeBlogRequest) (*pb.UnlikeBlogResponse, error) {
	blog, err := h.blogService.UnlikeBlog(ctx, req.GetBlogId(), req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &pb.UnlikeBlogResponse{Blog: toProtoBlog(blog)}, nil
}
