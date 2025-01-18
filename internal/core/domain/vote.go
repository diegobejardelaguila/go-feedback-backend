package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vote struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FeedbackID primitive.ObjectID `bson:"feedback_id" json:"feedback_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

