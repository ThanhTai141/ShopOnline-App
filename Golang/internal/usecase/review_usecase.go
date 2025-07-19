package usecase

import (
	"my-go-project/internal/domain"
)

type ReviewUsecase struct {
	reviewRepo domain.ReviewRepository
}

func NewReviewUsecase(reviewRepo domain.ReviewRepository) *ReviewUsecase {
	return &ReviewUsecase{reviewRepo: reviewRepo}
}

func (u *ReviewUsecase) CreateReview(review *domain.Review) error {
	return u.reviewRepo.Create(review)
}

func (u *ReviewUsecase) GetReviewsByProductID(productID uint) ([]*domain.Review, error) {
	return u.reviewRepo.GetByProductID(productID)
}
