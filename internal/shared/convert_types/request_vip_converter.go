package convert_types

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_vip/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_vip/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateRequestVipToModel(req *request.CreateRequestVip) *model.RequestVip {
	return &model.RequestVip{
		UserID:        uuid.MustParse(req.UserId),
		PaymentMethod: req.PaymentMethod,
		AtasNamaTf:    req.Name,
		Email:         req.Email,
		BuktiTf:       req.BuktiTfStr,
		StatusAcc:     req.StatusAcc,
	}
}

func UpdateRequestVipToModel(req *request.UpdateRequestVip) *model.RequestVip {
	return &model.RequestVip{
		ID:         req.ID,
		AtasNamaTf: req.Name,
		Email:      req.Email,
		StatusAcc:  req.StatusAcc,
		UpdatedBy:  uuid.MustParse(req.UpdatedBy),
	}
}

func VipModelToResponse(data []model.RequestVip) []response.RequestVipResponse {
	results := make([]response.RequestVipResponse, 0, len(data))
	for _, v := range data {
		results = append(results, response.RequestVipResponse{
			ID:            v.ID,
			UserId:        v.UserID.String(),
			Name:          v.AtasNamaTf,
			Email:         v.Email,
			BuktiTf:       fmt.Sprintf("%s/minio/bukti-tf/%s", config.AppUrl, v.BuktiTf),
			StatusAcc:     v.StatusAcc,
			UpdatedBy:     v.UpdatedBy.String(),
			CreatedAt:     v.CreatedAt,
			UpdatedAt:     v.UpdatedAt,
			PaymentMethod: v.PaymentMethod,
		})
	}
	return results
}
