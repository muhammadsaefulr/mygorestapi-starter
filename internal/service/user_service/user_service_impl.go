package service

import (
	"errors"

	"firebase.google.com/go/v4/auth"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/request"
	user_model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/user"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userService struct {
	Log          *logrus.Logger
	Validate     *validator.Validate
	UserRepo     repository.UserRepo
	FirebaseAuth *auth.Client
}

func NewUserService(userRepo repository.UserRepo, validate *validator.Validate, FirebaseAuth *auth.Client) UserService {
	return &userService{
		Log:          utils.Log,
		Validate:     validate,
		UserRepo:     userRepo,
		FirebaseAuth: FirebaseAuth,
	}
}

func (s *userService) GetAllUser(c *fiber.Ctx, params *request.QueryUser) ([]user_model.User, int64, error) {
	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	// Default fallback kalau user tidak isi
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	users, total, err := s.UserRepo.GetAllUser(c.Context(), params)
	if err != nil {
		s.Log.Errorf("Failed to get users: %+v", err)
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, "Get users failed")
	}

	return users, total, nil
}

func (s *userService) CreateUser(c *fiber.Ctx, req *request.CreateUser) (*user_model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	if err, _ := s.GetUserByEmail(c, req.Email); err != nil {
		s.Log.Errorf("User Email already exists: %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "Email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.Log.Errorf("Hash password failed: %+v", err)
		return nil, err
	}

	req.Password = hashedPassword

	user := convert_types.CreateUserToUserModel(req)

	if err := s.UserRepo.CreateUser(c.Context(), user); err != nil {
		s.Log.Errorf("CreateUser failed: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create user failed")
	}

	return user, nil
}

func (s *userService) LoginWithFirebaseToken(c *fiber.Ctx, idToken string) (*user_model.User, error) {
	token, err := s.FirebaseAuth.VerifyIDToken(c.Context(), idToken)
	if err != nil {
		s.Log.Errorf("Failed to verify Firebase ID token: %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
	}

	email, _ := token.Claims["email"].(string)
	name, _ := token.Claims["name"].(string)
	firebaseUID := token.UID

	if email == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Email not found in Firebase token")
	}

	user, err := s.UserRepo.GetUserByEmail(c.Context(), email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newUser := &user_model.User{
			Name:        name,
			FirebaseUID: &firebaseUID,
			Email:       email,
		}

		if err := s.UserRepo.CreateUser(c.Context(), newUser); err != nil {
			s.Log.Errorf("Failed to create user from Firebase token: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
		}

		return newUser, nil
	}

	if err != nil {
		s.Log.Errorf("Error fetching user by email: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	return user, nil
}

func (s *userService) GoogleAuthHandler(c *fiber.Ctx, req *request.GoogleLogin) (*user_model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	userFromDB, err := s.UserRepo.GetUserByEmail(c.Context(), req.Email)
	if err == nil && userFromDB != nil {
		if userFromDB.Password != "" {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Silakan login menggunakan email dan password")
		}

		return userFromDB, nil
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Errorf("Unexpected error checking user email: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	createUser := &user_model.User{
		Name:          req.Name,
		Email:         req.Email,
		VerifiedEmail: false,
	}

	if createErr := s.UserRepo.CreateUser(c.Context(), createUser); createErr != nil {
		s.Log.Errorf("Failed to create user: %+v", createErr)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create user failed")
	}

	return createUser, nil
}

func (s *userService) UpdatePassOrVerify(c *fiber.Ctx, req *request.UpdatePassOrVerify, id string) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	if req.Password == "" && !req.VerifiedEmail {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Request")
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}
		req.Password = hashedPassword
	}

	updatesBody := convert_types.UpdatePassOrVerifyToUserModel(req)

	err := s.UserRepo.UpdateUser(c.Context(), updatesBody)

	if err != nil {
		s.Log.Errorf("Failed to update user password or verifiedEmail: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Update user password or verifiedEmail failed")
	}

	return err
}

func (s *userService) GetUserByEmail(c *fiber.Ctx, email string) (*user_model.User, error) {
	user, err := s.UserRepo.GetUserByEmail(c.Context(), email)

	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if err != nil {
		s.Log.Errorf("GetUserByEmail failed: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get user by email failed")
	}
	return user, nil
}

func (s *userService) GetUserSession(c *fiber.Ctx) (*user_model.User, error) {
	results, err := s.GetUserByID(c, c.Locals("user").(*user_model.User).ID.String())

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get user session failed")
	}

	if results.UserBadge == nil || len(*results.UserBadge) == 0 {
		results.UserBadge = &[]user_model.UserBadgeInfo{
			{
				BadgeID: 0,
			},
		}
	}

	if results.UserPoint.UserID == uuid.Nil {
		results.UserPoint = &user_model.UserPoints{
			ID:     uuid.New(),
			UserID: results.ID,
			Value:  0,
		}
	}
	return results, nil
}

func (s *userService) GetUserByID(c *fiber.Ctx, id string) (*user_model.User, error) {
	user, err := s.UserRepo.GetUserByID(c.Context(), id)

	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if err != nil {
		s.Log.Errorf("GetUserByID failed: %+v", err)
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateUser(c *fiber.Ctx, id string, req *request.UpdateUser) (*user_model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	userLocal := c.Locals("user").(*user_model.User)

	existedMail, errMail := s.GetUserByEmail(c, req.Email)
	if errMail != nil && errMail != gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error checking email")
	}

	if userLocal.ID != existedMail.ID || existedMail.Role != "admin" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "You Cannot Edit This Data !")
	}

	if existedMail != nil && existedMail.ID.String() != id {
		return nil, fiber.NewError(fiber.StatusConflict, "Email already exists")
	}

	parsedID, erruuid := uuid.Parse(id)
	if erruuid != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid UUID")
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		req.Password = hashedPassword
	}

	user := convert_types.UpdateUserToUserModel(req)

	user.ID = parsedID

	err := s.UserRepo.UpdateUser(c.Context(), user)

	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update user failed")
	}

	usr, err := s.GetUserByID(c, id)

	if err != nil {
		return nil, err
	}

	return usr, err
}

func (s *userService) DeleteUser(c *fiber.Ctx, id string) error {
	_, errFind := s.UserRepo.GetUserByID(c.Context(), id)

	if errFind == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	err := s.UserRepo.DeleteUser(c.Context(), id)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete user failed")
	}

	return nil
}
