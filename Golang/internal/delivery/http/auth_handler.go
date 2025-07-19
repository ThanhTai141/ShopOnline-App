package http

import (
	"my-go-project/internal/common"
	"my-go-project/internal/domain"
	"my-go-project/internal/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type UserHandler struct {
	Usecase *usecase.UserUsecase
}

func NewUserHandler(app *fiber.App, uc *usecase.UserUsecase) {
	handler := &UserHandler{Usecase: uc}
	app.Post("/v1/users/login", limiter.New(limiter.Config{
		Max:        3,
		Expiration: 1 * time.Minute,
		LimitReached: func(_ *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "Too many login attempts. Please try again later.")
		},
	}), handler.Login)

	app.Post("/v1/users/register", handler.Create)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}

	if err := user.ValidateLogin(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := h.Usecase.CreateUser(&user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.RespondCreated(c, user.ID)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}

	if err := user.ValidateLogin(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	accessToken, refreshToken, err := h.Usecase.LoginUser(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return common.RespondseSuccess(c, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
