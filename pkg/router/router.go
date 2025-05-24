package router

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/config"

	"github.com/muhammadsaefulr/NimeStreamAPI/pkg/service"
	od_service "github.com/muhammadsaefulr/NimeStreamAPI/pkg/service/otakudesu_scrape"
	"github.com/muhammadsaefulr/NimeStreamAPI/pkg/shared/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
	validate := validation.Validator()

	healthCheckService := service.NewHealthCheckService(db)
	emailService := service.NewEmailService()
	userService := service.NewUserService(db, validate)
	tokenService := service.NewTokenService(db, validate, userService)
	authService := service.NewAuthService(db, validate, userService, tokenService)
	odServices := od_service.NewAnimeService()

	v1 := app.Group("/api/v1")

	HealthCheckRoutes(v1, healthCheckService)
	AuthRoutes(v1, authService, userService, tokenService, emailService)
	UserRoutes(v1, userService, tokenService)
	OdRoutes(v1, odServices)
	// TODO: add another routes here...

	if !config.IsProd {
		DocsRoutes(v1)
	}
}
