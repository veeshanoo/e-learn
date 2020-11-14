package mongodb

import (
	"time"
)

type MongoDb struct {
	Url        string `json:"url"`
	DbName     string `json:"name"`
	Users      string `json:"users"`
	Sessions   string `json:"sessions"`
	Workspaces string `json:"workspaces"`
	Categories string `json:"categories"`
	Courses    string `json:"courses"`
	Students   string `json:"students"`
	Teachers   string `json:"teachers"`
	Lessons    string `json:"lessons"`
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
	Id       string   `json:"_id" bson:"_id"`
	Email    string   `json:"email" bson:"email"`
	Password string   `json:"password" bson:"password"`
	Type     UserType `json:"user_type" bson:"user_type"`
}

type Session struct {
	Id        string    `json:"_id" bson:"_id"`
	Token     string    `json:"token" bson:"token"`
	Email     string    `json:"email" bson:"email"`
	Type      UserType  `json:"user_type" bson:"user_type"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Workspace struct {
	Id            string        `json:"_id" bson:"_id"`
	ImageUrl      string        `json:"url" bson:"url"`
	Name          string        `json:"name" bson:"name"`
	WorkspaceType WorkspaceType `json:"type" bson:"type"`
}

type Category struct {
	Id          string `json:"_id" bson:"_id"`
	ImageUrl    string `json:"url" bson:"url"`
	WorkspaceId string `json:"workspace_id" bson:"workspace_id"`
	Name        string `json:"name" bson:"name"`
}

type Course struct {
	Id          string    `json:"_id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	ImageUrl    string    `json:"url" bson:"url"`
	WorkspaceId string    `json:"workspace_id" bson:"workspace_id"`
	CategoryId  string    `json:"cat_id" bson:"cat_id"`
	TeacherId   string    `json:"teacher_id" bson:"teacher_id"`
	Duration    int32     `json:"duration" bson:"duration"`
	Difficulty  int32     `json:"difficulty" bson:"difficulty"`
	Lessons     []*Lesson `json:"lessons" bson:"lessons"`
	Students    []string  `json:"students" bson:"students"`
	Teachers    []string  `json:"teachers" bson:"teachers"`
	Reviews     []*Review `json:"reviews" bson:"reviews"`
}

type Lesson struct {
	Id        string   `json:"_id" bson:"_id"`
	Students  []string `json:"students" bson:"students"`
	Video     *Video   `json:"video" bson:"video"`
	Materials []string `json:"materials" bson:"materials"`
}

type Video struct {
	URL       string `json:"url" bson:"url"`
	Subtitles string `json:"subtitles" bson:"subtitles"`
}

type Student struct {
	Id          string   `json:"_id" bson:"_id"`
	Name        string   `json:"name" bson:"name"`
	Email       string   `json:"email"`
	Courses     []string `json:"courses" bson:"courses"`
	Assignments []string `json:"assignments" bson:"assignments"`
}

type Teacher struct {
	Id      string   `json:"_id" bson:"_id"`
	Name    string   `json:"name" bson:"name"`
	Email   string   `json:"email"`
	Courses []string `json:"courses" bson:"courses"`
}

type Assignment struct {
	Id        string `json:"_id" bson:"_id"`
	LessonId  string `json:"lesson_id" bson:"lesson_id"`
	StudentId string `json:"student_id" bson:"student_id"`
	Grade     int32  `json:"grade" bson:"grade"`
}

type Review struct {
	StudentId string `json:"student_id" bson:"student_id"`
	Name      string `json:"name" bson:"name"`
	Rating    int32  `json:"rating" bson:"rating"`
	Feedback  string `json:"feedback"`
}
