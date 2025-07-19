package postgres

import (
	"my-go-project/internal/domain"
	"strings"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) GetByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetAll() ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *ProductRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *ProductRepository) GetByCategory(category string) ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.Where("category = ?", category).Find(&products).Error
	return products, err
}

func (r *ProductRepository) SearchByName(name string) ([]*domain.Product, error) {
	var products []*domain.Product
	searchTerm := "%" + strings.ToLower(name) + "%"
	err := r.db.Where("LOWER(name) LIKE ?", searchTerm).Find(&products).Error
	return products, err
}
