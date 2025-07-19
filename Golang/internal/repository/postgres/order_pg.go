package postgres

import (
	"my-go-project/internal/domain"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) GetByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items.Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByUserID(userID uint) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateStatus(id uint, status domain.OrderStatus) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) GetAll() ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.Preload("Items.Product").Find(&orders).Error
	return orders, err
}
