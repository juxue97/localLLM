package chatbot

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

func (s *Store) StoreChatHistory(input string, output *types.Output, roomID string) error {
	collection := s.mongoClient.Database(config.Envs.MongoDatabase).Collection("chatrooms")
	objID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return fmt.Errorf("invalid roomID format: %v", err)
	}

	// Create the user and assistant message structures
	conversations := types.ConversationEntry{
		UserMessage: struct {
			Content    string `bson:"content" json:"content"`
			InputToken int    `bson:"inputToken" json:"inputToken"`
		}{
			Content:    input,             // User's message content
			InputToken: output.InputToken, // Input token count for the user message
		},
		AssistantMessage: struct {
			Content     string `bson:"content" json:"content"`
			OutputToken int    `bson:"outputToken" json:"outputToken"`
		}{
			Content:     output.Response,    // Assistant's response content
			OutputToken: output.OutputToken, // Output token count for the assistant message
		},
	}

	// Update the chatroom document with the new conversation history
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$push": bson.M{
			"conversations": bson.M{"$each": []interface{}{conversations}},
		},
		"$set": bson.M{
			"lastUpdatedAt": time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update chat history: %v", err)
	}

	return nil
}

func (s *Store) LoadChatHistory(userID string, roomID string) ([]map[string]string, error) {
	collectionChatrooms := s.mongoClient.Database(config.Envs.MongoDatabase).Collection("chatrooms")
	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, fmt.Errorf("invalid roomID format: %v", err)
	}

	filter := bson.M{"_id": objID, "userID": userID}

	var chatroom types.ChatbotHistory
	er := collectionChatrooms.FindOne(context.TODO(), filter).Decode(&chatroom)
	if er == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("chat history not found")
	} else if er != nil {
		return nil, fmt.Errorf("error loading chat history: %v", er)
	}

	var messages []map[string]string
	for _, conv := range chatroom.Conversations {
		// User message map
		userMessage := map[string]string{
			"role":    "user",
			"content": conv.UserMessage.Content,
		}
		messages = append(messages, userMessage)

		// Assistant message map
		assistantMessage := map[string]string{
			"role":    "assistant",
			"content": conv.AssistantMessage.Content,
		}
		messages = append(messages, assistantMessage)
	}

	return messages, nil
}

func (s *Store) CreateSessionRoom(userID string) (string, error) {
	collectionChatrooms := s.mongoClient.Database(config.Envs.MongoDatabase).Collection("chatrooms")

	room := bson.M{
		"userID":        userID,
		"createdAt":     time.Now(),
		"lastUpdatedAt": time.Now(),
	}
	result, err := collectionChatrooms.InsertOne(context.Background(), room)
	if err != nil {
		return "", fmt.Errorf("error creating session room: %v", err)
	}

	// Return the newly created roomID
	roomID := result.InsertedID.(primitive.ObjectID).Hex()

	return roomID, nil
}

func (s *Store) DeleteSessionRoom(userID primitive.ObjectID, roomID string) error {
	collection := s.mongoClient.Database(config.Envs.MongoDatabase).Collection("users")

	filter := bson.M{"_id": userID}

	deletion := bson.M{
		"$pull": bson.M{
			"sessionRoomIDs": bson.M{"roomID": roomID},
		},
	}

	result, err := collection.UpdateOne(context.TODO(), filter, deletion)
	if err != nil {
		return fmt.Errorf("error deleting session room: %v", err)
	}

	// Check if any document was matched and modified
	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found for the given userID")
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("roomID not found in sessionRoomIDs for this user")
	}
	return nil
}

func (s *Store) GetAllSessionRoomID(userID primitive.ObjectID) ([]string, error) {
	collection := s.mongoClient.Database(config.Envs.MongoDatabase).Collection("users")

	filter := bson.M{"_id": userID}

	var user types.User

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found with given userID")
		}
		return nil, fmt.Errorf("error retrieving session rooms: %v", err)
	}

	roomIDs := make([]string, len(user.SessionRoomIDs))
	for i, session := range user.SessionRoomIDs {
		roomIDs[i] = session.RoomID
	}
	return roomIDs, nil
}

func (s *Store) GetSessionRoomID(userID string) (*types.User, error) {
	collection := s.mongoClient.Database(config.Envs.MongoDatabase).Collection("users")

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid roomID format")
	}

	filter := bson.M{"_id": objID}

	var user types.User
	er := collection.FindOne(context.TODO(), filter).Decode(&user)
	if er == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("user not found")
	} else if er != nil {
		return nil, fmt.Errorf("error retrieving roomID: %v", er)
	}
	return &user, nil
}
