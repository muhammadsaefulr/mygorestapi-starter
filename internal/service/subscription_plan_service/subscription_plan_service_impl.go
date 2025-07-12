package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/subscription_plan/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/subscription_plan"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type SubscriptionPlanService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     repository.SubscriptionPlanRepo
}

func NewSubscriptionPlanService(repo repository.SubscriptionPlanRepo, validate *validator.Validate) SubscriptionPlanServiceInterface {
	return &SubscriptionPlanService{
		Log:      utils.Log,
		Validate: validate,
		Repo:     repo,
	}
}

func (s *SubscriptionPlanService) GetAll(c *fiber.Ctx, params *request.QuerySubscriptionPlan) ([]model.SubscriptionPlan, int64, error) {
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

func (s *SubscriptionPlanService) GetByID(c *fiber.Ctx, id uint) (*model.SubscriptionPlan, error) {
	data, err := s.Repo.GetByID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "SubscriptionPlan not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get subscription_plan by ID failed")
	}
	return data, nil
}

func (s *SubscriptionPlanService) Create(c *fiber.Ctx, req *request.CreateSubscriptionPlan) (*model.SubscriptionPlan, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.CreateSubscriptionPlanToModel(req)
	if err := s.Repo.Create(c.Context(), data); err != nil {
		s.Log.Errorf("Create error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create subscription_plan failed")
	}
	return data, nil
}

func (s *SubscriptionPlanService) Update(c *fiber.Ctx, id uint, req *request.UpdateSubscriptionPlan) (*model.SubscriptionPlan, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.UpdateSubscriptionPlanToModel(req)
	data.ID = id
	if err := s.Repo.Update(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update subscription_plan failed")
	}
	return s.GetByID(c, id)
}

func (s *SubscriptionPlanService) Delete(c *fiber.Ctx, id uint) error {
	if _, err := s.Repo.GetByID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "SubscriptionPlan not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete subscription_plan failed")
	}
	return nil
}
