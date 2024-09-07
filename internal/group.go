package internal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	Name      string               `bson:"name"`
	Admins    []primitive.ObjectID `bson:"admins"`
	Members   []primitive.ObjectID `bson:"members"`
	CreatedAt primitive.DateTime   `bson:"created_at"`
}

// CreateGroup creates a new group chat

func CreateGroup(db *Database, name string, creatorID primitive.ObjectID, memberIDs []primitive.ObjectID) (*Group, error) {
	collection := db.Collection("groups")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	group := Group{
		Name:      name,
		Admins:    []primitive.ObjectID{creatorID},
		Members:   append(memberIDs, creatorID),
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := collection.InsertOne(ctx, group)
	if err != nil {
		return nil, err
	}

	group.ID = result.InsertedID.(primitive.ObjectID)
	return &group, nil
}

// AddMember adds a member to the group chat
func AddMember(db *Database, groupID, userID primitive.ObjectID) error {
	collection := db.Collection("groups")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"_id": groupID}, bson.M{
		"$addToSet": bson.M{"members": userID},
	})
	return err
}

// RemoveMember removes a member from the group chat
func RemoveMember(db *Database, groupID, userID primitive.ObjectID) error {
	collection := db.Collection("groups")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"_id": groupID}, bson.M{
		"$pull": bson.M{"members": userID},
	})
	return err
}

func DeleteGroup(db *Database, groupID primitive.ObjectID) error {
	collection := db.Collection("groups")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": groupID})
	return err
}

func LeaveGroup(db *Database, groupID, userID primitive.ObjectID) error {
	collection := db.Collection("groups")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"_id": groupID}, bson.M{"$pull": bson.M{"members": userID}})
	return err
}

func InviteFriendToGroup(db *Database, groupID, friendID primitive.ObjectID) error {
	collection := db.Collection("groups")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"_id": groupID}, bson.M{"$addToSet": bson.M{"members": friendID}})
	return err
}
