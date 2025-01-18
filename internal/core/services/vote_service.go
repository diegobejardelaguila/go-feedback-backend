package services

import (
	"context"
	"time"

	"github.com/yourusername/go-feedback-backend/internal/core/domain"
	"github.com/yourusername/go-feedback-backend/internal/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type voteService struct {
	voteRepo ports.VoteRepository
}

func NewVoteService(voteRepo ports.VoteRepository) ports.VoteService {
	return &voteService{voteRepo: voteRepo}
}

func (s *voteService) CreateVote(ctx context.Context, vote *domain.Vote) error {
	vote.CreatedAt = time.Now()
	return s.voteRepo.Create(ctx, vote)
}

func (s *voteService) DeleteVote(ctx context.Context, id primitive.ObjectID) error {
	return s.voteRepo.Delete(ctx, id)
}

