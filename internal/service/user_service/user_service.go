package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/request"
	user_model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model/user"
)

type UserService interface {
	CreateUser(c *fiber.Ctx, req *request.CreateUser) (*user_model.User, error)
	CreateGoogleUser(c *fiber.Ctx, req *request.GoogleLogin) (*user_model.User, error)
	GetUserByEmail(c *fiber.Ctx, email string) (*user_model.User, error)
	GetUserByID(c *fiber.Ctx, id string) (*user_model.User, error)
	UpdatePassOrVerify(c *fiber.Ctx, req *request.UpdatePassOrVerify, id string) error
	UpdateUser(c *fiber.Ctx, id string, req *request.UpdateUser) (*user_model.User, error)
	GetAllUser(c *fiber.Ctx, params *request.QueryUser) ([]user_model.User, int64, error)
	DeleteUser(c *fiber.Ctx, id string) error
}
