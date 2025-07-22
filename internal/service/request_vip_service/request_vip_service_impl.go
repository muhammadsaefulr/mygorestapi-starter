package service

import (
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_vip/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/request_vip"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type RequestVipService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     repository.RequestVipRepo
	S3       *utils.S3Uploader
}

func NewRequestVipService(repo repository.RequestVipRepo, validate *validator.Validate, S3 *utils.S3Uploader) RequestVipService {
	return RequestVipService{
		Log:      utils.Log,
		Validate: validate,
		Repo:     repo,
		S3:       S3,
	}
}

func (s *RequestVipService) GetAll(c *fiber.Ctx, params *request.QueryRequestVip) ([]model.RequestVip, int64, error) {
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

func (s *RequestVipService) GetByID(c *fiber.Ctx, id uint) (*model.RequestVip, error) {
	data, err := s.Repo.GetByID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "RequestVip not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get request_vip by ID failed")
	}
	return data, nil
}

func (s *RequestVipService) Create(c *fiber.Ctx, req *request.CreateRequestVip) (*model.RequestVip, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Validasi gagal: "+err.Error())
	}

	if req.BuktiTf == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "File bukti transfer tidak ditemukan")
	}

	file, err := req.BuktiTf.Open()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal membuka file bukti transfer")
	}
	defer file.Close()

	fileName := fmt.Sprintf("bukti-tf/%d_%s", time.Now().Unix(), req.BuktiTf.Filename)
	contentType := req.BuktiTf.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	log.Printf("Form data info: %+v", req)

	filePath, _, err := s.S3.UploadFile("bukti-transfer", file, fileName, contentType)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal upload file ke MinIO")
	}

	req.BuktiTfStr = filePath

	data := convert_types.CreateRequestVipToModel(req)

	if err := s.Repo.Create(c.Context(), data); err != nil {
		s.Log.Errorf("Create error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create request_vip failed")
	}

	return data, nil
}

func (s *RequestVipService) Update(c *fiber.Ctx, id uint, req *request.UpdateRequestVip) (*model.RequestVip, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.UpdateRequestVipToModel(req)
	data.ID = id
	if err := s.Repo.Update(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update request_vip failed")
	}
	return s.GetByID(c, id)
}

func (s *RequestVipService) Delete(c *fiber.Ctx, id uint) error {
	if _, err := s.Repo.GetByID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "RequestVip not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete request_vip failed")
	}
	return nil
}
