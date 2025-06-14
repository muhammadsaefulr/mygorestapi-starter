package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"{{.ModulePath}}/internal/domain/dto/{{.Name}}/request"
	model "{{.ModulePath}}/internal/domain/model/{{.Name}}"
	"{{.ModulePath}}/internal/repository/{{.Name}}"
	"{{.ModulePath}}/internal/shared/convert_types"
	"{{.ModulePath}}/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type {{.Name}}Service struct {
	Log        *logrus.Logger
	Validate   *validator.Validate
	Repo repository.{{.PascalName}}Repo
}

func New{{.PascalName}}Service(repo repository.{{.PascalName}}Repo, validate *validator.Validate) {{.PascalName}}Service {
	return {{.Name}}Service{
		Log:        utils.Log,
		Validate:   validate,
		Repo: repo,
	}
}

func (s *{{.Name}}Service) GetAll{{.PascalName}}(c *fiber.Ctx, params *request.Query{{.PascalName}}) ([]model.{{.PascalName}}, int64, error) {
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

func (s *{{.Name}}Service) Get{{.PascalName}}ByID(c *fiber.Ctx, id uint) (*model.{{.PascalName}}, error) {
	data, err := s.Repo.GetByID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "{{.PascalName}} not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get {{.Name}} by ID failed")
	}
	return data, nil
}

func (s *{{.Name}}Service) Create{{.PascalName}}(c *fiber.Ctx, req *request.Create{{.PascalName}}) (*model.{{.PascalName}}, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.Create{{.PascalName}}ToModel(req)
	if err := s.Repo.Create(c.Context(), data); err != nil {
		s.Log.Errorf("Create error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create {{.Name}} failed")
	}
	return data, nil
}

func (s *{{.Name}}Service) Update{{.PascalName}}(c *fiber.Ctx, id uint, req *request.Update{{.PascalName}}) (*model.{{.PascalName}}, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	data := convert_types.Update{{.PascalName}}ToModel(req)
	data.ID = id
	if err := s.Repo.Update(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update {{.Name}} failed")
	}
	return s.GetByID(c, id)
}

func (s *{{.Name}}Service) Delete{{.PascalName}}(c *fiber.Ctx, id uint) error {
	if _, err := s.Repo.GetByID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "{{.PascalName}} not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete {{.Name}} failed")
	}
	return nil
}
