package convert_types

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/banner_app/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateBannerAppToModel(req *request.CreateBannerApp) *model.BannerApp {
	return &model.BannerApp{
		Title:      req.Title,
		BannerType: req.BannerType,
		ImageUrl:   req.ImageUrl,
		UpdatedBy:  req.UpdatedBy,
		DetailURL:  req.DetailURL,
	}
}

func UpdateBannerAppToModel(req *request.UpdateBannerApp) *model.BannerApp {
	return &model.BannerApp{
		ID:         req.ID,
		Title:      req.Title,
		BannerType: req.BannerType,
		UpdatedBy:  req.UpdatedBy,
		ImageUrl:   req.ImageUrl,
		DetailURL:  req.DetailURL,
	}
}
