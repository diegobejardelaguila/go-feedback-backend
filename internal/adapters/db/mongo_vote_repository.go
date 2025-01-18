package db

import (
	"context"

	"github.com/yourusername/go-feedback-backend/internal/core/domain"
	"github.com/yourusername/go-feedback-backend/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoVoteRepository struct {
	collection *mongo.Collection
}

func NewMongoVoteRepository(db *mongo.Database) ports.VoteRepository {
	return &mongoVoteRepository{
		collection: db.Collection("votes"),
	}
}

func (r *mongoVoteRepository) Create(ctx context.Context, vote *domain.Vote) error {
	_, err := r.collection.InsertOne(ctx, vote)
	return err
}

func (r *mongoVoteRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mongoVoteRepository) GetByFeedbackID(ctx context.Context, feedbackID primitive.ObjectID) ([]*domain.Vote, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"feedback_id": feedbackID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var votes []*domain.Vote
	for cursor.Next(ctx) {
		var vote domain.Vote
		if err := cursor.Decode(&vote); err != nil {
			return nil, err
		}
		votes = append(votes, &vote)
	}

	return votes, nil
}

