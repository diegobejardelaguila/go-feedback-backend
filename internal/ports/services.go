package ports

import (
	"context"

	"github.com/yourusername/go-feedback-backend/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	Register(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, email, password string) (string, error)
	GetUser(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	GenerateToken(user *domain.User) (string, error)
}

type FeedbackService interface {
	CreateFeedback(ctx context.Context, feedback *domain.Feedback) error
	GetFeedback(ctx context.Context, id primitive.ObjectID) (*domain.Feedback, error)
	ListFeedback(ctx context.Context, limit, offset int) ([]*domain.Feedback, error)
	UpdateFeedback(ctx context.Context, feedback *domain.Feedback) error
	DeleteFeedback(ctx context.Context, id primitive.ObjectID) error
}

type VoteService interface {
	CreateVote(ctx context.Context, vote *domain.Vote) error
	DeleteVote(ctx context.Context, id primitive.ObjectID) error
}

type AnalyticsService interface {
	GetAnalytics(ctx context.Context) (map[string]interface{}, error)
}

