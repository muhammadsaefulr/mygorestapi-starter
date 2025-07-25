package service

import (
	"github.com/gofiber/fiber/v2"
)

type NotificationServiceInterface interface {
	SendNotificationToUser(c *fiber.Ctx, userID string, title string, body string) error
	BroadcastToTopic(c *fiber.Ctx, topic string, title string, body string, data map[string]string) error
}
