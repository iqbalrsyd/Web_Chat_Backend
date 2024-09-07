package internal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ChatID    primitive.ObjectID `bson:"chat_id"`
	SenderID  primitive.ObjectID `bson:"sender_id"`
	Content   string             `bson:"content"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

// SendMessage sends a message to a chat (individual or group)
func SendMessage(db *Database, chatID, senderID primitive.ObjectID, content string) (*Message, error) {
	collection := db.Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message := Message{
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   content,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := collection.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}

	message.ID = result.InsertedID.(primitive.ObjectID)
	return &message, nil
}

// GetMessages retrieves all messages from a chat
func GetMessages(db *Database, chatID primitive.ObjectID) ([]Message, error) {
	collection := db.Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var messages []Message
	cursor, err := collection.Find(ctx, bson.M{"chat_id": chatID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var message Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
