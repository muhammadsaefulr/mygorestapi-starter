package test

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	module "github.com/muhammadsaefulr/NimeStreamAPI/internal"
	database "github.com/muhammadsaefulr/NimeStreamAPI/internal/infrastructure/persistence"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var App = fiber.New(fiber.Config{
	CaseSensitive: true,
	ErrorHandler:  utils.ErrorHandler,
})
var DB *gorm.DB
var Redis = database.RedisClient
var Firebase = config.InitFirebaseAuthClient()
var Log = utils.Log
var firebaseMessaging = config.InitFirebaseMessagingClient()

func init() {
	// TODO: You can modify host and database configuration for tests
	DB = database.Connect("localhost", "testdb")
	module.InitModule(App, DB, Redis, Firebase, firebaseMessaging)
	App.Use(utils.NotFoundHandler)
}
