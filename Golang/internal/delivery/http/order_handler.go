package http

import (
	"my-go-project/internal/common"
	"my-go-project/internal/domain"
	"my-go-project/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	usecase *usecase.OrderUseCase
}

func NewOrderHandler(app *fiber.App, uc *usecase.OrderUseCase) {
	handler := &OrderHandler{usecase: uc}

	// User routes (require authentication)
	app.Post("/v1/orders", common.AuthMiddleware, handler.CreateOrder)
	app.Get("/v1/orders", common.AuthMiddleware, handler.GetUserOrders)
	app.Get("/v1/orders/:id", common.AuthMiddleware, handler.GetOrderByID)

	// Admin routes (require authentication)
	app.Get("/v1/admin/orders", common.AuthMiddleware, handler.GetAllOrders)
	app.Put("/v1/admin/orders/:id/status", common.AuthMiddleware, handler.UpdateOrderStatus)
}

type CreateOrderRequest struct {
	ShippingAddress string `json:"shipping_address"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.ShippingAddress == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Shipping address is required",
		})
	}

	order, err := h.usecase.CreateOrder(userID, req.ShippingAddress)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	orders, err := h.usecase.GetOrdersByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve orders",
		})
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	orderID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	order, err := h.usecase.GetOrderByID(uint(orderID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	// Check if the order belongs to the authenticated user
	if order.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	return c.Status(fiber.StatusOK).JSON(order)
}

func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	orders, err := h.usecase.GetAllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve orders",
		})
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	orderID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	var req UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Convert string status to domain.OrderStatus
	var status domain.OrderStatus
	switch req.Status {
	case "pending":
		status = domain.OrderStatusPending
	case "confirmed":
		status = domain.OrderStatusConfirmed
	case "shipped":
		status = domain.OrderStatusShipped
	case "delivered":
		status = domain.OrderStatusDelivered
	case "cancelled":
		status = domain.OrderStatusCancelled
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order status",
		})
	}

	if err := h.usecase.UpdateOrderStatus(uint(orderID), status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order status updated successfully",
	})
}
