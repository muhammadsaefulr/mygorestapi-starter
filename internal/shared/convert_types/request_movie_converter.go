package convert_types

import (
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_movie/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_movie/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateRequestMovieToModel(req *request.CreateRequestMovie) *model.RequestMovie {
	return &model.RequestMovie{
		UserIdRequest: uuid.MustParse(req.UserIdRequest),
		Title:         req.Title,
		TypeMovie:     req.TypeMovie,
		Genre:         req.Genre,
		Description:   req.Description,
		StatusMovie:   req.StatusMovie,
		StatusRequest: req.StatusRequest,
	}
}

func UpdateRequestMovieToModel(req *request.UpdateRequestMovie) *model.RequestMovie {
	return &model.RequestMovie{
		ID:            req.ID,
		Title:         req.Title,
		TypeMovie:     req.TypeMovie,
		Genre:         req.Genre,
		Description:   req.Description,
		StatusMovie:   req.StatusMovie,
		StatusRequest: req.StatusRequest,
		UpdatedAt:     req.UpdatedAt,
	}
}

func ModelRequestMovieToResponse(req *model.RequestMovie) *response.RequestMovieResponse {
	return &response.RequestMovieResponse{
		ID:            req.ID,
		UserIdRequest: req.UserIdRequest,
		Title:         req.Title,
		TypeMovie:     req.TypeMovie,
		Genre:         req.Genre,
		Description:   req.Description,
		StatusMovie:   req.StatusMovie,
		StatusRequest: req.StatusRequest,
		RequestedBy: response.UserResponse{
			ID:   req.UserIdRequest,
			Name: req.RequestedBy.Name,
			Role: req.RequestedBy.Role,
		},
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}
}
