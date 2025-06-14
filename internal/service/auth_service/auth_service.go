package service

import (
	"github.com/gofiber/fiber/v2"
	auth_request_dto "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/auth/request"
	auth_response_dto "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/auth/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type AuthService interface {
	Register(c *fiber.Ctx, req *auth_request_dto.Register) (*model.User, error)
	Login(c *fiber.Ctx, req *auth_request_dto.Login) (*model.User, error)
	Logout(c *fiber.Ctx, req *auth_request_dto.Logout) error
	RefreshAuth(c *fiber.Ctx, req *auth_request_dto.RefreshToken) (*auth_response_dto.Tokens, error)
	ResetPassword(c *fiber.Ctx, query *auth_request_dto.Token, req *request.UpdatePassOrVerify) error
	VerifyEmail(c *fiber.Ctx, query *auth_request_dto.Token) error
}
