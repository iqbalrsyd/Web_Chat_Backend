package internal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// GetChatByID fetches a chat by its ID
func GetChatByID(db *Database, chatID primitive.ObjectID) (*Chat, error) {
	collection := db.Collection("chats")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var chat Chat
	err := collection.FindOne(ctx, bson.M{"_id": chatID}).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Chat not found
		}
		return nil, err
	}

	return &chat, nil
}
