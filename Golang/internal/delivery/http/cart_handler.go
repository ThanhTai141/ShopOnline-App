package http

import (
	"my-go-project/internal/common"
	"my-go-project/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	usecase *usecase.CartUseCase
}

func NewCartHandler(app *fiber.App, uc *usecase.CartUseCase) {
	handler := &CartHandler{usecase: uc}

	// All cart routes require authentication
	app.Get("/v1/cart", common.AuthMiddleware, handler.GetCart)
	app.Post("/v1/cart/items", common.AuthMiddleware, handler.AddToCart)
	app.Put("/v1/cart/items/:productId", common.AuthMiddleware, handler.UpdateCartItem)
	app.Delete("/v1/cart/items/:productId", common.AuthMiddleware, handler.RemoveFromCart)
	app.Delete("/v1/cart", common.AuthMiddleware, handler.ClearCart)
}

type AddToCartRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity"`
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	cart, err := h.usecase.GetCart(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get cart",
		})
	}

	return c.Status(fiber.StatusOK).JSON(cart)
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.ProductID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	if req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Quantity must be greater than 0",
		})
	}

	if err := h.usecase.AddToCart(userID, req.ProductID, req.Quantity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Item added to cart successfully",
	})
}

func (h *CartHandler) UpdateCartItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	productID, err := strconv.ParseUint(c.Params("productId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var req UpdateCartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Quantity < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Quantity cannot be negative",
		})
	}

	if err := h.usecase.UpdateCartItemQuantity(userID, uint(productID), req.Quantity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cart item updated successfully",
	})
}

func (h *CartHandler) RemoveFromCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	productID, err := strconv.ParseUint(c.Params("productId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	if err := h.usecase.RemoveFromCart(userID, uint(productID)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Item removed from cart successfully",
	})
}

func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	if err := h.usecase.ClearCart(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear cart",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cart cleared successfully",
	})
}
