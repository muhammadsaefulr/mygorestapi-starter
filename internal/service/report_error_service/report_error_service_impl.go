package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/report_error/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/report_error"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type ReportErrorService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     repository.ReportErrorRepo
}

func NewReportErrorService(repo repository.ReportErrorRepo, validate *validator.Validate) ReportErrorServiceInterface {
	return &ReportErrorService{
		Log:      utils.Log,
		Validate: validate,
		Repo:     repo,
	}
}

func (s *ReportErrorService) GetAll(c *fiber.Ctx, params *request.QueryReportError) ([]model.ReportError, int64, error) {
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

func (s *ReportErrorService) GetByID(c *fiber.Ctx, id string) (*model.ReportError, error) {
	data, err := s.Repo.GetByID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "Error Report not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get Error Report by ID failed")
	}
	return data, nil
}

func (s *ReportErrorService) Create(c *fiber.Ctx, req *request.CreateReportError) (*model.ReportError, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	user := c.Locals("user").(*model.User)
	req.ReportedBy = user.ID.String()

	data := convert_types.CreateReportErrorToModel(req)
	if err := s.Repo.Create(c.Context(), data); err != nil {
		s.Log.Errorf("Create error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create Error Report failed")
	}
	return data, nil
}

func (s *ReportErrorService) Update(c *fiber.Ctx, id string, req *request.UpdateReportError) (*model.ReportError, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	user := c.Locals("user").(*model.User)
	req.HandledBy = user.ID.String()

	data := convert_types.UpdateReportErrorToModel(req)

	convUUID, errParUUID := uuid.Parse(id)
	if errParUUID != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a valid User UUID")
	}

	data.ReportId = convUUID
	if err := s.Repo.Update(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update Error Report failed")
	}
	return s.GetByID(c, id)
}

func (s *ReportErrorService) Delete(c *fiber.Ctx, id string) error {
	if _, err := s.Repo.GetByID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "ReportError not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete Error Report failed")
	}
	return nil
}
