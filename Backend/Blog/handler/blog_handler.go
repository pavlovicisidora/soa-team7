package handler

import (
	"context"
	"log"

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
	log.Printf("HANDLER: Received CreateBlog request for user ID: %s", req.GetUserId())
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
		log.Printf("ERROR: Failed to create blog: %v", err)
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

	log.Printf("HANDLER: Successfully created blog with ID: %s", createdBlog.ID.Hex())
	return &pb.CreateBlogResponse{
		Blog: protoBlog,
	}, nil
}

func (h *BlogHandler) GetBlog(ctx context.Context, req *pb.GetBlogRequest) (*pb.GetBlogResponse, error) {
	log.Printf("HANDLER: Received GetBlog request for blog ID: %s", req.GetId())
	blog, err := h.blogService.GetBlogByID(ctx, req.GetId())
	if err != nil {
		log.Printf("ERROR: Failed to get blog with ID %s: %v", req.GetId(), err)
		return nil, err
	}

	log.Printf("HANDLER: Successfully fetched blog with ID: %s", req.GetId())
	return &pb.GetBlogResponse{Blog: toProtoBlog(blog)}, nil
}

func (h *BlogHandler) GetAllBlogs(ctx context.Context, req *pb.GetAllBlogsRequest) (*pb.GetAllBlogsResponse, error) {
	log.Printf("HANDLER: Received GetAllBlogs request for %d authors.", len(req.GetFollowedUserIds()))

	followedUserIDs := req.GetFollowedUserIds()

	blogs, err := h.blogService.GetAllBlogs(ctx, followedUserIDs)
	if err != nil {
		log.Printf("ERROR: Failed to get all blogs: %v", err)
		return nil, err
	}

	var protoBlogs []*pb.Blog
	for _, blog := range blogs {
		protoBlogs = append(protoBlogs, toProtoBlog(&blog))
	}

	log.Printf("HANDLER: Successfully fetched %d blogs.", len(protoBlogs))
	return &pb.GetAllBlogsResponse{Blogs: protoBlogs}, nil
}

func (h *BlogHandler) LikeBlog(ctx context.Context, req *pb.LikeBlogRequest) (*pb.LikeBlogResponse, error) {
	log.Printf("HANDLER: Received LikeBlog request for blog ID %s from user ID %s.", req.GetBlogId(), req.GetUserId())
	blog, err := h.blogService.LikeBlog(ctx, req.GetBlogId(), req.GetUserId())
	if err != nil {
		log.Printf("ERROR: Failed to like blog %s: %v", req.GetBlogId(), err)
		return nil, err
	}
	log.Printf("HANDLER: Successfully liked blog %s.", req.GetBlogId())
	return &pb.LikeBlogResponse{Blog: toProtoBlog(blog)}, nil
}

func (h *BlogHandler) UnlikeBlog(ctx context.Context, req *pb.UnlikeBlogRequest) (*pb.UnlikeBlogResponse, error) {
	log.Printf("HANDLER: Received UnlikeBlog request for blog ID %s from user ID %s.", req.GetBlogId(), req.GetUserId())
	blog, err := h.blogService.UnlikeBlog(ctx, req.GetBlogId(), req.GetUserId())
	if err != nil {
		log.Printf("ERROR: Failed to unlike blog %s: %v", req.GetBlogId(), err)
		return nil, err
	}
	log.Printf("HANDLER: Successfully unliked blog %s.", req.GetBlogId())
	return &pb.UnlikeBlogResponse{Blog: toProtoBlog(blog)}, nil
}
