package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(user User) (primitive.ObjectID, error)
	GetUserByID(id string) (*User, error)
	UpdateUserSessionRooms(userID primitive.ObjectID, roomID string) error
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=16"`
}

// Credential holds the user's login credentials
type Credential struct {
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Email     string `bson:"email" json:"email"`
	Password  string `bson:"password" json:"password"`
}

// SessionRoomID represents a user's session with a room
type SessionRoomID struct {
	RoomID    string             `bson:"roomID" json:"roomID"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
}

// User represents a user in the database
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Credential     Credential         `bson:"credential" json:"credential"`
	CreatedAt      primitive.DateTime `bson:"createdAt" json:"createdAt"`
	LastLoginAt    primitive.DateTime `bson:"lastLoginAt" json:"lastLoginAt"`
	Roles          string             `bson:"roles,omitempty" json:"roles"`
	SessionRoomIDs []SessionRoomID    `bson:"sessionRoomIDs,omitempty" json:"sessionRoomIDs"`
	LastSessionIP  string             `bson:"lastSessionIP,omitempty" json:"lastSessionIP"`
	SessionToken   string             `bson:"sessionToken,omitempty" json:"sessionToken"`
}
