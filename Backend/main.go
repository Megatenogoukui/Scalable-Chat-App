package main

import (
	"Backend/Routers"
	"fmt"

	utils "Backend/Utils"

	upgrade_to_websocket "Backend/Controllers"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redis_client *redis.Client

func main() {
	//Initializing Default gin router `instance`
	router := gin.Default()

	// Set up CORS middleware
	router.Use(corsMiddleware())

	//Setting up routes
	Routers.MapRoutes(router)

	//Initializing Redis
	utils.Initialize_redis()

	//Initializing Mongo
	utils.Init_mongo()

	//Consuming messages from kafka to mongodb
	go utils.ConsumeMessages()

	// Start a Goroutine for Redis subscription
	go upgrade_to_websocket.ListenToRedis()

	//Getting the  redis instance
	redis_client = utils.GetRedisClient()

	//Starting the server
	fmt.Println("Starting the server")
	if err := router.Run(); err != nil {
		fmt.Println("Error in starting the server")
	}
}

// corsMiddleware returns a gin middleware for handling CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Change '*' to your specific frontend domain if needed
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // To allow cookies

		// Allow OPTIONS method for pre-flight request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No content
			return
		}

		// Proceed to next middleware or handler
		c.Next()
	}
}
