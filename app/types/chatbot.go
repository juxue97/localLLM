package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatbotStore interface {
	CreateSessionRoom(userID string) (string, error)
	DeleteSessionRoom(userID primitive.ObjectID, roomID string) error
	GetAllSessionRoomID(userID primitive.ObjectID) ([]string, error)
	StoreChatHistory(input string, output *Output, roomID string) error
	LoadChatHistory(userID string, roomID string) ([]map[string]string, error)
}

type ChatbotPayload struct {
	Query         string `json:"query" validate:"required"`
	SessionRoomID string `json:"sessionRoomID,omitempty"`
}

type ConversationEntry struct {
	UserMessage struct {
		Content    string `bson:"content" json:"content"`
		InputToken int    `bson:"inputToken" json:"inputToken"`
	} `bson:"userMessage" json:"userMessage"`

	AssistantMessage struct {
		Content     string `bson:"content" json:"content"`
		OutputToken int    `bson:"outputToken" json:"outputToken"`
	} `bson:"assistantMessage" json:"assistantMessage"`
}

type ChatbotHistory struct {
	SessionRoomID primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID  `bson:"userID" json:"userID"`
	Conversations []ConversationEntry `bson:"conversations" json:"conversations"`
	CreatedAt     time.Time           `bson:"createdAt" json:"createdAt"`
	LastUpdatedAt time.Time           `bson:"lastUpdatedAt" json:"lastUpdatedAt"`
}
