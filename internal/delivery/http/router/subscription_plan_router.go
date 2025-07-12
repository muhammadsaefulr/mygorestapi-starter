package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/subscription_plan_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/subscription_plan_service"
)

func SubscriptionPlanRoutes(v1 fiber.Router, c service.SubscriptionPlanServiceInterface) {
	subscription_planController := controller.NewSubscriptionPlanController(c)

	group := v1.Group("/subscription/plans")

	group.Get("/", subscription_planController.GetAllSubscriptionPlan)
	group.Post("/", m.Auth("postSubsPlan"), subscription_planController.CreateSubscriptionPlan)
	group.Get("/:id", subscription_planController.GetSubscriptionPlanByID)
	group.Put("/:id", m.Auth("updateSubsPlan"), subscription_planController.UpdateSubscriptionPlan)
	group.Delete("/:id", m.Auth("deleteSubsPlan"), subscription_planController.DeleteSubscriptionPlan)
}
