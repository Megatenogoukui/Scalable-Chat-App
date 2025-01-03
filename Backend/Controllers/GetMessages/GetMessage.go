package get_messages

import (
	upgrade_to_websocket "Backend/Controllers"
	utils "Backend/Utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMessage(c *gin.Context) {
	channel_name := c.Query("channel_name")
	user_id1 := c.Query("userId1")
	user_id2 := c.Query("userId2")

	// Take initial 50 messages from Redis
	recent_messages := utils.Get_Redis_Message(channel_name)

	conversation, err2 := upgrade_to_websocket.FindConversation(user_id2, user_id1)
	if err2 != nil {
		c.JSON(400, gin.H{"message": "Error finding conversation"})
		return
	}

	var conversation_id_obj primitive.ObjectID
	if err2 == nil || conversation != nil {
		fmt.Println("Conversation found: ", conversation)
		conversation_id_obj, _ = conversation["_id"].(primitive.ObjectID)

	}

	recent_messages_bson, err := stringArrayToBsonArray(recent_messages, conversation_id_obj)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error in converting string array to bson array"})
		return
	}

	fmt.Println(len(recent_messages_bson))

	// Get the latest message time from the last message in `recent_messages`
	var latest_message primitive.M
	var latest_time_str string
	var latest_time time.Time
	if len(recent_messages_bson) > 0 {
		latest_message = recent_messages_bson[len(recent_messages_bson)-1]
		latest_time_str = latest_message["time"].(string)
		layout := "2006-01-02T15:04:05.999999999-07:00" // Go reference layout

		latest_time, err = time.Parse(layout, latest_time_str)
		if err != nil {
			fmt.Println("Error parsing time:", err)
		}
	}

	// Take remaining messages from MongoDB
	collM := utils.GetCollection("messages")
	filter := bson.M{
		"created_at": bson.M{
			"$lt": latest_time,
		},
		"conversation_id": conversation_id_obj,
	}

	sortOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}) // Descending order
	cursor, err := collM.Find(context.Background(), filter, sortOptions)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error in finding data from MongoDB"})
		return
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		fmt.Println("Inside curosoe")
		var message bson.M
		err := cursor.Decode(&message)
		if err != nil {
			fmt.Printf("Error decoding MongoDB message: %v", err)
			continue
		}

		recent_messages_bson = append(recent_messages_bson, message)
	}

	// Reverse the order of messages
	for i, j := 0, len(recent_messages_bson)-1; i < j; i, j = i+1, j-1 {
		recent_messages_bson[i], recent_messages_bson[j] = recent_messages_bson[j], recent_messages_bson[i]
	}

	// // Extract only the first part of each message in `recent_messages`
	// var final_messages []string
	// for _, msg := range recent_messages {
	// 	msg_parts := strings.Split(msg, "|")
	// 	if len(msg_parts) > 0 {
	// 		final_messages = append(final_messages, msg_parts[0]) // Add only the first part
	// 	}
	// }

	c.JSON(200, gin.H{"messages": recent_messages_bson, "count": len(recent_messages_bson)})
}

// Convert an array of strings (each string being a JSON object) to an array of bson.M
func stringArrayToBsonArray(jsonStrings []string, conversation_id primitive.ObjectID) ([]bson.M, error) {
	var bsonArray []bson.M
	for _, jsonString := range jsonStrings {
		var bsonObj bson.M
		err := json.Unmarshal([]byte(jsonString), &bsonObj) // Convert JSON string to bson.M
		if err != nil {
			return nil, err
		}

		conversation_id2, _ := bsonObj["conversationId"].(string)

		fmt.Println("Conversation ID1: ", conversation_id)
		fmt.Println("Conversation ID2: ", conversation_id2)
		if conversation_id2 != conversation_id.Hex() {
			continue
		}

		bsonArray = append(bsonArray, bsonObj)
	}
	return bsonArray, nil
}
