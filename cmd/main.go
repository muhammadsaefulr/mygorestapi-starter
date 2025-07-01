package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	"github.com/muhammadsaefulr/NimeStreamAPI/docs"
	module "github.com/muhammadsaefulr/NimeStreamAPI/internal"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	database "github.com/muhammadsaefulr/NimeStreamAPI/internal/infrastructure/persistence"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"gorm.io/gorm"
)

// ketika mau upload ganti ke dev.msaepul.my.id
// @title						NimeStream API documentation
// @version						1.0.0
// @BasePath					/api/v1
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description					Example Value: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := setupFiberApp()
	db := setupDatabase()
	defer closeDatabase(db)
	setupModule(app, db)

	address := fmt.Sprintf("%s:%d", config.AppHost, config.AppPort)

	if config.IsProd {
		docs.SwaggerInfo.Host = "dev.msaepul.my.id"
	} else {
		docs.SwaggerInfo.Host = "localhost:8080"
	}

	// Start server and handle graceful shutdown
	serverErrors := make(chan error, 1)
	go startServer(app, address, serverErrors)
	handleGracefulShutdown(ctx, app, serverErrors)
}

func setupFiberApp() *fiber.App {
	app := fiber.New(config.FiberConfig(), fiber.Config{
		BodyLimit: 4 * 1024 * 1024 * 1024,
	})

	// Middleware setup
	app.Use("/v1/auth", middleware.LimiterConfig())
	app.Use(middleware.LoggerConfig())
	app.Use(helmet.New())
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(middleware.RecoverConfig())

	return app
}

func setupDatabase() *gorm.DB {
	db := database.Connect(config.DBHost, config.DBName)
	// Add any additional database setup if needed
	return db
}

func setupModule(app *fiber.App, db *gorm.DB) {
	module.InitModule(app, db)
	app.Use(utils.NotFoundHandler)
}

func startServer(app *fiber.App, address string, errs chan<- error) {
	if err := app.Listen(address); err != nil {
		errs <- fmt.Errorf("error starting server: %w", err)
	}
}

func closeDatabase(db *gorm.DB) {
	sqlDB, errDB := db.DB()
	if errDB != nil {
		utils.Log.Errorf("Error getting database instance: %v", errDB)
		return
	}

	if err := sqlDB.Close(); err != nil {
		utils.Log.Errorf("Error closing database connection: %v", err)
	} else {
		utils.Log.Info("Database connection closed successfully")
	}
}

func handleGracefulShutdown(ctx context.Context, app *fiber.App, serverErrors <-chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		utils.Log.Fatalf("Server error: %v", err)
	case <-quit:
		utils.Log.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			utils.Log.Fatalf("Error during server shutdown: %v", err)
		}
	case <-ctx.Done():
		utils.Log.Info("Server exiting due to context cancellation")
	}

	utils.Log.Info("Server exited")
}
