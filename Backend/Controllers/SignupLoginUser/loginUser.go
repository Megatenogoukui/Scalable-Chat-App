package signup_login_user

import (
	utils "Backend/Utils"
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(c *gin.Context) {

	var body bson.M
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error in binding JSON data"})
		return
	}

	// Getting the user details from the frontend
	email := body["email"].(string)
	password := body["password"].(string)

	// Checking if the user already exists in the databas	e
	coll := utils.GetCollection("users")
	filter := bson.M{
		"email": email,
	}
	var result bson.M
	err2 := coll.FindOne(context.Background(), filter).Decode(&result)
	if err2 != nil {
		c.JSON(400, gin.H{"message": "User already exists"})
		return
	}

	og_passsword := []byte(result["password"].(string))

	password2 := []byte(password)

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(og_passsword, password2)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid password"})
	}

	token, _ := createToken(email, result["_id"].(primitive.ObjectID).Hex())

	// Sending the success message back to the frontend
	c.JSON(200, gin.H{"message": "User signed up successfully", "token": token})

}
