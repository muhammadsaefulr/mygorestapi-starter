package service

import (
	"context"
	"log"

	"firebase.google.com/go/v4/messaging"
	"github.com/gofiber/fiber/v2"
)

type NotificationService struct {
	firebaseMessaging *messaging.Client
}

func NewNotificationService(firebaseMessaging *messaging.Client) NotificationServiceInterface {
	return &NotificationService{
		firebaseMessaging: firebaseMessaging,
	}
}

func (n *NotificationService) SendNotificationToUser(c *fiber.Ctx, userToken string, title string, body string) error {
	message := &messaging.Message{
		Token: userToken,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},

		Data: map[string]string{
			"title": title,
			"body":  body,
		},
	}

	response, err := n.firebaseMessaging.Send(context.Background(), message)
	if err != nil {
		log.Printf("error sending notification: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to send notification")
	}

	log.Printf("Successfully sent message: %s", response)
	return nil
}

func (n *NotificationService) BroadcastToTopic(c *fiber.Ctx, topic string, title string, body string, data map[string]string) error {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Topic: topic,
		Data:  data,
	}

	response, err := n.firebaseMessaging.Send(c.Context(), message)
	if err != nil {
		log.Printf("error sending broadcast notification to topic %s: %v", topic, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to send broadcast notification")
	}

	log.Printf("Successfully sent broadcast message to topic %s: %s", topic, response)
	return nil
}
