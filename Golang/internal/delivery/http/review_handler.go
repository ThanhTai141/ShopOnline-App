package http

import (
	"my-go-project/internal/domain"
	"my-go-project/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ReviewHandler struct {
	reviewUsecase *usecase.ReviewUsecase
}

func NewReviewHandler(app *fiber.App, uc *usecase.ReviewUsecase) {
	handler := &ReviewHandler{reviewUsecase: uc}
	app.Post("/v1/reviews", handler.CreateReview)
	app.Get("/v1/reviews/:productID", handler.GetReviewsByProductID)
}

func (h *ReviewHandler) CreateReview(c *fiber.Ctx) error {
	var review domain.Review
	if err := c.BodyParser(&review); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// Lấy userID từ context nếu có xác thực JWT
	if userID := c.Locals("user_id"); userID != nil {
		if id, ok := userID.(uint); ok {
			review.UserID = id
		}
	}
	if err := h.reviewUsecase.CreateReview(&review); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(review)
}

func (h *ReviewHandler) GetReviewsByProductID(c *fiber.Ctx) error {
	productIDStr := c.Params("productID")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	reviews, err := h.reviewUsecase.GetReviewsByProductID(uint(productID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(reviews)
}
