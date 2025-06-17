package convert_types

import (
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/comment/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/comment/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateCommentToModel(req *request.CreateComment) *model.Comment {
	return &model.Comment{
		UserId:  uuid.MustParse(req.UserId),
		MovieId: req.MovieId,
		Content: req.Content,
	}
}

func UpdateCommentToModel(req *request.UpdateComment) *model.Comment {
	return &model.Comment{
		ID:      req.ID,
		UserId:  uuid.MustParse(req.UserId),
		MovieId: req.MovieId,
		Content: req.Content,
	}
}

func CommentResponseToModel(res *response.CommentResponse) *model.Comment {
	return &model.Comment{
		ID:      res.ID,
		UserId:  res.UserId,
		MovieId: res.MovieId,
		Content: res.Content,
	}
}
