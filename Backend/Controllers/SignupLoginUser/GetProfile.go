package signup_login_user

import (
	utils "Backend/Utils"
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProfile(c *gin.Context) {

	// Getting the user details from the frontend
	user_id := c.Query("_id")

	user_id_hex, _ := primitive.ObjectIDFromHex(user_id)

	coll_users := utils.GetCollection("users")

	filter := bson.M{
		"_id": user_id_hex,
	}

	var result bson.M
	err := coll_users.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		c.JSON(400, gin.H{"message": "User does not exist"})
		return
	}

	c.JSON(200, gin.H{"message": "User fetched", "result": result})

}
