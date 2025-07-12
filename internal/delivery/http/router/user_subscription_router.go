package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/user_subscription_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_subscription_service"
)

func UserSubscriptionRoutes(v1 fiber.Router, c service.UserSubscriptionServiceInterface) {
	user_subscriptionController := controller.NewUserSubscriptionController(c)

	group := v1.Group("/user/subscriptions")

	group.Get("/", m.Auth("userSubscriptionGet"), user_subscriptionController.GetAllUserSubscription)
	group.Post("/", m.Auth("userSubscriptionPost"), user_subscriptionController.CreateUserSubscription)
	group.Get("/:user_id", m.Auth("userSubscriptionGet"), user_subscriptionController.GetUserSubscriptionByID)
	group.Put("/:user_id", m.Auth("userSubscriptionPut"), user_subscriptionController.UpdateUserSubscription)
	group.Delete("/:user_id", m.Auth("userSubscriptionDelete"), user_subscriptionController.DeleteUserSubscription)
}
