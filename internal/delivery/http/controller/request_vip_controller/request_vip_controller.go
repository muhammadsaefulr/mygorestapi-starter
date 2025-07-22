package controller

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_vip/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/request_vip_service"
)

type RequestVipController struct {
	Service service.RequestVipService
}

func NewRequestVipController(service service.RequestVipService) *RequestVipController {
	return &RequestVipController{
		Service: service,
	}
}

// @Tags         request_vip
// @Summary      Get all request_vip
// @Security     BearerAuth
// @Produce      json
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        search query     string  false  "Search term"
// @Router       /request_vip [get]
func (h *RequestVipController) GetAllRequestVip(c *fiber.Ctx) error {
	query := &request.QueryRequestVip{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 10),
	}

	data, total, err := h.Service.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.RequestVip]{
		Code:         fiber.StatusOK,
		Status:       "success",
		Message:      "Successfully retrieved data",
		Results:      data,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   int64(math.Ceil(float64(total) / float64(query.Limit))),
		TotalResults: total,
	})
}

// @Tags         request_vip
// @Summary      Get a request_vip by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "RequestVip ID (uint)"
// @Router       /request-vip/{id} [get]
func (h *RequestVipController) GetRequestVipByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	result, err := h.Service.GetByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.RequestVip]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         request_vip
// @Summary      Create a new request_vip
// @Accept       multipart/form-data
// @Security     BearerAuth
// @Produce      json
// @Param        payment_method formData string true "Payment method"
// @Param        atas_nama_tf   formData string true "Nama pengirim"
// @Param        email          formData string true "Email"
// @Param        bukti_tf       formData file   true "Bukti transfer (image)"
// @Router       /request_vip [post]
func (h *RequestVipController) CreateRequestVip(c *fiber.Ctx) error {
	req := new(request.CreateRequestVip)

	fileHeader, err := c.FormFile("bukti_tf")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bukti transfer tidak ditemukan")
	}

	session := c.Locals("user").(*model.User)
	req.UserId = session.ID.String()
	req.BuktiTf = fileHeader

	_, err = h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.Common{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
	})
}

// @Tags         request_vip
// @Summary      Update a request_vip
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "RequestVip ID (uint)"
// @Param        request  body  request.UpdateRequestVip  true  "Request body"
// @Router       /request-vip/{id} [put]
func (h *RequestVipController) UpdateRequestVip(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	req := new(request.UpdateRequestVip)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Update(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.RequestVip]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         request_vip
// @Summary      Delete a request_vip
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "RequestVip ID (uint)"
// @Router       /request-vip/{id} [delete]
func (h *RequestVipController) DeleteRequestVip(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	if err := h.Service.Delete(c, id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data deleted successfully",
	})
}
