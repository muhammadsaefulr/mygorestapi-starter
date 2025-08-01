package middleware

import (
	"strings"

	"github.com/muhammadsaefulr/NimeStreamAPI/config"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"

	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_service"

	"github.com/gofiber/fiber/v2"
)

var userService service.UserService

func InitAuthMiddleware(us service.UserService) {
	userService = us
}

func Auth(requiredRights ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("access_token")
		if token == "" {
			authHeader := c.Get("Authorization")
			token = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		}

		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
		}

		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
		}

		userID, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
		}

		user, err := userService.GetUserByID(c, userID)
		if err != nil || user == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
		}

		c.Locals("user", user)

		// log.Printf("user permission: %+v", user.UserRole.Permissions)

		var userRights []string
		for _, p := range user.UserRole.Permissions {
			userRights = append(userRights, p.PermissionName)
		}

		if len(requiredRights) > 0 {
			if (!hasAllRights(userRights, requiredRights)) && c.Params("userId") != userID {
				return fiber.NewError(fiber.StatusForbidden, "You don't have permission to access this resource")
			}
		}

		return c.Next()
	}
}

func hasAllRights(userRights, requiredRights []string) bool {
	rightSet := make(map[string]struct{}, len(userRights))
	for _, right := range userRights {
		rightSet[right] = struct{}{}
	}

	for _, right := range requiredRights {
		if _, exists := rightSet[right]; !exists {
			return false
		}
	}
	return true
}
