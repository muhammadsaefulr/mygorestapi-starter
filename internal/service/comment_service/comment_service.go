package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/comment/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/comment/response"
)

type CommentService interface {
	CreateComment(ctx *fiber.Ctx, req *request.CreateComment) error
	GetCommentByID(ctx *fiber.Ctx, id uint) (*response.CommentResponse, error)
	GetCommentsMovieId(ctx *fiber.Ctx, movieId string) ([]response.CommentResponse, error)
	UpdateComment(ctx *fiber.Ctx, req *request.UpdateComment, id string) (*response.CommentResponse, error)
	DeleteComment(ctx *fiber.Ctx, id uint) error
}
