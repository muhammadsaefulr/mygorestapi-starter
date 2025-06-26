package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/comment/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/comment/response"
	responseUser "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/response"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/comment"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"gorm.io/gorm"
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

	if err := c.commentRepository.CreateComment(ctx.Context(), convert_types.CreateCommentToModel(req)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create comment")
	}

	return nil
}

func (c *commentService) GetCommentByID(ctx *fiber.Ctx, id uint) (*response.CommentResponse, error) {
	comment, err := c.commentRepository.GetCommentByID(ctx.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "Comment not found")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get comment")
	}

	return &response.CommentResponse{
		ID:      comment.ID,
		UserId:  comment.UserId,
		MovieId: comment.MovieId,
		Content: comment.Content,
		UserDetal: &responseUser.GetUsersResponse{
			ID:              comment.UserId,
			Name:            comment.UserDetail.Name,
			Email:           comment.UserDetail.Email,
			Role:            comment.UserDetail.Role,
			IsEmailVerified: comment.UserDetail.VerifiedEmail,
		},
	}, nil
}

func (c *commentService) GetCommentsMovieId(ctx *fiber.Ctx, movieId string) ([]response.CommentResponse, error) {
	comments, err := c.commentRepository.GetCommentByMovieID(ctx.Context(), movieId)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get comments for movie")
	}

	dummyComments := []response.CommentResponse{
		{
			ID:        1,
			UserId:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			MovieId:   movieId,
			Content:   "Halo Aku Dummy !",
			CreatedAt: time.Now().Add(-2 * time.Hour),
			UserDetal: &responseUser.GetUsersResponse{
				ID:              uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				Name:            "Dummy user 1",
				Email:           "dummy1@example.com",
				Role:            "user",
				IsEmailVerified: true,
			},
		},
		{
			ID:        2,
			UserId:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			MovieId:   movieId,
			Content:   "Hidup Jokowwieee",
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UserDetal: &responseUser.GetUsersResponse{
				ID:              uuid.MustParse("22222222-2222-2222-2222-222222222222"),
				Name:            "Dummy user 2",
				Email:           "dummy2@example.com",
				Role:            "user",
				IsEmailVerified: false,
			},
		},
		{
			ID:        3,
			UserId:    uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			MovieId:   movieId,
			Content:   "Gemoy Joget",
			CreatedAt: time.Now(),
			UserDetal: &responseUser.GetUsersResponse{
				ID:              uuid.MustParse("33333333-3333-3333-3333-333333333333"),
				Name:            "Dummy user 3",
				Email:           "dummy3@example.com",
				Role:            "user",
				IsEmailVerified: true,
			},
		},
	}

	responses := make([]response.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		responses = append(responses, response.CommentResponse{
			ID:      comment.ID,
			UserId:  comment.UserId,
			MovieId: comment.MovieId,
			Content: comment.Content,
			UserDetal: &responseUser.GetUsersResponse{
				ID:              comment.UserId,
				Name:            comment.UserDetail.Name,
				Email:           comment.UserDetail.Email,
				Role:            comment.UserDetail.Role,
				IsEmailVerified: comment.UserDetail.VerifiedEmail,
			},
		})
	}

	responses = append(responses, dummyComments...)

	return responses, nil
}

func (c *commentService) UpdateComment(ctx *fiber.Ctx, req *request.UpdateComment, id string) (*response.CommentResponse, error) {
	commentID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid comment ID Because : %s", err.Error()))
	}

	commentsDetails, errGetComm := c.GetCommentByID(ctx, uint(commentID))
	if errGetComm != nil {
		if errors.Is(errGetComm, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Comment not found")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get comment")
	}

	commentsDetails.Content = req.Content

	comment := convert_types.CommentResponseToModel(commentsDetails)
	comment.ID = uint(commentID)

	updated, err := c.commentRepository.UpdateComment(ctx.Context(), comment)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to update comment")
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
