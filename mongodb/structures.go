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
type WorkspaceType int

const (
	UserType_Student UserType = 0
	UserType_Doctor  UserType = 1
)

const (
	WorkspaceType_Student   WorkspaceType = 0
	WorkspaceType_Universal WorkspaceType = 1
)

type UserAuth struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Type     UserType           `json:"user_type" bson:"user_type"`
}

type Session struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Token     string             `json:"token" bson:"token"`
	Email     string             `json:"email" bson:"email"`
	Type      UserType           `json:"user_type" bson:"user_type"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Workspace struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	Name          string             `json:"name" bson:"name"`
	WorkspaceType WorkspaceType      `json:"type" bson:"type"`
}

type Category struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	WorkspaceId primitive.ObjectID `json:"workspace_id" bson:"workspace_id"`
	Name        string             `json:"name" bson:"name"`
}

type Course struct {
	Id         primitive.ObjectID   `json:"_id" bson:"_id"`
	CategoryId primitive.ObjectID   `json:"cat_id" bson:"cat_id"`
	TeacherId  primitive.ObjectID   `json:"teacher_id" bson:"teacher_id"`
	Duration   int32                `json:"duration" bson:"duration"`
	Difficulty int32                `json:"difficulty" bson:"difficulty"`
	Lessons    []*Lesson            `json:"lessons" bson:"lessons"`
	Students   []primitive.ObjectID `json:"students" bson:"students"`
	Reviews    []*Review            `json:"reviews" bson:"reviews"`
}

type Lesson struct {
	Id        primitive.ObjectID   `json:"_id" bson:"_id"`
	Students  []primitive.ObjectID `json:"students" bson:"students"`
	Video     Video                `json:"video" bson:"video"`
	Materials []string             `json:"materials" bson:"materials"`
}

type Video struct {
	URL       string `json:"url" bson:"url"`
	Subtitles string `json:"subtitles" bson:"subtitles"`
}

type Student struct {
	Id          primitive.ObjectID   `json:"_id" bson:"_id"`
	Name        string               `json:"name" bson:"name"`
	Courses     []primitive.ObjectID `json:"courses" bson:"courses"`
	Assignments []primitive.ObjectID `json:"assignments" bson:"assignments"`
}

type Teacher struct {
	Id      primitive.ObjectID   `json:"_id" bson:"_id"`
	Name    string               `json:"name" bson:"name"`
	Courses []primitive.ObjectID `json:"courses" bson:"courses"`
}

type Assignment struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	LessonId  primitive.ObjectID `json:"lesson_id" bson:"lesson_id"`
	StudentId primitive.ObjectID `json:"student_id" bson:"student_id"`
	Grade     int32              `json:"grade" bson:"grade"`
}

type Review struct {
	StudentId primitive.ObjectID `json:"student_id" bson:"student_id"`
	Rating    int32              `json:"rating" bson:"rating"`
	Feedback  string             `json:"feedback"`
}
