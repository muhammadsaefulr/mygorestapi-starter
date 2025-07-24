package router

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/config"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	auth_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/auth_service"
	system_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/system_service"
	user_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(
	v1 fiber.Router, a auth_service.AuthService, u user_service.UserService,
	t system_service.TokenService, e system_service.EmailService,
) {
	authController := controller.NewAuthController(a, u, t, e)
	config.GoogleConfig()

	auth := v1.Group("/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)
	auth.Post("/refresh-tokens", authController.RefreshTokens)
	auth.Post("/forgot-password", authController.ForgotPassword)
	auth.Post("/reset-password", authController.ResetPassword)
	auth.Post("/send-verification-email", m.Auth(), authController.SendVerificationEmail)
	auth.Post("/verify-email", authController.VerifyEmail)
	auth.Get("/google", authController.GoogleLogin)
	auth.Post("/google/signin", authController.FirebaseGoogleSignIn)
	auth.Get("/google-callback", authController.GoogleCallback)
}
