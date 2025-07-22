package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	subs_plan_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/subscription_plan_service"
	badge_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_badge_service"
	user_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_service"

	request_bdge "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_badge/request"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_subscription/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/user_subscription"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type UserSubscriptionService struct {
	Log         *logrus.Logger
	Validate    *validator.Validate
	Repo        repository.UserSubscriptionRepo
	UserSvc     user_service.UserService
	SubsPlanSvc subs_plan_service.SubscriptionPlanServiceInterface
	BadgeSvc    badge_service.UserBadgeServiceInterface
}

func NewUserSubscriptionService(repo repository.UserSubscriptionRepo, validate *validator.Validate, userSvc user_service.UserService, subsPlanSvc subs_plan_service.SubscriptionPlanServiceInterface, badgeSvc badge_service.UserBadgeServiceInterface) UserSubscriptionServiceInterface {
	return &UserSubscriptionService{
		Log:         utils.Log,
		Validate:    validate,
		Repo:        repo,
		UserSvc:     userSvc,
		SubsPlanSvc: subsPlanSvc,
		BadgeSvc:    badgeSvc,
	}
}

func (s *UserSubscriptionService) GetAll(c *fiber.Ctx, params *request.QueryUserSubscription) ([]model.UserSubscription, int64, error) {
	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	return s.Repo.GetAll(c.Context(), params)
}

func (s *UserSubscriptionService) GetByUserID(c *fiber.Ctx, id string) (*model.UserSubscription, error) {
	data, err := s.Repo.GetByUserID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "User Subscription not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get User Subscription by ID failed")
	}
	return data, nil
}

func (s *UserSubscriptionService) Create(c *fiber.Ctx, req *request.CreateUserSubscription) (*model.UserSubscription, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	_, err := s.SubsPlanSvc.GetByID(c, req.SubscriptionPlanId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Subscription Plan tidak ditemukan atau tidak valid")
		}
		return nil, err
	}

	_, err = s.UserSvc.GetUserByID(c, req.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusBadRequest, "User tidak ditemukan atau tidak valid")
		}
		return nil, err
	}

	sessionNow := c.Locals("user").(*model.User)
	data := convert_types.CreateUserSubscriptionToModel(req)
	updatedBy := uuid.MustParse(sessionNow.ID.String())
	data.UpdatedBy = updatedBy

	badgeSubmt := &request_bdge.CreateUserBadgeInfo{
		UserID:    data.UserID.String(),
		BadgeID:   1,
		Note:      "User Subscription",
		HandledBy: sessionNow.ID.String(),
	}

	if err := s.BadgeSvc.CreateUserBadgeInfo(c, badgeSubmt); err != nil {
		// s.Log.Errorf("Create User Badge Info error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create User Badge Info failed")
	}

	if err := s.Repo.Create(c.Context(), data); err != nil {
		// s.Log.Errorf("Create error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create User Subscription failed")
	}
	return data, nil
}

func (s *UserSubscriptionService) UpdateByUserId(c *fiber.Ctx, id string, req *request.UpdateUserSubscription) (*model.UserSubscription, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	data := convert_types.UpdateUserSubscriptionToModel(req)
	data.UserID = uuid.MustParse(id)

	_, err := s.SubsPlanSvc.GetByID(c, req.SubscriptionPlanId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Subscription Plan tidak ditemukan atau tidak valid")
		}
		return nil, err
	}

	_, err = s.UserSvc.GetUserByID(c, req.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusBadRequest, "User tidak ditemukan atau tidak valid")
		}
		return nil, err
	}

	updatedBy := uuid.MustParse(c.Locals("user").(string))
	data.UpdatedBy = updatedBy

	if err := s.Repo.UpdateByUserId(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update User Subscription failed")
	}
	return s.GetByUserID(c, id)
}

func (s *UserSubscriptionService) DeleteByUserId(c *fiber.Ctx, id string) error {
	if _, err := s.Repo.GetByUserID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "UserSubscription not found")
	}
	if err := s.Repo.DeleteByUserId(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete User Subscription failed")
	}
	return nil
}
