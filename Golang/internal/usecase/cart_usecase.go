package usecase

import (
	"errors"
	"my-go-project/internal/domain"
)

type CartUseCase struct {
	CartRepo    domain.CartRepository
	ProductRepo domain.ProductRepository
}

func NewCartUseCase(cartRepo domain.CartRepository, productRepo domain.ProductRepository) *CartUseCase {
	return &CartUseCase{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}
}

func (uc *CartUseCase) GetCart(userID uint) (*domain.Cart, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	cart, err := uc.CartRepo.GetByUserID(userID)
	if err != nil {
		// If cart doesn't exist, create one
		if err := uc.CartRepo.CreateCart(userID); err != nil {
			return nil, err
		}
		return uc.CartRepo.GetByUserID(userID)
	}

	return cart, nil
}

func (uc *CartUseCase) AddToCart(userID uint, productID uint, quantity int) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}
	if productID == 0 {
		return errors.New("invalid product ID")
	}
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// Check if product exists
	product, err := uc.ProductRepo.GetByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Check if product has enough stock
	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	// Get or create cart
	cart, err := uc.GetCart(userID)
	if err != nil {
		return err
	}

	return uc.CartRepo.AddItem(cart.ID, productID, quantity)
}

func (uc *CartUseCase) UpdateCartItemQuantity(userID uint, productID uint, quantity int) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}
	if productID == 0 {
		return errors.New("invalid product ID")
	}
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	// Check if product exists
	product, err := uc.ProductRepo.GetByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Check if product has enough stock
	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	// Get cart
	cart, err := uc.GetCart(userID)
	if err != nil {
		return err
	}

	if quantity == 0 {
		return uc.CartRepo.RemoveItem(cart.ID, productID)
	}

	return uc.CartRepo.UpdateItemQuantity(cart.ID, productID, quantity)
}

func (uc *CartUseCase) RemoveFromCart(userID uint, productID uint) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}
	if productID == 0 {
		return errors.New("invalid product ID")
	}

	cart, err := uc.GetCart(userID)
	if err != nil {
		return err
	}

	return uc.CartRepo.RemoveItem(cart.ID, productID)
}

func (uc *CartUseCase) ClearCart(userID uint) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}

	cart, err := uc.GetCart(userID)
	if err != nil {
		return err
	}

	return uc.CartRepo.ClearCart(cart.ID)
}
