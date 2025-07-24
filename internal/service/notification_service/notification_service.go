package service

import (
	"github.com/gofiber/fiber/v2"
)

type NotificationServiceInterface interface {
	SendNotificationToUser(c *fiber.Ctx, userID string, title string, body string) error
}
