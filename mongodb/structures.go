package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MongoDb struct {
	Url      string `json:"url"`
	DbName   string `json:"name"`
	Users    string `json:"users"`
	Sessions string `json:"sessions"`
}

type UserType int

const (
	UserType_Student UserType = 0
	UserType_Doctor  UserType = 1
)

type UserAuth struct {
	Id       primitive.ObjectID   `json:"_id" bson:"_id"`
	Email    string   `json:"email" bson:"email"`
	Password string   `json:"password" bson:"password"`
	Type     UserType `json:"user_type" bson:"user_type"`
}

type Session struct {
	Id        primitive.ObjectID    `json:"_id" bson:"_id"`
	Token     string    `json:"token" bson:"token"`
	Email     string    `json:"email" bson:"email"`
	Type      UserType  `json:"user_type" bson:"user_type"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
