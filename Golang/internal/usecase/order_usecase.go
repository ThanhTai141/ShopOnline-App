package usecase

import (
	"errors"
	"my-go-project/internal/domain"
)

type OrderUseCase struct {
	OrderRepo   domain.OrderRepository
	CartRepo    domain.CartRepository
	ProductRepo domain.ProductRepository
}

func NewOrderUseCase(orderRepo domain.OrderRepository, cartRepo domain.CartRepository, productRepo domain.ProductRepository) *OrderUseCase {
	return &OrderUseCase{
		OrderRepo:   orderRepo,
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}
}

func (uc *OrderUseCase) CreateOrder(userID uint, shippingAddress string) (*domain.Order, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	if shippingAddress == "" {
		return nil, errors.New("shipping address is required")
	}

	// Get user's cart
	cart, err := uc.CartRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Calculate total amount and validate stock
	var totalAmount float64
	for _, item := range cart.Items {
		product, err := uc.ProductRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, errors.New("product not found")
		}

		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		totalAmount += product.Price * float64(item.Quantity)
	}

	// Create order
	order := &domain.Order{
		UserID:          userID,
		Status:          domain.OrderStatusPending,
		TotalAmount:     totalAmount,
		ShippingAddress: shippingAddress,
	}

	// Create order items
	for _, item := range cart.Items {
		product, err := uc.ProductRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, err
		}

		orderItem := domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		order.Items = append(order.Items, orderItem)
	}

	if err := uc.OrderRepo.Create(order); err != nil {
		return nil, err
	}

	// Clear cart after successful order creation
	if err := uc.CartRepo.ClearCart(cart.ID); err != nil {
		// Log error but don't fail the order creation
		// In production, you might want to handle this differently
	}

	return order, nil
}

func (uc *OrderUseCase) GetOrderByID(id uint) (*domain.Order, error) {
	if id == 0 {
		return nil, errors.New("invalid order ID")
	}

	return uc.OrderRepo.GetByID(id)
}

func (uc *OrderUseCase) GetOrdersByUserID(userID uint) ([]*domain.Order, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	return uc.OrderRepo.GetByUserID(userID)
}

func (uc *OrderUseCase) UpdateOrderStatus(id uint, status domain.OrderStatus) error {
	if id == 0 {
		return errors.New("invalid order ID")
	}

	// Validate status
	switch status {
	case domain.OrderStatusPending, domain.OrderStatusConfirmed, domain.OrderStatusShipped, domain.OrderStatusDelivered, domain.OrderStatusCancelled:
		// Valid status
	default:
		return errors.New("invalid order status")
	}

	return uc.OrderRepo.UpdateStatus(id, status)
}

func (uc *OrderUseCase) GetAllOrders() ([]*domain.Order, error) {
	return uc.OrderRepo.GetAll()
}
