package convert_types

import (
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_badge/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateUserBadgeToModel(req *request.CreateUserBadge) *model.UserBadge {
	return &model.UserBadge{
		BadgeName: req.BadgeName,
		IconURL:   req.IconURL,
		Color:     req.Color,
	}
}

func UpdateUserBadgeToModel(req *request.UpdateUserBadge) *model.UserBadge {
	return &model.UserBadge{
		ID:        req.ID,
		BadgeName: req.BadgeName,
		IconURL:   req.IconURL,
		Color:     req.Color,
	}
}

func CreateUserInfoBadgeToModel(req *request.CreateUserBadgeInfo) *model.UserBadgeInfo {
	return &model.UserBadgeInfo{
		UserID:    uuid.MustParse(req.UserID),
		BadgeID:   req.BadgeID,
		Note:      req.Note,
		HandledBy: uuid.MustParse(req.HandledBy),
	}
}

func UpdateUserInfoBadgeToModel(req *request.UpdateUserBadgeInfo) *model.UserBadgeInfo {
	return &model.UserBadgeInfo{
		UserID:    uuid.MustParse(req.UserID),
		BadgeID:   req.BadgeID,
		Note:      req.Note,
		HandledBy: uuid.MustParse(req.HandledBy),
	}
}
