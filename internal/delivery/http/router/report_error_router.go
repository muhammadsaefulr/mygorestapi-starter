package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/report_error_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/report_error_service"
)

func ReportErrorRoutes(v1 fiber.Router, c service.ReportErrorServiceInterface) {
	report_errorController := controller.NewReportErrorController(c)

	group := v1.Group("/report-episode")

	group.Get("/", m.Auth("getAllReportError"), report_errorController.GetAllReportError)
	group.Post("/", m.Auth("postReportError"), report_errorController.CreateReportError)
	group.Get("/:id", m.Auth("getReportErrorByID"), report_errorController.GetReportErrorByID)
	group.Patch("/:id", m.Auth("updateReportError"), report_errorController.UpdateReportError)
	group.Delete("/:id", m.Auth("deleteReportError"), report_errorController.DeleteReportError)
}
