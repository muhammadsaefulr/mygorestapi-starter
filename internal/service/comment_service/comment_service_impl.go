package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/comment/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/comment/response"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/comment"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
)

type commentService struct {
	commentRepository repository.CommentRepository
}

func NewCommentService(commentRepository repository.CommentRepository) CommentService {
	return &commentService{
		commentRepository: commentRepository,
	}
}

func (c *commentService) CreateComment(ctx *fiber.Ctx, req *request.CreateComment) error {
	authHeader := ctx.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	IdUsr, errVerToken := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if errVerToken != nil {
		return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("Error verifying token: %s", errVerToken.Error()))
	}

	req.UserId = IdUsr

	err := c.commentRepository.CreateComment(ctx.Context(), convert_types.CreateCommentToModel(req))
	if err != nil {
		return err
	}

	return nil
}

func (c *commentService) GetCommentByID(ctx *fiber.Ctx, id uint) (*response.CommentResponse, error) {
	comment, err := c.commentRepository.GetCommentByID(ctx.Context(), id)
	if err != nil {
		return nil, err
	}

	result := &response.CommentResponse{
		ID:      comment.ID,
		UserId:  comment.UserId,
		MovieId: comment.MovieId,
		Content: comment.Content,
	}

	return result, nil
}

func (c *commentService) GetCommentsMovieId(ctx *fiber.Ctx, movieId string) ([]response.CommentResponse, error) {
	comments, err := c.commentRepository.GetCommentByMovieID(ctx.Context(), movieId)
	if err != nil {
		return nil, err
	}

	responses := make([]response.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		responses = append(responses, response.CommentResponse{
			ID:      comment.ID,
			UserId:  comment.UserId,
			MovieId: comment.MovieId,
			Content: comment.Content,
		})
	}

	return responses, nil
}

func (c *commentService) UpdateComment(ctx *fiber.Ctx, req *request.UpdateComment, id string) (*response.CommentResponse, error) {
	commentID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	comment := convert_types.UpdateCommentToModel(req)
	comment.ID = uint(commentID)

	updated, err := c.commentRepository.UpdateComment(ctx.Context(), comment)
	if err != nil {
		return nil, err
	}

	result := &response.CommentResponse{
		ID:      updated.ID,
		UserId:  updated.UserId,
		MovieId: updated.MovieId,
		Content: updated.Content,
	}

	return result, nil
}

func (c *commentService) DeleteComment(ctx *fiber.Ctx, id uint) error {
	return c.commentRepository.DeleteComment(ctx.Context(), id)
}
