package services

import (
	"context"
	"time"

	"github.com/yourusername/go-feedback-backend/internal/core/domain"
	"github.com/yourusername/go-feedback-backend/internal/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type feedbackService struct {
	feedbackRepo ports.FeedbackRepository
}

func NewFeedbackService(feedbackRepo ports.FeedbackRepository) ports.FeedbackService {
	return &feedbackService{feedbackRepo: feedbackRepo}
}

func (s *feedbackService) CreateFeedback(ctx context.Context, feedback *domain.Feedback) error {
	feedback.CreatedAt = time.Now()
	feedback.UpdatedAt = time.Now()
	return s.feedbackRepo.Create(ctx, feedback)
}

func (s *feedbackService) GetFeedback(ctx context.Context, id primitive.ObjectID) (*domain.Feedback, error) {
	return s.feedbackRepo.GetByID(ctx, id)
}

func (s *feedbackService) ListFeedback(ctx context.Context, limit, offset int) ([]*domain.Feedback, error) {
	return s.feedbackRepo.List(ctx, limit, offset)
}

func (s *feedbackService) UpdateFeedback(ctx context.Context, feedback *domain.Feedback) error {
	feedback.UpdatedAt = time.Now()
	return s.feedbackRepo.Update(ctx, feedback)
}

func (s *feedbackService) DeleteFeedback(ctx context.Context, id primitive.ObjectID) error {
	return s.feedbackRepo.Delete(ctx, id)
}

