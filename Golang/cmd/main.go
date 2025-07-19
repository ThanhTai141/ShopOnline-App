package main

import (
	"my-go-project/internal/delivery/http"
	"my-go-project/internal/repository/postgres"
	"my-go-project/internal/usecase"
	"my-go-project/pkg/cache"
	config "my-go-project/pkg/database"

	_ "my-go-project/docs" // Import generated docs

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Shopping App API
// @version 1.0
// @description A shopping app API with authentication, products, cart, and orders
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	app := fiber.New()
	db := config.InitDB()
	cache.InitRedis()
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Task handlers
	taskRepo := postgres.NewTaskRepository(db)
	taskUC := usecase.NewTaskUseCase(taskRepo)

	http.NewTaskHandler(app, taskUC)
	http.GetAllTasksHandler(app, taskUC)
	http.GetTaskByIDHandler(app, taskUC)
	http.UpdateTaskHandler(app, taskUC)
	http.DeleteTaskHandler(app, taskUC)
	http.DeleteAllTasksHandler(app, taskUC)

	// User handlers
	userRepo := postgres.NewUserPostgresRepo(db)
	userUC := usecase.NewUserUsecase(userRepo)
	http.NewUserHandler(app, userUC)

	// Product handlers
	productRepo := postgres.NewProductRepository(db)
	productUC := usecase.NewProductUseCase(productRepo)
	http.NewProductHandler(app, productUC)

	// Cart handlers
	cartRepo := postgres.NewCartRepository(db)
	cartUC := usecase.NewCartUseCase(cartRepo, productRepo)
	http.NewCartHandler(app, cartUC)

	// Order handlers
	orderRepo := postgres.NewOrderRepository(db)
	orderUC := usecase.NewOrderUseCase(orderRepo, cartRepo, productRepo)
	http.NewOrderHandler(app, orderUC)

	// Review handlers
	reviewRepo := postgres.NewReviewPG(db)
	reviewUC := usecase.NewReviewUsecase(reviewRepo)
	http.NewReviewHandler(app, reviewUC)

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
