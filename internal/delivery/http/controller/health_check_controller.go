package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/dto/util/response"
	service "github.com/muhammadsaefulr/mygorestapi-starter/internal/service/system_service"
)

type HealthCheckController struct {
	HealthCheckService service.HealthCheckService
}

func NewHealthCheckController(healthCheckService service.HealthCheckService) *HealthCheckController {
	return &HealthCheckController{
		HealthCheckService: healthCheckService,
	}
}

func (h *HealthCheckController) addServiceStatus(
	serviceList *[]response.HealthCheck, name string, isUp bool, message *string,
) {
	status := "Up"

	if !isUp {
		status = "Down"
	}

	*serviceList = append(*serviceList, response.HealthCheck{
		Name:    name,
		Status:  status,
		IsUp:    isUp,
		Message: message,
	})
}
func (h *HealthCheckController) Check(c *fiber.Ctx) error {
	isHealthy := true   // All service
	coreHealthy := true // Postgre + Memory
	var serviceList []response.HealthCheck

	// PostgreSQL check (core)
	if err := h.HealthCheckService.GormCheck(); err != nil {
		isHealthy = false
		coreHealthy = false
		errMsg := err.Error()
		h.addServiceStatus(&serviceList, "Postgre", false, &errMsg)
	} else {
		h.addServiceStatus(&serviceList, "Postgre", true, nil)
	}

	// Memory check (core)
	if err := h.HealthCheckService.MemoryHeapCheck(); err != nil {
		isHealthy = false
		coreHealthy = false
		errMsg := err.Error()
		h.addServiceStatus(&serviceList, "Memory", false, &errMsg)
	} else {
		h.addServiceStatus(&serviceList, "Memory", true, nil)
	}

	// S3 / MinIO check (non-core)
	if err := h.HealthCheckService.S3Check(); err != nil {
		isHealthy = false
		errMsg := err.Error()
		h.addServiceStatus(&serviceList, "S3 Object Storage", false, &errMsg)
	} else {
		h.addServiceStatus(&serviceList, "S3 Object Storage", true, nil)
	}

	statusCode := fiber.StatusOK
	status := "success"

	// Return 500 if core service are down
	if !coreHealthy {
		statusCode = fiber.StatusInternalServerError
		status = "error"
	}

	return c.Status(statusCode).JSON(response.HealthCheckResponse{
		Status:    status,
		Message:   "Health check completed",
		Code:      statusCode,
		IsHealthy: isHealthy,
		Result:    serviceList,
	})
}
