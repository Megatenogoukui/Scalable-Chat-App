package utils

import (
	"Backend/Model/Message"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProduceMessage(message Message.KafkaMessage) {
	dialer := KafkaAuthentication()

	// Define Kafka writer configuration
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka-2b529cca-chat-app-megate.e.aivencloud.com:11789"},
		Topic:   "Messages",
		Dialer:  dialer,
	})
	// Ensure the writer is closed when finished
	defer writer.Close()

	time1 := message.CreatedAt
	formattedTime := time1.Format("2006-01-02T15:04:05") // Format time1

	sender_id := message.SenderID
	value := message.TextMessage
	message_id := message.MessageId
	recipient_id := message.RecipientID

	fmt.Println("message_id :  ", message_id)

	message_to_send := bson.M{
		"message":        value,
		"sender_id":      sender_id,
		"recipient_id":   recipient_id,
		"conversationId": message.ConversationID,
		"time":           formattedTime,
		"messageId":      message_id,
		"connection":     message.ConnectionID,
	}

	// Convert the message to a JSON string
	messageJSON, _ := json.Marshal(message_to_send)

	// Message to send
	messageValue := value

	// Write the message
	err2 := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(messageJSON),
			Value: []byte(messageValue),
		},
	)
	if err2 != nil {
		log.Fatalf("Failed to write messages: %s", err2)
	}
	log.Println("Messages were successfully written to Kafka!")

}

// ConsumeMessages consumes messages from Kafka and sends them to MongoDB
func ConsumeMessages() {
	for {
		dialer := KafkaAuthentication()

		// Define Kafka reader configuration
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{"kafka-2b529cca-chat-app-megate.e.aivencloud.com:11789"},
			Topic:     "Messages",
			Dialer:    dialer,
			Partition: 0,    // Set the partition to 0 (or another partition if needed)
			MinBytes:  1,    // Minimum bytes to read
			MaxBytes:  10e6, // Maximum bytes to read
		})

		// Ensure the reader is closed when finished
		defer reader.Close()

		var messagesToInsert []Message.Message

		// MongoDB collection
		collM := GetCollection("messages")

		// Consume messages
		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Failed to read message: %s", err)
				continue
			}

			log.Printf("Received message: %s %s", string(msg.Key), string(msg.Value))

			// Process the message
			key := string(msg.Key)
			value := string(msg.Value)

			// Parse the key
			var keybson bson.M
			_s := json.Unmarshal([]byte(key), &keybson)
			if _s != nil {
				fmt.Printf("Failed to unmarshal key: %s", _s)
				continue
			}

			fmt.Println("Keybson: ", keybson)

			timestamp, _ := keybson["time"].(string)
			message_id_str, _ := keybson["messageId"].(string)
			message_id, _ := primitive.ObjectIDFromHex(message_id_str)
			fmt.Println("Message ID : ", message_id)
			sender_id_str, _ := keybson["sender_id"].(string)
			sender_id, _ := primitive.ObjectIDFromHex(sender_id_str)
			recipient_id_str, _ := keybson["recipient_id"].(string)
			recipient_id, _ := primitive.ObjectIDFromHex(recipient_id_str)
			conversation_id_str, _ := keybson["conversationId"].(string)
			conversation_id, _ := primitive.ObjectIDFromHex(conversation_id_str)
			connection_id, _ := keybson["connection"].(string)

			// Parse the timestamp
			created_at, err := time.Parse("2006-01-02T15:04:05", timestamp)
			if err != nil {
				log.Printf("Failed to parse created_at: %s", err)
				continue
			}

			// Check for duplicates in MongoDB
			filter := bson.M{"message_id": message_id} // Assuming `message_id` is the unique identifier
			var existingMessage bson.M
			err = collM.FindOne(context.Background(), filter).Decode(&existingMessage)
			if err == nil {
				// Message already exists, skip it
				log.Printf("Duplicate message found, skipping: %s", message_id)
				continue

			} else if err != mongo.ErrNoDocuments {
				// Unexpected error, log and continue
				log.Printf("Error checking for duplicate message: %v", err)
				continue
			}

			// Add the message to the batch to be inserted
			data_to_insert := Message.Message{
				TextMessage:    value,
				ClientId:       sender_id,
				CreatedAt:      created_at,
				MessageId:      message_id,      // Ensure you include MessageId in your struct
				ConversationID: conversation_id, // Add the conversation ID
				ConnectionID:   connection_id,
				SenderID:       sender_id,
				RecipientID:    recipient_id,
			}
			messagesToInsert = append(messagesToInsert, data_to_insert)

			// Insert messages into MongoDB if we have enough for a batch
			if len(messagesToInsert) >= 1 {
				BatchInsert(messagesToInsert)          // Use your batch insert function
				messagesToInsert = []Message.Message{} // Clear the batch
			}
		}
	}
}
func DeleteMessagesFromKafka(messages []Message.Message) error {

	dialer := KafkaAuthentication()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka-2b529cca-chat-app-megate.e.aivencloud.com:11789"},
		Dialer:  dialer,
	})

	for _, message := range messages {
		messageKey := fmt.Sprintf(`message-%s-%s`, message.CreatedAt.Format("2006-01-02T15:04:05"), message.ClientId)
		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(messageKey), // The key used to delete messages from Kafka
				Topic: "Messages",
			},
		)
		if err != nil {
			return fmt.Errorf("failed to delete message: %w", err)
		}
	}

	writer.Close()
	return nil
}

func KafkaAuthentication() *kafka.Dialer {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	if username == "" || password == "" {
		log.Fatalf("Kafka credentials are not set in .env")
	}

	// Load the CA certificate if needed
	rootCA, err := ioutil.ReadFile("./Utils/ca.pem")
	if err != nil {
		log.Fatalf("Failed to load root CA certificate: %s", err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(rootCA)

	dialer := &kafka.Dialer{
		SASLMechanism: plain.Mechanism{
			Username: username,
			Password: password,
		},
		TLS: &tls.Config{
			RootCAs: certPool,
		},
	}

	return dialer
}
