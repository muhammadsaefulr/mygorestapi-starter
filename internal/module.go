package module

import (
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/mygorestapi-starter/config"
	"github.com/muhammadsaefulr/mygorestapi-starter/internal/delivery/http/router"
	"github.com/muhammadsaefulr/mygorestapi-starter/internal/delivery/middleware"
	userRepo "github.com/muhammadsaefulr/mygorestapi-starter/internal/repository/user"
	"github.com/redis/go-redis/v9"

	userService "github.com/muhammadsaefulr/mygorestapi-starter/internal/service/user_service"

	// User Auth And Role

	userRoleRepo "github.com/muhammadsaefulr/mygorestapi-starter/internal/repository/user_role"
	userRoleService "github.com/muhammadsaefulr/mygorestapi-starter/internal/service/user_role_service"

	rolePermissionRepo "github.com/muhammadsaefulr/mygorestapi-starter/internal/repository/role_permissions"
	authService "github.com/muhammadsaefulr/mygorestapi-starter/internal/service/auth_service"
	rolePermissionService "github.com/muhammadsaefulr/mygorestapi-starter/internal/service/role_permissions_service"

	systemService "github.com/muhammadsaefulr/mygorestapi-starter/internal/service/system_service"

	"github.com/muhammadsaefulr/mygorestapi-starter/internal/shared/utils"
	"github.com/muhammadsaefulr/mygorestapi-starter/internal/shared/validation"

	"gorm.io/gorm"
)

func InitModule(app *fiber.App, db *gorm.DB, redis *redis.Client, firebase *auth.Client, firebaseMessaging *messaging.Client) {
	validate := validation.Validator()

	uploader, err := utils.NewS3Uploader(
		"http://minio:9000", // Endpoint (MinIO or AWS)
		"admin", "4dm1n3rs", // Access key
		"https://dev.msaepul.my.id/minio", // Endpoint (MinIO or AWS)
	)

	if err != nil {
		utils.Log.Errorf("Failed to create S3 uploader: %v", err)
		return
	}

	// Init services
	// fcmService := fcmService.NewNotificationService(firebaseMessaging)
	userRepo := userRepo.NewUserRepositryImpl(db)
	userSvc := userService.NewUserService(userRepo, validate, firebase)

	// User Auth And Role

	userRoleRepo := userRoleRepo.NewUserRoleRepositoryImpl(db)
	userRoleSvc := userRoleService.NewUserRoleService(userRoleRepo, validate)

	rolePermissionRepo := rolePermissionRepo.NewRolePermissionsRepositoryImpl(db)
	rolePermissionSvc := rolePermissionService.NewRolePermissionsService(rolePermissionRepo, validate)

	tokenSvc := systemService.NewTokenService(db, validate, userSvc)
	authSvc := authService.NewAuthService(db, validate, userSvc, tokenSvc)
	emailSvc := systemService.NewEmailService()
	healthSvc := systemService.NewHealthCheckService(db, uploader)

	middleware.InitAuthMiddleware(userSvc)

	v1 := app.Group("/api/v1")

	router.AuthRoutes(v1, authSvc, userSvc, tokenSvc, emailSvc)
	router.UserRoutes(v1, userSvc, tokenSvc)
	router.HealthCheckRoutes(v1, healthSvc)
	router.DocsRoutes(v1)
	router.UserRoleRoutes(v1, userRoleSvc)
	router.RolePermissionsRoutes(v1, rolePermissionSvc)

	if !config.IsProd {
		v1.Get("/docs", func(c *fiber.Ctx) error {
			return c.SendString("API Docs here")
		})
	}
}
