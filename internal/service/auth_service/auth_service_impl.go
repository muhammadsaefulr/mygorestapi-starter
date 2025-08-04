package service

import (
	"errors"

	"github.com/muhammadsaefulr/NimeStreamAPI/config"

	auth_request_dto "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/auth/request"
	auth_response_dto "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/auth/response"
	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/request"
	userpts_request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_points/request"
	user_model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	system_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/system_service"
	user_point_svc "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_points_service"
	user_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_service"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type authService struct {
	Log          *logrus.Logger
	DB           *gorm.DB
	Validate     *validator.Validate
	UserService  user_service.UserService
	TokenService system_service.TokenService
	UsPointSvc   user_point_svc.UserPointsServiceInterface
}

func NewAuthService(
	db *gorm.DB, validate *validator.Validate, userService user_service.UserService, tokenService system_service.TokenService, user_point_svc user_point_svc.UserPointsServiceInterface,
) AuthService {
	return &authService{
		Log:          utils.Log,
		DB:           db,
		Validate:     validate,
		UserService:  userService,
		UsPointSvc:   user_point_svc,
		TokenService: tokenService,
	}
}

func (s *authService) Register(c *fiber.Ctx, req *auth_request_dto.Register) (*user_model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	user := &request.CreateUser{
		Name:     req.Name,
		Email:    req.Email,
		RoleId:   1,
		Role:     "user",
		Password: req.Password,
	}

	result, err := s.UserService.CreateUser(c, user)

	if err != nil {
		s.Log.Errorf("Failed create user: %+v", err)

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fiber.NewError(fiber.StatusConflict, "Email Or Username already taken")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create user failed")
	}

	s.UsPointSvc.Update(c, &userpts_request.UserPoints{
		UserId:     result.ID.String(),
		TypeUpdate: "add",
		Value:      0,
	})

	return result, nil
}

func (s *authService) Login(c *fiber.Ctx, req *auth_request_dto.Login) (*user_model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	user, err := s.UserService.GetUserByEmail(c, req.Email)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	return user, nil
}
func (s *authService) Logout(c *fiber.Ctx, req *auth_request_dto.Logout) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	token, err := s.TokenService.GetTokenByUserID(c, req.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Token not found")
	}

	err = s.TokenService.DeleteToken(c, config.TokenTypeRefresh, token.UserID.String())

	return err
}

func (s *authService) RefreshAuth(c *fiber.Ctx, req *auth_request_dto.RefreshToken) (*auth_response_dto.Tokens, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	token, err := s.TokenService.GetTokenByUserID(c, req.RefreshToken)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	user, err := s.UserService.GetUserByID(c, token.UserID.String())
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	newTokens, err := s.TokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return newTokens, err
}

func (s *authService) ResetPassword(c *fiber.Ctx, query *auth_request_dto.Token, req *request.UpdatePassOrVerify) error {
	if err := s.Validate.Struct(query); err != nil {
		return err
	}

	userID, err := utils.VerifyToken(query.Token, config.JWTSecret, config.TokenTypeResetPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Token")
	}

	user, err := s.UserService.GetUserByID(c, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Password reset failed")
	}

	if errUpdate := s.UserService.UpdatePassOrVerify(c, req, user.ID.String()); errUpdate != nil {
		return errUpdate
	}

	if errToken := s.TokenService.DeleteToken(c, config.TokenTypeResetPassword, user.ID.String()); errToken != nil {
		return errToken
	}

	return nil
}

func (s *authService) VerifyEmail(c *fiber.Ctx, query *auth_request_dto.Token) error {
	if err := s.Validate.Struct(query); err != nil {
		return err
	}

	userID, err := utils.VerifyToken(query.Token, config.JWTSecret, config.TokenTypeVerifyEmail)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Token")
	}

	user, err := s.UserService.GetUserByID(c, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Verify email failed")
	}

	if errToken := s.TokenService.DeleteToken(c, config.TokenTypeVerifyEmail, user.ID.String()); errToken != nil {
		return errToken
	}

	updateBody := &request.UpdatePassOrVerify{
		VerifiedEmail: true,
	}

	if errUpdate := s.UserService.UpdatePassOrVerify(c, updateBody, user.ID.String()); errUpdate != nil {
		return errUpdate
	}

	return nil
}
