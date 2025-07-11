package controller

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"

	request "{{.ModulePath}}/internal/domain/dto/{{.Name}}/request"
	"{{.ModulePath}}/internal/domain/model"
	"{{.ModulePath}}/internal/domain/dto/util/response"
	"{{.ModulePath}}/internal/service/{{.Name}}_service"
)

type {{.PascalName}}Controller struct {
	Service service.{{.PascalName}}Service
}

func New{{.PascalName}}Controller(service service.{{.PascalName}}Service) *{{.PascalName}}Controller {
	return &{{.PascalName}}Controller{Service: service}
}

// @Tags         {{.Name}}
// @Summary      Get all {{.Name}}
// @Produce      json
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        search query     string  false  "Search term"
// @Router       /{{.Name}} [get]
func (h *{{.PascalName}}Controller) GetAll{{.PascalName}}(c *fiber.Ctx) error {
	query := &request.Query{{.PascalName}}{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
	}

	data, total, err := h.Service.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.{{.PascalName}}]{
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

// @Tags         {{.Name}}
// @Summary      Get a {{.Name}} by ID
// @Produce      json
// @Param        id  path  int  true  "{{.PascalName}} ID (uint)"
// @Router       /{{.Name}}/{id} [get]
func (h *{{.PascalName}}Controller) Get{{.PascalName}}ByID(c *fiber.Ctx) error {
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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.{{.PascalName}}]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         {{.Name}}
// @Summary      Create a new {{.Name}}
// @Accept       json
// @Produce      json
// @Param        request  body  request.Create{{.PascalName}}  true  "Request body"
// @Router       /{{.Name}} [post]
func (h *{{.PascalName}}Controller) Create{{.PascalName}}(c *fiber.Ctx) error {
	req := new(request.Create{{.PascalName}})
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.{{.PascalName}}]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
		Data:    *result,
	})
}

// @Tags         {{.Name}}
// @Summary      Update a {{.Name}}
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "{{.PascalName}} ID (uint)"
// @Param        request  body  request.Update{{.PascalName}}  true  "Request body"
// @Router       /{{.Name}}/{id} [put]
func (h *{{.PascalName}}Controller) Update{{.PascalName}}(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	req := new(request.Update{{.PascalName}})
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Update(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.{{.PascalName}}]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         {{.Name}}
// @Summary      Delete a {{.Name}}
// @Produce      json
// @Param        id  path  int  true  "{{.PascalName}} ID (uint)"
// @Router       /{{.Name}}/{id} [delete]
func (h *{{.PascalName}}Controller) Delete{{.PascalName}}(c *fiber.Ctx) error {
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