package module

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/router"
	userRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/user"
	authService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/auth_service"
	odService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/otakudesu_scrape"
	systemService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/system_service"
	userService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_service"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/validation"

	"gorm.io/gorm"
)

func InitModule(app *fiber.App, db *gorm.DB) {
	validate := validation.Validator()

	// Init services
	userRepo := userRepo.NewUserRepositryImpl(db)
	userSvc := userService.NewUserService(userRepo, validate)

	tokenSvc := systemService.NewTokenService(db, validate, userSvc)

	authSvc := authService.NewAuthService(db, validate, userSvc, tokenSvc)

	emailSvc := systemService.NewEmailService()
	healthSvc := systemService.NewHealthCheckService(db)

	animeSvc := odService.NewAnimeService()

	v1 := app.Group("/api/v1")

	router.AuthRoutes(v1, authSvc, userSvc, tokenSvc, emailSvc)
	router.UserRoutes(v1, userSvc, tokenSvc)
	router.OdRoutes(v1, animeSvc)
	router.HealthCheckRoutes(v1, healthSvc)
	router.DocsRoutes(v1)

	if !config.IsProd {
		v1.Get("/docs", func(c *fiber.Ctx) error {
			return c.SendString("API Docs here")
		})
	}
}
