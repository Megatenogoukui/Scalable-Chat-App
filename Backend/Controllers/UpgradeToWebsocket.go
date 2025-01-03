package upgrade_to_websocket

import (
	"Backend/Model/Message"
	utils "Backend/Utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// object to upgrade https connection to socket connection
var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for testing purposes
	},
}

// Creating an empty object to store the data of clients which are connected
type Client struct {
	Conn   *websocket.Conn
	UserID string
}

var clients = make(map[string]*Client) // Map userID (string) to Client

func Handle_socket_connection(c *gin.Context) {
	w := c.Writer
	r := c.Request

	// Upgrading my HTTP connection to a WebSocket connection
	connection, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection to WebSocket:", err)
		return // Exit if upgrade fails
	}

	// Expect user_id as a query parameter during connection
	userID := c.Query("userId")
	if userID == "" {
		fmt.Println("User ID is required")
		connection.Close()
		return
	}

	// Add the new connection to the clients map
	clients[userID] = &Client{
		Conn:   connection,
		UserID: userID,
	}

	// Printing when user connected
	fmt.Printf("User Connected, current count: %d\n", len(clients))

	defer func() {
		connection.Close()

		//Deleting the associated users from the map
		for userID, client := range clients {
			if client.Conn == connection {
				delete(clients, userID)
				fmt.Printf("User Disconnected: %s, current count: %d\n", userID, len(clients))
				break
			}
		}
	}()

	// Running an infinite loop to listen for messages
	for {

		//Reading the message from the connection
		_, data, err := connection.ReadMessage()
		if err != nil {
			// Find and delete the user associated with this connection
			for userID, client := range clients {
				if client.Conn == connection {
					delete(clients, userID)
					fmt.Printf("User Disconnected: %s, current count: %d\n", userID, len(clients))
					break
				}
			}
			fmt.Println("Error retrieving message:", err)
			break // Exit loop on error
		}

		//Unmarshaling the data from string to object
		var data_obj bson.M
		err3 := json.Unmarshal(data, &data_obj)
		if err3 != nil {
			fmt.Println("Error unmarshalling data:", err3)
			continue
		}

		//Extracting information from that unmarshaled object
		message := data_obj["message"].(string)
		current_time := time.Now()
		message_id := primitive.NewObjectID()
		message_id_str := message_id.Hex()
		sender_id := data_obj["sender_id"].(string) //User Id of the User to whom the message is received
		sender_id_obj, _ := primitive.ObjectIDFromHex(sender_id)
		recipient_id := data_obj["recipient_id"].(string) //User Id of the User to whom the message is to be sent
		recipient_id_obj, _ := primitive.ObjectIDFromHex(recipient_id)

		// Find or create the conversation
		conversation, err2 := FindConversation(sender_id, recipient_id)
		var conversation_id string
		var conversation_id_obj primitive.ObjectID

		//If conversation does not exist then create a new connection
		if err2 != nil || conversation["_id"] == nil {
			inserted_doc, _ := utils.GetCollection("conversations").InsertOne(context.TODO(), bson.M{
				"users":           []string{sender_id, recipient_id},
				"created_at":      time.Now(),
				"last_message_id": message_id,
			})
			conversation_id = inserted_doc.InsertedID.(primitive.ObjectID).Hex()
			conversation_id_obj, _ = primitive.ObjectIDFromHex(conversation_id)

		} else {
			conversation_id = conversation["_id"].(primitive.ObjectID).Hex()
			conversation_id_obj, _ = primitive.ObjectIDFromHex(conversation_id)
		}

		//Structuring data for redis
		data_sent_to_redis := bson.M{
			"message":        message,
			"sender_id":      sender_id,
			"recipient_id":   recipient_id,
			"conversationId": conversation_id,
			"time":           current_time,
			"messageId":      message_id_str,
			"connection":     fmt.Sprintf(`%p`, connection),
		}

		//Converting from objectid to string
		dataWithSender, _ := json.Marshal(data_sent_to_redis)

		// Send the message to Redis
		utils.Publish_Message("Messages", string(dataWithSender))

		// Insert message into Kafka
		data_to_insert := Message.KafkaMessage{
			TextMessage:    string(message),
			ClientId:       sender_id_obj.Hex(),
			SenderID:       sender_id_obj.Hex(),
			RecipientID:    recipient_id_obj.Hex(),
			CreatedAt:      time.Now(),
			MessageId:      message_id.Hex(),
			ConnectionID:   fmt.Sprintf("%p", connection),
			ConversationID: conversation_id_obj.Hex(),
		}

		//Sending data to kafka
		utils.ProduceMessage(data_to_insert)
	}
}

func broadcastMessage(data []byte) {
	// Unmarshal the data into a BSON object
	var dataObj bson.M
	err := json.Unmarshal(data, &dataObj)
	if err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return
	}

	// Extract necessary fields
	actualData := dataObj["message"].(string)
	senderConnection := dataObj["connection"].(string)
	targetUserID := dataObj["recipient_id"].(string) // The target user ID

	// Broadcast the message to the target user
	for userID, client := range clients {

		// Skip if the userID is not the target user
		if userID != targetUserID {
			fmt.Println("not the target")
			continue
		}

		// Ensure we don't send the message to the sender
		if fmt.Sprintf("%p", client.Conn) == senderConnection {
			fmt.Println("same user")
			continue
		}

		// Attempt to send the message
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte(actualData))
		if err != nil {
			fmt.Println("Error broadcasting to client:", err)
			client.Conn.Close()
			delete(clients, userID) // Remove the client if it fails
		}
	}
}

// ListenToRedis listens to Redis and broadcasts messages to WebSocket clients
func ListenToRedis() {
	subscription := utils.Subscribe_Message("Messages")
	defer subscription.Close()

	fmt.Println("Listening to Redis channel: Messages")

	for {
		// Receive message from Redis
		message, err := subscription.ReceiveMessage(context.Background())
		if err != nil {
			fmt.Println("Error receiving message from Redis:", err)
			continue
		}

		fmt.Println("Received message from Redis:", message.Payload)

		// Broadcast message to all WebSocket clients
		broadcastMessage([]byte(message.Payload))
	}
}

func FindConversation(userId1, userId2 string) (bson.M, error) {
	// Get the conversations collection
	collection := utils.GetCollection("conversations")

	// Define the filter to check if a conversation exists
	filter := bson.M{
		"users": bson.M{
			"$all": []string{userId1, userId2},
		},
	}

	fmt.Println("Filter: ", userId1, userId2)

	var conversation bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&conversation)
	fmt.Println("Conversation found2: ", conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No conversation found
			return nil, nil
		}
		// Some other error occurred
		return nil, err
	}

	// Return the found conversation
	return conversation, nil
}
