package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"my-go-project/internal/domain"
	"my-go-project/pkg/cache"
	"time"
)

type TaskUseCase struct {
	TaskRepository domain.TaskRepository
}

func NewTaskUseCase(taskRepo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{
		TaskRepository: taskRepo,
	}
}
func (uc *TaskUseCase) CreateTask(task *domain.Task) error {
	return uc.TaskRepository.Create(task)
}

func (uc *TaskUseCase) GetTaskByID(id uint) (*domain.Task, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("task:%d", id)

	cachedTask, err := cache.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var task domain.Task
		if err := json.Unmarshal([]byte(cachedTask), &task); err == nil {
			return &task, nil
		}
	}

	task, err := uc.TaskRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %w", err)
	}
	cache.RedisClient.Set(ctx, cacheKey, taskJSON, 10*time.Minute)

	return task, nil

}
func (uc *TaskUseCase) GetAllTasks() ([]domain.Task, error) {
	return uc.TaskRepository.GetAll()
}
func (uc *TaskUseCase) UpdateTask(task *domain.Task) error {
	return uc.TaskRepository.Update(task)
}
func (uc *TaskUseCase) DeleteTask(id uint) error {
	return uc.TaskRepository.Delete(id)
}

func (uc *TaskUseCase) DeleteAllTasks() ([]domain.Task, error) {
	tasks, err := uc.TaskRepository.DeleteAll()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
