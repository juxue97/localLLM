package user

import (
	"context"
	"fmt"
	"time"

	"chatbot/config"
	"chatbot/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	mongoClient *mongo.Client
}

func NewStore(mongoClient *mongo.Client) *Store {
	return &Store{mongoClient: mongoClient}
}

func (s *Store) CreateUser(user types.User) (primitive.ObjectID, error) {
	collection := s.mongoClient.Database(config.Envs.MongoDatabase).Collection("users")

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	userID := result.InsertedID.(primitive.ObjectID)
	return userID, nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	collection := s.mongoClient.Database(config.Envs.MongoDatabase).Collection("users")
	filter := bson.M{"credential.email": email}
	var user types.User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByID(id string) (*types.User, error) {
	collection := s.mongoClient.Database(config.Envs.MongoDatabase).Collection("users")

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format")
	}

	filter := bson.M{"_id": objID}

	var user types.User
	er := collection.FindOne(context.TODO(), filter).Decode(&user)
	if er == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("user not found")
	} else if er != nil {
		return nil, er
	}
	return &user, nil
}

func (s *Store) UpdateUserSessionRooms(userID primitive.ObjectID, roomID string) error {
	collectionUser := s.mongoClient.Database(config.Envs.MongoDatabase).Collection("users")

	sessionRoomID := types.SessionRoomID{
		RoomID:    roomID,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$push": bson.M{
			"sessionRoomIDs": bson.M{"$each": []interface{}{sessionRoomID}},
		},
		"$set": bson.M{
			"lastUpdatedAt": time.Now(),
		},
	}

	_, err := collectionUser.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update chat history: %w", err)
	}
	return nil
}
