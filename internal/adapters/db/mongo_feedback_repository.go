package db

import (
	"context"

	"github.com/yourusername/go-feedback-backend/internal/core/domain"
	"github.com/yourusername/go-feedback-backend/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoFeedbackRepository struct {
	collection *mongo.Collection
}

func NewMongoFeedbackRepository(db *mongo.Database) ports.FeedbackRepository {
	return &mongoFeedbackRepository{
		collection: db.Collection("feedback"),
	}
}

func (r *mongoFeedbackRepository) Create(ctx context.Context, feedback *domain.Feedback) error {
	_, err := r.collection.InsertOne(ctx, feedback)
	return err
}

func (r *mongoFeedbackRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Feedback, error) {
	var feedback domain.Feedback
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&feedback)
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *mongoFeedbackRepository) List(ctx context.Context, limit, offset int) ([]*domain.Feedback, error) {
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var feedbacks []*domain.Feedback
	for cursor.Next(ctx) {
		var feedback domain.Feedback
		if err := cursor.Decode(&feedback); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, &feedback)
	}

	return feedbacks, nil
}

func (r *mongoFeedbackRepository) Update(ctx context.Context, feedback *domain.Feedback) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": feedback.ID}, bson.M{"$set": feedback})
	return err
}

func (r *mongoFeedbackRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

