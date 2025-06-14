package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/watchlist"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type newWatchlistService struct {
	Log        *logrus.Logger
	Validate   *validator.Validate
	Repository repository.WatchlistRepo
}

func NewWatchlistService(repo repository.WatchlistRepo, validate *validator.Validate) WatchlistService {
	return &newWatchlistService{
		Log:        utils.Log,
		Validate:   validate,
		Repository: repo,
	}
}

func (s *newWatchlistService) GetAllWatchlist(c *fiber.Ctx, params *request.QueryWatchlist) ([]model.Watchlist, int64, error) {
	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	return s.Repository.GetAllWatchlist(c.Context(), params)
}

func (s *newWatchlistService) GetWatchlistByID(c *fiber.Ctx, id uint) (*model.Watchlist, error) {
	data, err := s.Repository.GetWatchlistByID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "Watchlist not found")
	}
	if err != nil {
		s.Log.Errorf("GetWatchlistByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve watchlist")
	}
	return data, nil
}

func (s *newWatchlistService) CreateWatchlist(c *fiber.Ctx, req *request.CreateWatchlist) (*model.Watchlist, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	data := convert_types.CreateWatchlistToModel(req)

	if err := s.Repository.CreateWatchlist(c.Context(), data); err != nil {
		s.Log.Errorf("CreateWatchlist error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create watchlist")
	}
	return data, nil
}

func (s *newWatchlistService) UpdateWatchlist(c *fiber.Ctx, id uint, req *request.UpdateWatchlist) (*model.Watchlist, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	data := convert_types.UpdateWatchlistToModel(req)
	data.ID = id

	if err := s.Repository.UpdateWatchlist(c.Context(), data); err != nil {
		s.Log.Errorf("UpdateWatchlist error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to update watchlist")
	}
	return s.GetWatchlistByID(c, id)
}

func (s *newWatchlistService) DeleteWatchlist(c *fiber.Ctx, id uint) error {
	if _, err := s.Repository.GetWatchlistByID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "Watchlist not found")
	}

	if err := s.Repository.DeleteWatchlist(c.Context(), id); err != nil {
		s.Log.Errorf("DeleteWatchlist error: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete watchlist")
	}
	return nil
}
