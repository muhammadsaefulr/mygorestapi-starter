package convert_types

import (
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_vip/request"
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
		BuktiTf:    req.BuktiTf,
		StatusAcc:  req.StatusAcc,
		UpdatedBy:  uuid.MustParse(req.UpdatedBy),
	}
}
