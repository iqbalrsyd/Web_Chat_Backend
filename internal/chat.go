package internal

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Chat struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	Type      string               `bson:"type"` // "individual" or "group"
	Members   []primitive.ObjectID `bson:"members"`
	CreatedAt primitive.DateTime   `bson:"created_at"`
}

// CreateChat creates an individual or group chat
func CreateChat(db *Database, chatType string, members []primitive.ObjectID) (*Chat, error) {
	collection := db.Collection("chats")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chat := Chat{
		Type:      chatType,
		Members:   members,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := collection.InsertOne(ctx, chat)
	if err != nil {
		return nil, err
	}

	chat.ID = result.InsertedID.(primitive.ObjectID)
	return &chat, nil
}

func GetChatByID(db *Database, chatID primitive.ObjectID) (*Chat, error) {
	collection := db.Collection("chats")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var chat Chat
	err := collection.FindOne(ctx, bson.M{"_id": chatID}).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("No chat found with ID:", chatID.Hex()) // Debugging information
			return nil, nil                                     // Chat not found
		}
		log.Println("Error fetching chat:", err) // Log the error for debugging
		return nil, err
	}

	return &chat, nil
}
