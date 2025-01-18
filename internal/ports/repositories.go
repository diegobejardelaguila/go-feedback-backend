package ports

import (
	"context"

	"github.com/yourusername/go-feedback-backend/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *domain.Feedback) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Feedback, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Feedback, error)
	Update(ctx context.Context, feedback *domain.Feedback) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type VoteRepository interface {
	Create(ctx context.Context, vote *domain.Vote) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetByFeedbackID(ctx context.Context, feedbackID primitive.ObjectID) ([]*domain.Vote, error)
}

