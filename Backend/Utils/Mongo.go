package utils

import (
	"Backend/Model/Message"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongo_client *mongo.Client
var database_name *mongo.Database

func Init_mongo() {

	//my mongo url
	const uri = "mongodb+srv://abhisheknaikworkspace:Megate%4019102001@cluster0.xctmb.mongodb.net/"

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	mongo_client = client

	//Getting database name
	database_name = client.Database("Cluster0")

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.M{"ping": 1}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

}

func GetCollection(collection_name string) *mongo.Collection {
	if mongo_client == nil {
		fmt.Println("Mongo client is not initialized")
	}
	if database_name == nil {
		fmt.Println("Mongo database is not initialized")
	}

	//Getting collection
	return database_name.Collection(collection_name)
}

func BatchInsert(messages []Message.Message) {
	// Getting the collection
	var messagesCollection = GetCollection("messages")

	var mongoDocuments []interface{}
	for _, message := range messages {
		mongoDocuments = append(mongoDocuments, message)
	}

	// Insert batch into MongoDB
	insertResult, err := messagesCollection.InsertMany(context.Background(), mongoDocuments)
	if err != nil {
		log.Fatalf("Failed to insert messages to MongoDB: %s", err)
	}

	log.Printf("Inserted %d messages into MongoDB with IDs: %v", len(insertResult.InsertedIDs), insertResult.InsertedIDs)

}
