package usecase

import (
	"errors"
	"my-go-project/internal/domain"
	"strings"
)

type ProductUseCase struct {
	Repo domain.ProductRepository
}

func NewProductUseCase(r domain.ProductRepository) *ProductUseCase {
	return &ProductUseCase{Repo: r}
}

func (uc *ProductUseCase) CreateProduct(p *domain.Product) error {
	if p.Name == "" {
		return errors.New("product name is required")
	}
	if p.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	if p.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}

	return uc.Repo.Create(p)
}

func (uc *ProductUseCase) GetProductByID(id uint) (*domain.Product, error) {
	if id == 0 {
		return nil, errors.New("invalid product ID")
	}

	return uc.Repo.GetByID(id)
}

func (uc *ProductUseCase) GetAllProducts() ([]*domain.Product, error) {
	return uc.Repo.GetAll()
}

func (uc *ProductUseCase) UpdateProduct(p *domain.Product) error {
	if p.ID == 0 {
		return errors.New("invalid product ID")
	}
	if p.Name == "" {
		return errors.New("product name is required")
	}
	if p.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	if p.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}

	return uc.Repo.Update(p)
}

func (uc *ProductUseCase) DeleteProduct(id uint) error {
	if id == 0 {
		return errors.New("invalid product ID")
	}

	return uc.Repo.Delete(id)
}

func (uc *ProductUseCase) GetProductsByCategory(category string) ([]*domain.Product, error) {
	if category == "" {
		return nil, errors.New("category is required")
	}

	return uc.Repo.GetByCategory(category)
}

func (uc *ProductUseCase) SearchProductsByName(name string) ([]*domain.Product, error) {
	if name == "" {
		return nil, errors.New("search term is required")
	}

	// Convert to lowercase for case-insensitive search
	searchTerm := strings.ToLower(strings.TrimSpace(name))
	return uc.Repo.SearchByName(searchTerm)
}
