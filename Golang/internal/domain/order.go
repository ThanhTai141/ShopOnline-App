package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id" gorm:"not null"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Order struct {
	ID              uint        `json:"id" gorm:"primaryKey"`
	UserID          uint        `json:"user_id" gorm:"not null"`
	Status          OrderStatus `json:"status" gorm:"not null;default:'pending'"`
	TotalAmount     float64     `json:"total_amount" gorm:"not null"`
	Items           []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	ShippingAddress string      `json:"shipping_address"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id uint) (*Order, error)
	GetByUserID(userID uint) ([]*Order, error)
	UpdateStatus(id uint, status OrderStatus) error
	GetAll() ([]*Order, error)
}
