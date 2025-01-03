package users

import (
	upgrade_to_websocket "Backend/Controllers"
	utils "Backend/Utils"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllUsers(c *gin.Context) {

	coll_users := utils.GetCollection("users")
	user_id1 := c.Query("userId")

	filter := bson.M{}
	var users []bson.M

	cursor, err := coll_users.Find(context.Background(), filter)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error in finding data from MongoDB"})
		return
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user bson.M
		err := cursor.Decode(&user)
		if err != nil {
			fmt.Printf("Error decoding MongoDB message: %v", err)
			continue
		}

		user_id2 := user["_id"].(primitive.ObjectID).Hex()

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

		user["conversation_id"] = conversation_id_obj

		users = append(users, user)
	}

	if users == nil {
		users = []bson.M{}
	}

	c.JSON(200, gin.H{"message": "User fetched", "result": users})

}
