package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson: "_id,omitempty" json: "id"`
	Name           string             `bson:"name" json: "name"`
	Email          string             `bson:"email" json: "email"`
	TelegramHandle string             `bson:"telegramHandle,omitempty" json: "telegramHandle,omitempty"`
	PasswordHash   string             `bson:"passwordHash" json: "-"` // don't return in API
	CreatedAt      time.Time          `bson:"createdAt" json: "createdAt"`
}
