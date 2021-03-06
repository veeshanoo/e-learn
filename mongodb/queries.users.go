package mongodb

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (mc *MongoClient) GetUser(email string, password string, checkPwd bool) (*UserAuth, error) {
	filter := bson.M{"email": email}

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Users)

	user := &UserAuth{}
	if err := collection.FindOne(ctx, filter).Decode(user); err != nil {
		return nil, err
	}

	if !checkPwd {
		user.Password = ""
		return user, nil
	}

	if err := comparePasswords(user.Password, password); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (mc *MongoClient) InsertNewUser(user *UserAuth) (*Session, error) {
	if _, err := mc.GetUser(user.Email, "", false); err == nil {
		return nil, errors.New("email already exists")
	}

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Users)

	pwd := user.Password
	if hashedPwd, err := hashAndSalt(pwd); err != nil {
		return nil, err
	} else {
		user.Password = hashedPwd
	}

	user.Id = primitive.NewObjectID().Hex()
	if _, err := collection.InsertOne(ctx, user); err != nil {
		return nil, err
	}

	return mc.InsertSession(user.Email)
}

func (mc *MongoClient) InsertSession(email string) (*Session, error) {
	user, err := mc.GetUser(email, "", false)
	if err != nil {
		return nil, err
	}

	token := uuid.New().String()

	session := Session{
		Id:        primitive.NewObjectID().Hex(),
		Token:     token,
		Email:     user.Email,
		CreatedAt: time.Now(),
		Type:      user.Type,
	}

	//if err := mc.ClearSession(email); err != nil {
	//	return nil, err
	//}

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Sessions)
	_, err = collection.InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (mc *MongoClient) GetSession(token string) (*Session, error) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Sessions)

	session := &Session{}
	err := collection.FindOne(ctx, bson.M{"token": token}).Decode(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (mc *MongoClient) ClearSession(email string) error {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Sessions)

	filter := bson.M{"email": email}
	if _, err := collection.DeleteMany(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (mc *MongoClient) AddCourse(email string, userType UserType, courseId string) error {
	if _, err := mc.GetCourse(courseId); err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := &mongo.Collection{}

	var courses []string
	if userType == UserType_Student {
		collection = mc.Client.Database(MyDb.DbName).Collection(MyDb.Students)

		student, err := mc.GetStudent(bson.M{"email": email})
		if err != nil {
			return err
		}

		courses = student.Courses
	} else {
		collection = mc.Client.Database(MyDb.DbName).Collection(MyDb.Teachers)

		teacher, err := mc.GetTeacher(bson.M{"email": email})
		if err != nil {
			return err
		}

		courses = teacher.Courses
	}

	addStudent(&courses, courseId)

	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"courses": courses}}

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}
