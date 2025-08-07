package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"github.com/muhammadsaefulr/mygorestapi-starter/config"
	"github.com/muhammadsaefulr/mygorestapi-starter/docs"
	module "github.com/muhammadsaefulr/mygorestapi-starter/internal"
	"github.com/redis/go-redis/v9"

	"github.com/muhammadsaefulr/mygorestapi-starter/internal/delivery/middleware"
	database "github.com/muhammadsaefulr/mygorestapi-starter/internal/infrastructure/persistence"
	"github.com/muhammadsaefulr/mygorestapi-starter/internal/shared/utils"

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
	redis := setupRedis()
	firebase := setupFirebaseAuthClient()
	firebaseMessaging := setupFirebaseMessagingClient()

	defer closeDatabase(db)
	setupModule(app, db, redis, firebase, firebaseMessaging)

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

	baseConfig := config.FiberConfig()
	baseConfig.BodyLimit = 4 * 1024 * 1024 * 1024
	baseConfig.Prefork = false

	app := fiber.New(baseConfig)

	// Middleware setup
	app.Use("/v1/auth", middleware.LimiterConfig())
	app.Use(middleware.LoggerConfig())
	app.Use(helmet.New())
	app.Use(compress.New())

	origin := config.ClientFeHost
	if origin == "" {
		origin = "http://localhost:3000"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     origin,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(middleware.RecoverConfig())

	return app
}

func setupDatabase() *gorm.DB {
	db := database.Connect(config.DBHost, config.DBName)
	// Add any additional database setup if needed
	return db
}

func setupFirebaseAuthClient() *auth.Client {
	return config.InitFirebaseAuthClient()
}

func setupFirebaseMessagingClient() *messaging.Client {
	return config.InitFirebaseMessagingClient()
}

func setupRedis() *redis.Client {
	client := database.ConnectRedis()
	return client
}

func setupModule(app *fiber.App, db *gorm.DB, redis *redis.Client, firebase *auth.Client, firebaseMessaging *messaging.Client) {
	module.InitModule(app, db, redis, firebase, firebaseMessaging)
	app.Use(utils.NotFoundHandler)
}

func startServer(app *fiber.App, address string, errs chan<- error) {
	utils.Log.Infof("Starting server at %s", address)
	if err := app.Listen(address); err != nil {
		utils.Log.Errorf("Fiber Listen failed: %v", err)
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
