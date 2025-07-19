package domain

import "time"

type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type ReviewRepository interface {
	Create(review *Review) error
	GetByProductID(productID uint) ([]*Review, error)
}
