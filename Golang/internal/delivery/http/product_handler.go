package http

import (
	"my-go-project/internal/common"
	"my-go-project/internal/domain"
	"my-go-project/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	usecase *usecase.ProductUseCase
}

func NewProductHandler(app *fiber.App, uc *usecase.ProductUseCase) {
	handler := &ProductHandler{usecase: uc}

	// Public routes
	app.Get("/v1/products", handler.GetAll)
	app.Get("/v1/products/:id", handler.GetByID)
	app.Get("/v1/products/category/:category", handler.GetByCategory)
	app.Get("/v1/products/search/:name", handler.SearchByName)

	// Protected routes (admin only)
	app.Post("/v1/products", common.AuthMiddleware, handler.Create)
	app.Put("/v1/products/:id", common.AuthMiddleware, handler.Update)
	app.Delete("/v1/products/:id", common.AuthMiddleware, handler.Delete)
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var product domain.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.usecase.CreateProduct(&product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	product, err := h.usecase.GetProductByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	products, err := h.usecase.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	product, err := h.usecase.GetProductByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.usecase.UpdateProduct(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	if err := h.usecase.DeleteProduct(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ProductHandler) GetByCategory(c *fiber.Ctx) error {
	category := c.Params("category")
	products, err := h.usecase.GetProductsByCategory(category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (h *ProductHandler) SearchByName(c *fiber.Ctx) error {
	name := c.Params("name")
	products, err := h.usecase.SearchProductsByName(name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}
