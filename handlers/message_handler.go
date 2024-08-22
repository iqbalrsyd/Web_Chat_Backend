package handlers

import (
	"chat-backend/database"
	"chat-backend/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var messageCollection *mongo.Collection = database.Client.Database("chatDB").Collection("messages")

func SendMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	_ = json.NewDecoder(r.Body).Decode(&message)

	message.CreatedAt = time.Now()

	_, err := messageCollection.InsertOne(context.TODO(), message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	cur, err := messageCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer cur.Close(context.TODO())

	var messages []models.Message
	for cur.Next(context.TODO()) {
		var message models.Message
		_ = cur.Decode(&message)
		messages = append(messages, message)
	}

	json.NewEncoder(w).Encode(messages)
}
