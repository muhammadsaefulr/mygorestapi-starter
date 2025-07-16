package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/comment/request"
	dto_response "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/comment/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/comment_service"
)

type CommentController struct {
	CommentService service.CommentService
}

func NewCommentController(commentService service.CommentService) *CommentController {
	return &CommentController{
		CommentService: commentService,
	}
}

// @Tags         Comments
// @Summary      Create a comment
// @Description  Create a comment. For replying to a comment, fill `parent_id` with the parent comment's ID. If replying to a child comment (nested reply), just mention the user using @username — the frontend should convert it into a hyperlink pointing to /users/info/:username.
// @Security     BearerAuth
// @Produce      json
// @Param        request  body  request.CreateComment  true  "Request body"
// @Router       /comments [post]
func (co *CommentController) CreateComment(c *fiber.Ctx) error {
	var req request.CreateComment
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Common{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	if err := co.CommentService.CreateComment(c, &req); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.Common{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Create comment successfully",
	})
}

// @Tags         Comments
// @Summary      Get a comment by ID
// @Description  Get a comment by ID
// @Produce      json
// @Param        id  path  uint  true  "Comment ID"
// @Router       /comments/{id} [get]
func (co *CommentController) GetCommentByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	comment, err := co.CommentService.GetCommentByID(c, uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[dto_response.CommentResponse]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Get comment by ID successfully",
		Data:    *comment,
	})
}

// @Tags         Comments
// @Summary      Get comments by movie ID
// @Description  Get comments by movie ID
// @Produce      json
// @Param        movieEpsId  path  string  true  "Movie Eps ID"  Example(drstn-s4-episode-8-sub-indo)
// @Router       /comments/movie/{movieEpsId} [get]
func (co *CommentController) GetCommentsMovieId(c *fiber.Ctx) error {
	movieId := c.Params("movieId")
	comments, err := co.CommentService.GetCommentsMovieId(c, movieId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[[]dto_response.CommentResponse]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Get comments by movie ID successfully",
		Data:    comments,
	})
}

// @Tags         Comments
// @Summary      Update a comment
// @Description  Update a comment
// @Security BearerAuth
// @Produce      json
// @Param        id  path  uint  true  "Comment ID"
// @Param        request  body  request.UpdateComment  true  "Request body"
// @Router       /comments/{id} [put]
func (co *CommentController) UpdateComment(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id := idStr
	var req request.UpdateComment

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Common{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	comment, err := co.CommentService.UpdateComment(c, &req, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[dto_response.CommentResponse]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Update comment successfully",
		Data:    *comment,
	})
}

// @Tags         Comments
// @Summary      Delete a comment
// @Description  Delete a comment
// @Security BearerAuth
// @Produce      json
// @Param        id  path  uint  true  "Comment ID"
// @Router       /comments/{id} [delete]
func (co *CommentController) DeleteComment(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	if err := co.CommentService.DeleteComment(c, uint(id)); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Delete comment successfully",
	})
}

// @Tags         Comments
// @Summary      Like a comment
// @Description  Like a comment
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Comment id"
// @Router       /comments/{id}/like [post]
func (co *CommentController) LikeComment(c *fiber.Ctx) error {
	commentID := c.Params("id")

	cmntIdUint, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		return err
	}

	if err := co.CommentService.LikeComment(c, uint(cmntIdUint)); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Like comment successfully",
	})
}

// @Tags         Comments
// @Summary      Dislike a comment
// @Description  Dislike a comment
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Comment id"
// @Router       /comments/{id}/dislike [post]
func (co *CommentController) DislikeComment(c *fiber.Ctx) error {
	commentID := c.Params("id")

	cmntIdUint, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		return err
	}

	if err := co.CommentService.DislikeComment(c, uint(cmntIdUint)); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Dislike comment successfully",
	})
}
