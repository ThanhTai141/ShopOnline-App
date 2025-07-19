package postgres

import (
	"my-go-project/internal/domain"

	"gorm.io/gorm"
)

type ReviewPG struct {
	db *gorm.DB
}

func NewReviewPG(db *gorm.DB) *ReviewPG {
	return &ReviewPG{db: db}
}

func (r *ReviewPG) Create(review *domain.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewPG) GetByProductID(productID uint) ([]*domain.Review, error) {
	var reviews []*domain.Review
	err := r.db.Where("product_id = ?", productID).Order("created_at desc").Find(&reviews).Error
	return reviews, err
}
