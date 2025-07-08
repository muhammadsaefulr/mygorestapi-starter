package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/mdl/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	modules "github.com/muhammadsaefulr/NimeStreamAPI/internal/infrastructure/modules/mdl"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type MdlService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
}

func NewMdlService(validate *validator.Validate) MdlServiceInterface {
	return &MdlService{
		Log:      utils.Log,
		Validate: validate,
	}
}

func (s *MdlService) GetAll(c *fiber.Ctx, params *request.QueryMdl) ([]response.MovieDetailOnlyResponse, int64, int64, error) {
	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, 0, fiber.NewError(fiber.StatusBadRequest, "Invalid request param")
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	chromeCtx, cancel := modules.NewChromeContext()
	defer cancel()

	result, total, totalPages, err := modules.FetchMDLMedia(chromeCtx, params.Category, params.Search, params.Page, params.Limit)

	if err != nil {
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, "Internal server error: "+err.Error())
	}

	return result, int64(total), int64(totalPages), nil
}
