package signup_login_user

import (
	user "Backend/Model/User"
	utils "Backend/Utils"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func SignUpUser(c *gin.Context) {

	var body bson.M
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error in binding JSON data"})
		return
	}

	// Getting the user details from the frontend
	user_name := body["user_name"].(string)
	email := body["email"].(string)
	password := body["password"].(string)

	// Checking if the user already exists in the databas	e
	coll := utils.GetCollection("users")
	filter := bson.M{
		"email": email,
	}
	var result bson.M
	err2 := coll.FindOne(context.Background(), filter).Decode(&result)
	if err2 == nil {
		c.JSON(400, gin.H{"message": "User already exists"})
		return
	}

	password2 := []byte(password)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password2, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	// Creating a new user object
	user := user.User{
		UserName: user_name,
		Email:    email,
		Password: string(hashedPassword),
	}

	// Inserting the user object in the database
	inserted_id, err := coll.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error in inserting data in MongoDB"})
		return
	}

	token, _ := createToken(email, inserted_id.InsertedID.(primitive.ObjectID).Hex())

	// Sending the success message back to the frontend
	c.JSON(200, gin.H{"message": "User signed up successfully", "token": token})

}

var secretKey = []byte("secret-code")

func createToken(username string, inserted_id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"_id":      inserted_id,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	fmt.Println("id : ", bson.M{
		"username": username,
		"_id":      inserted_id,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
