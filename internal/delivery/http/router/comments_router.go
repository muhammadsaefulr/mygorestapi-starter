package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/comment_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/comment_service"
)

func CommentsRoutes(v1 fiber.Router, c service.CommentService) {
	commentController := controller.NewCommentController(c)

	group := v1.Group("/comments")

	group.Get("/movie/:movieId", commentController.GetCommentsMovieId)

	group.Post("/", m.Auth(), commentController.CreateComment)
	group.Put("/:id", m.Auth(), commentController.UpdateComment)
	group.Get("/:id", commentController.GetCommentByID)
	group.Delete("/:id", m.Auth(), commentController.DeleteComment)
}
