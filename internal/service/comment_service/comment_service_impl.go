package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/comment/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/comment/response"
	responseUser "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
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

	log.Printf("Comment: %+v", req)

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
			Content:   "Menarik...",
			CreatedAt: time.Now().Add(-2 * time.Hour),
			UserDetal: &responseUser.GetUsersResponse{
				ID:              uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				Name:            "Dummy user 1",
				Email:           "dummy1@example.com",
				Role:            "user",
				IsEmailVerified: true,
			},
		},
	}

	responses := make([]response.CommentResponse, 0, len(comments))

	for _, comment := range comments {
		// Mapping untuk replies
		replies := make([]response.CommentResponse, 0, len(comment.Replies))
		for _, reply := range comment.Replies {
			replies = append(replies, response.CommentResponse{
				ID:        reply.ID,
				UserId:    reply.UserId,
				MovieId:   reply.MovieId,
				Content:   reply.Content,
				CreatedAt: reply.CreatedAt,
				Likes:     len(reply.Likes),
				UserDetal: &responseUser.GetUsersResponse{
					ID:              reply.UserId,
					Name:            reply.UserDetail.Name,
					Email:           reply.UserDetail.Email,
					Role:            reply.UserDetail.Role,
					IsEmailVerified: reply.UserDetail.VerifiedEmail,
				},
			})
		}

		// Mapping comment utama
		responses = append(responses, response.CommentResponse{
			ID:           comment.ID,
			UserId:       comment.UserId,
			MovieId:      comment.MovieId,
			Content:      comment.Content,
			CreatedAt:    comment.CreatedAt,
			Likes:        len(comment.Likes),
			CommentReply: replies,
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

func (s *commentService) LikeComment(ctx *fiber.Ctx, commentID uint) error {
	userSession := ctx.Locals("user")
	if userSession == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	user := userSession.(*model.User)

	comment, err := s.commentRepository.GetCommentByID(ctx.Context(), commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Comment not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check comment")
	}

	liked, err := s.commentRepository.HasUserLiked(ctx.Context(), comment.ID, user.ID.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check like status")
	}
	if liked {
		return nil
	}

	newLike := model.CommentLike{
		UserId:    user.ID,
		CommentID: comment.ID,
	}
	if err := s.commentRepository.LikeComment(ctx.Context(), newLike); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to like comment")
	}

	return nil
}

func (c *commentService) DislikeComment(ctx *fiber.Ctx, id uint) error {
	userSession := ctx.Locals("user")
	if userSession == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	user := userSession.(*model.User)

	err := c.commentRepository.DislikeComment(ctx.Context(), id, user.ID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Like not found for this user and comment")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to dislike comment")
	}

	return nil
}
