package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name string
	Age  int
	City string
}

func main() {
	credential := options.Credential{
		Username: "root",
		Password: "root",
	}

	// MongoDB 클라이언트 설정
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetAuth(credential)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// 데이터베이스 및 컬렉션 선택
	collection := client.Database("sos").Collection("metadata")

	// Create
	person := Person{"Alice", 25, "New York"}
	insertResult, err := collection.InsertOne(context.TODO(), person)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// 반환된 ObjectID 검증
	objectID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Fatal("InsertedID is not of type ObjectID")
	}
	fmt.Println("Inserted a single document with ObjectID: ", objectID.Hex())

	// Read
	var result Person
	filter := bson.D{{"name", "Alice"}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)

	// Update
	update := bson.D{
		{"$set", bson.D{
			{"age", 26},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Delete
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the people collection\n", deleteResult.DeletedCount)
}
