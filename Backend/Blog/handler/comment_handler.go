package handler

import (
	"context"

	"github.com/pavlovicisidora/soa-team7/Backend/Blog/model"
	"github.com/pavlovicisidora/soa-team7/Backend/Blog/service"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/pavlovicisidora/soa-team7/Backend/Blog/proto"
)

func toProtoComment(comment *model.Comment) *pb.Comment {
	var createdAtProto, updatedAtProto *timestamppb.Timestamp

	if !comment.CreatedAt.IsZero() {
		createdAtProto = timestamppb.New(comment.CreatedAt)
	}

	if comment.UpdatedAt != nil {
		updatedAtProto = timestamppb.New(*comment.UpdatedAt)
	}

	return &pb.Comment{
		Id:             comment.ID.Hex(),
		BlogId:         comment.BlogID.Hex(),
		AuthorId:       comment.AuthorID,
		Text:           comment.Text,
		CreatedAt:      createdAtProto,
		UpdatedAt:      updatedAtProto,
	}
}

type CommentHandler struct {
	pb.UnimplementedCommentServiceServer
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) AddComment(ctx context.Context, req *pb.AddCommentRequest) (*pb.AddCommentResponse, error) {
	createdComment, err := h.commentService.AddComment(
		ctx,
		req.GetBlogId(),
		req.GetAuthorId(),
		req.GetText(),
	)
	if err != nil {
		return nil, err
	}

	return &pb.AddCommentResponse{
		Comment: toProtoComment(createdComment),
	}, nil
}

func (h *CommentHandler) GetComments(ctx context.Context, req *pb.GetCommentsRequest) (*pb.GetCommentsResponse, error) {
	comments, err := h.commentService.GetComments(ctx, req.GetBlogId())
	if err != nil {
		return nil, err
	}

	var protoComments []*pb.Comment
	for _, c := range comments {
		// ovde je c već model.Comment, zato ide &c
		protoComments = append(protoComments, toProtoComment(&c))
	}

	return &pb.GetCommentsResponse{
		Comments: protoComments,
	}, nil
}

func (h *CommentHandler) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	updatedComment, err := h.commentService.UpdateComment(
		ctx,
		req.GetCommentId(), // koristi samo commentID
		req.GetText(),
	)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateCommentResponse{
		Comment: toProtoComment(updatedComment),
	}, nil
}

func (h *CommentHandler) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	err := h.commentService.DeleteComment(
		ctx,
		req.GetCommentId(), // koristi samo commentID
	)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteCommentResponse{Success: true}, nil
}
