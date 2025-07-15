package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/banner_app/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/banner_app"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type BannerAppService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     repository.BannerAppRepo
}

func NewBannerAppService(repo repository.BannerAppRepo, validate *validator.Validate) BannerAppService {
	return BannerAppService{
		Log:      utils.Log,
		Validate: validate,
		Repo:     repo,
	}
}

func (s *BannerAppService) GetAll(c *fiber.Ctx, params *request.QueryBannerApp) ([]model.BannerApp, int64, error) {
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

func (s *BannerAppService) GetByID(c *fiber.Ctx, id uint) (*model.BannerApp, error) {
	data, err := s.Repo.GetByID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "BannerApp not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get banner_app by ID failed")
	}
	return data, nil
}

func (s *BannerAppService) Create(c *fiber.Ctx, req *request.CreateBannerApp) (*model.BannerApp, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.CreateBannerAppToModel(req)
	if err := s.Repo.Create(c.Context(), data); err != nil {
		s.Log.Errorf("Create error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create banner_app failed")
	}
	return data, nil
}

func (s *BannerAppService) Update(c *fiber.Ctx, id uint, req *request.UpdateBannerApp) (*model.BannerApp, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.UpdateBannerAppToModel(req)
	data.ID = id
	if err := s.Repo.Update(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update banner_app failed")
	}
	return s.GetByID(c, id)
}

func (s *BannerAppService) Delete(c *fiber.Ctx, id uint) error {
	if _, err := s.Repo.GetByID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "BannerApp not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete banner_app failed")
	}
	return nil
}
