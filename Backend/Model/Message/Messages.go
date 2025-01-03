package Message

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	TextMessage    string             `bson:text_message`
	CreatedAt      time.Time          `bson:"created_at"`
	ClientId       primitive.ObjectID `bson:"client_id"`
	MessageId      primitive.ObjectID `bson:"message_id"`
	ConnectionID   string             `bson:"connection_id"`
	ConversationID primitive.ObjectID `bson:"conversation_id"`
	SenderID       primitive.ObjectID `bson:"sender_id"`
	RecipientID    primitive.ObjectID `bson:"recipient_id"`
}

type KafkaMessage struct {
	TextMessage    string    `bson:text_message`
	CreatedAt      time.Time `bson:"created_at"`
	ClientId       string    `bson:"client_id"`
	MessageId      string    `bson:"message_id"`
	ConnectionID   string    `bson:"connection_id"`
	ConversationID string    `bson:"conversation_id"`
	SenderID       string    `bson:"sender_id"`
	RecipientID    string    `bson:"recipient_id"`
}
