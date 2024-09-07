package internal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Settings struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"` // Menghubungkan ke User
	Notifications bool               `bson:"notifications" json:"notifications"`
	Privacy       string             `bson:"privacy" json:"privacy"`   // Contoh: "public", "private"
	Theme         string             `bson:"theme" json:"theme"`       // Contoh: "dark", "light"
	Language      string             `bson:"language" json:"language"` // Contoh: "en", "id"
	LastUpdated   time.Time          `bson:"last_updated" json:"last_updated"`
}

// Fungsi untuk membuat pengaturan baru
func CreateSettings(db *mongo.Database, userID primitive.ObjectID) (*Settings, error) {
	collection := db.Collection("settings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newSettings := Settings{
		UserID:        userID,
		Notifications: true,     // Default notification on
		Privacy:       "public", // Default privacy public
		Theme:         "light",  // Default theme light
		Language:      "en",     // Default language English
		LastUpdated:   time.Now(),
	}

	_, err := collection.InsertOne(ctx, newSettings)
	if err != nil {
		return nil, err
	}

	return &newSettings, nil
}

// Fungsi untuk mendapatkan pengaturan berdasarkan user ID
func GetSettings(db *mongo.Database, userID primitive.ObjectID) (*Settings, error) {
	collection := db.Collection("settings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var settings Settings
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Tidak ada pengaturan ditemukan
		}
		return nil, err
	}

	return &settings, nil
}

// Fungsi untuk mengupdate pengaturan pengguna
func UpdateSettings(db *mongo.Database, userID primitive.ObjectID, updates bson.M) error {
	collection := db.Collection("settings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": updates})
	if err != nil {
		return err
	}

	return nil
}

// Fungsi untuk menghapus pengaturan pengguna (misalnya saat pengguna menghapus akun)
func DeleteSettings(db *mongo.Database, userID primitive.ObjectID) error {
	collection := db.Collection("settings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		return err
	}

	return nil
}
