package postgres

import (
	"my-go-project/internal/domain"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetByUserID(userID uint) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) CreateCart(userID uint) error {
	cart := &domain.Cart{
		UserID: userID,
	}
	return r.db.Create(cart).Error
}

func (r *CartRepository) AddItem(cartID uint, productID uint, quantity int) error {
	// Check if item already exists
	var existingItem domain.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&existingItem).Error

	if err == gorm.ErrRecordNotFound {
		// Create new item
		cartItem := &domain.CartItem{
			CartID:    cartID,
			ProductID: productID,
			Quantity:  quantity,
		}
		return r.db.Create(cartItem).Error
	} else if err != nil {
		return err
	}

	// Update existing item quantity
	existingItem.Quantity += quantity
	return r.db.Save(&existingItem).Error
}

func (r *CartRepository) UpdateItemQuantity(cartID uint, productID uint, quantity int) error {
	if quantity == 0 {
		return r.RemoveItem(cartID, productID)
	}

	return r.db.Model(&domain.CartItem{}).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Update("quantity", quantity).Error
}

func (r *CartRepository) RemoveItem(cartID uint, productID uint) error {
	return r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).
		Delete(&domain.CartItem{}).Error
}

func (r *CartRepository) ClearCart(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&domain.CartItem{}).Error
}
