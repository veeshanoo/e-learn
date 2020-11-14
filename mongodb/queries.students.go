package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (mc *MongoClient) GetStudent(filter bson.M) (*Student, error) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Students)

	student := &Student{}
	student.Id = primitive.NewObjectID().Hex()
	if err := collection.FindOne(ctx, filter).Decode(student); err != nil {
		return nil, err
	}

	return student, nil
}

func (mc *MongoClient) InsertStudent(student *Student) error {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Students)

	user, err := mc.GetUser(student.Email, "", false)
	if err != nil {
		return err
	}

	student.Id = user.Id
	if _, err := collection.InsertOne(ctx, student); err != nil {
		return err
	}

	return nil
}
