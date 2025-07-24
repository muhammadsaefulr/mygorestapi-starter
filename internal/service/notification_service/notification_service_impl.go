package service

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"github.com/gofiber/fiber/v2"
	"log"
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
		return err
	}

	log.Printf("Successfully sent message: %s", response)
	return nil
}
