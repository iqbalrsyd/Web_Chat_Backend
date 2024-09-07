package internal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Settings struct {
	Notifications bool   `bson:"notifications"`
	Privacy       string `bson:"privacy"`
	Image         string `bson:"image"`
	Status        string `bson:"status"`
}

// UpdateSettings updates user settings in the database
func UpdateSettings(db *Database, userID primitive.ObjectID, newSettings Settings) error {
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"settings": newSettings}})
	return err
}

// GetSettings retrieves user settings from the database
func GetSettings(db *Database, userID primitive.ObjectID) (*Settings, error) {
	var user User
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user.Settings, nil
}

// CreateSetting initializes default settings for new users
func CreateSetting(db *Database, userID primitive.ObjectID, settings Settings) error {
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"settings": settings}})
	return err
}

// DeleteSetting resets user settings to default
func DeleteSetting(db *Database, userID primitive.ObjectID) error {
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defaultSettings := Settings{}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"settings": defaultSettings}})
	return err
}

// ViewSettings retrieves and returns user settings
func ViewSettings(db *Database, userID primitive.ObjectID) (*Settings, error) {
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user.Settings, nil
}
