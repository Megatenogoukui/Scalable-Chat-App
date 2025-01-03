package Conversation

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Users         []string           `json:"users" bson:"users"` // Array of user IDs
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	LastMessageID primitive.ObjectID `json:"last_message_id" bson:"last_message_id,omitempty"`
}
