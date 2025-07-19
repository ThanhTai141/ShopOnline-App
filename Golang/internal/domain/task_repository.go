package domain

type TaskRepository interface {
	Create(task *Task) error
	GetByID(id uint) (*Task, error)
	GetAll() ([]Task, error)
	Update(task *Task) error
	Delete(id uint) error
	DeleteAll() ([]Task, error)
}
