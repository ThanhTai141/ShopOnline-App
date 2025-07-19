package http

import (
	"fmt"
	"my-go-project/internal/common"
	"my-go-project/internal/domain"
	"my-go-project/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	usecase *usecase.TaskUseCase
}

func NewTaskHandler(app *fiber.App, uc *usecase.TaskUseCase) {
	handler := &TaskHandler{usecase: uc}

	app.Post("/tasks", common.AuthMiddleware, handler.Create)
}

func GetAllTasksHandler(app *fiber.App, uc *usecase.TaskUseCase) {
	handler := &TaskHandler{usecase: uc}

	app.Get("/tasks", common.AuthMiddleware, handler.GetAll)
}
func GetTaskByIDHandler(app *fiber.App, uc *usecase.TaskUseCase) {
	handler := &TaskHandler{usecase: uc}

	app.Get("/tasks/:id", handler.GetByID)
}
func UpdateTaskHandler(app *fiber.App, uc *usecase.TaskUseCase) {
	handler := &TaskHandler{usecase: uc}

	app.Put("/tasks/:id", handler.Update)
}
func DeleteAllTasksHandler(app *fiber.App, uc *usecase.TaskUseCase) {
	handler := &TaskHandler{usecase: uc}

	app.Delete("/tasks", handler.DeleteAll)
}

// DeleteAll godoc
// @Summary Delete all tasks
// @Description Delete all tasks for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Task "List of deleted tasks"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /tasks [delete]
func (h *TaskHandler) DeleteAll(c *fiber.Ctx) error {
	tasks, err := h.usecase.DeleteAllTasks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete all tasks",
		})
	}
	return c.Status(fiber.StatusOK).JSON(tasks)
}

// GetAll godoc
// @Summary Get all tasks
// @Description Get all tasks for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json

// @Success 200 {object} []domain.Task "List of tasks"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /tasks [get]
func (h *TaskHandler) GetAll(c *fiber.Ctx) error {
	tasks, err := h.usecase.GetAllTasks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve tasks",
		})
	}
	return c.Status(fiber.StatusOK).JSON(tasks)

}

// Update godoc
// @Summary Update a task
// @Description Update a task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body domain.Task true "Task object"
// @Success 200 {object} domain.Task "Updated task"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /tasks/{id} [put]
func (h *TaskHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}
	task, err := h.usecase.GetTaskByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if err := h.usecase.UpdateTask(task); err != nil {
		fmt.Println("DEBUG: error =", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update task",
		})
	}

	return c.Status(fiber.StatusOK).JSON(task)

}

// DeleteTaskHandler godoc
// @Summary Delete a task
// @Description Delete a task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /tasks/{id} [delete]
func DeleteTaskHandler(app *fiber.App, uc *usecase.TaskUseCase) {
	handler := &TaskHandler{usecase: uc}

	app.Delete("/tasks/:id", handler.Delete)
}
func (h *TaskHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	if err := h.usecase.DeleteTask(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetByID godoc
// @Summary Get a task by ID
// @Description Get a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} domain.Task "Task object"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	task, err := h.usecase.GetTaskByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(task)
}

// Create godoc
// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body domain.Task true "Task object"
// @Success 201 {object} map[string]interface{} "Created task ID"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 422 {object} map[string]interface{} "Validation Error"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /tasks [post]
func (h *TaskHandler) Create(c *fiber.Ctx) error {
	var task domain.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.usecase.CreateTask(&task); err != nil {
		fmt.Println("DEBUG: error =", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create task",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}
