package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// Creating a global variable
var redis_client *redis.Client

// sync.Once ensures that the code inside this runs only once
var (
	redisOnce sync.Once
)

// This function ensures that there is only one instance of redis client
func SetRedisClient(client *redis.Client) {
	redisOnce.Do(func() {
		redis_client = client
	})
}

// returns the redis client instance
func GetRedisClient() *redis.Client {
	return redis_client
}

func Initialize_redis() {

	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Get Redis URL from environment variables
	url := os.Getenv("REDIS_URL")
	if url == "" {
		log.Fatalf("Redis URL is not set in .env")
	}

	// Parse the Redis URL into options
	options, err1 := redis.ParseURL(url)
	if err1 != nil {
		panic(fmt.Sprintf("Error parsing Redis URL: %v", err1))
	}

	//Creating a new redis instance
	redis_client = redis.NewClient(options)

	//Testing the redis instance if it is connecting or not
	_, err2 := redis_client.Ping(context.Background()).Result()
	if err2 != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err2))
	}

	fmt.Println("Connected to Redis successfully with connection pooling")

	//Using this function so that only one instance of redis_client is made
	SetRedisClient(redis_client)

	fmt.Println("Redis Client : ", redis_client)
}

func Publish_Message(channel string, data string) {
	//Publishing message to redis
	err := redis_client.Publish(context.Background(), channel, data).Err()
	if err != nil {
		fmt.Println("Error while publishing : ", err)
	}

	fmt.Println("Data : ", data)
	//Add messages to redis list
	err2 := redis_client.LPush(context.Background(), fmt.Sprintf(`%s-list`, channel), data).Err()
	if err2 != nil {
		fmt.Println("Error while publishing : ", err)
	}

	//Trim the redis list so that only 50 messages are stored at a time
	err3 := redis_client.LTrim(context.Background(), fmt.Sprintf(`%s-list`, channel), 0, 9).Err()
	if err3 != nil {
		fmt.Println("Error while publishing : ", err)
	}
}

func Subscribe_Message(channel string) *redis.PubSub {
	//Subscribing to the channel
	pubsub := redis_client.Subscribe(context.Background(), channel)

	// Close the subscription when we are done.
	// defer pubsub.Close()

	return pubsub
}

func Get_Redis_Message(channel string) []string {
	// Get recent 50 messages from Redis
	recentMessages, err := redis_client.LRange(context.Background(), fmt.Sprintf("%s-list", channel), 0, -1).Result()
	if err != nil {
		log.Printf("Error fetching from Redis list: %v", err)
	}

	return recentMessages
}
