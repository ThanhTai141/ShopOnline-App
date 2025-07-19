package domain

import "time"

type CartItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CartID    uint      `json:"cart_id" gorm:"not null"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null;default:1"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Cart struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null;unique"`
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartRepository interface {
	GetByUserID(userID uint) (*Cart, error)
	AddItem(cartID uint, productID uint, quantity int) error
	UpdateItemQuantity(cartID uint, productID uint, quantity int) error
	RemoveItem(cartID uint, productID uint) error
	ClearCart(cartID uint) error
	CreateCart(userID uint) error
}
