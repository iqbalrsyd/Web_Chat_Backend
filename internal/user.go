package internal

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	Username     string               `bson:"username"`
	Password     string               `bson:"password"`
	PasswordHash string               `bson:"password_hash"`
	Firstname    string               `bson:"firstname"`
	Lastname     string               `bson:"lastname"`
	Birthdate    time.Time            `bson:"birthdate"`
	ProfileName  string               `bson:"profile_name"`
	EmailID      string               `bson:"email_id"`
	Status       string               `bson:"status"`
	PhoneNumber  string               `bson:"phone_number"`
	ProfilePic   string               `bson:"image"`
	BlockedUsers []primitive.ObjectID `bson:"blocked_users"`
	Settings     Settings             `bson:"settings"`
}

// RegisterUser registers a new user
func RegisterUser(db *Database, username, password string) (*User, error) {
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := User{
		Username:     username,
		PasswordHash: string(passwordHash),
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

// AuthenticateUser authenticates a user with their username and password
func AuthenticateUser(db *Database, username, password string) (*User, error) {
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil // Invalid password
	}

	return &user, nil
}

func LoginUser(db *Database, username, password string) (string, error) {
	var user User
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Validasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token, err := GenerateJWT(user.ID.Hex())
	if err != nil {
		return "", err
	}

	return token, nil
}

func BlockUser(db *Database, userID, blockUserID primitive.ObjectID) error {
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$addToSet": bson.M{"blocked_users": blockUserID}})
	return err
}

// SearchUsers mencari pengguna berdasarkan username atau email
func SearchUsers(db *Database, query string) ([]User, error) {
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"username": bson.M{"$regex": query, "$options": "i"}},
			{"email_id": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func GetActiveList(db *Database) ([]User, error) {
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"status": "online"}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
