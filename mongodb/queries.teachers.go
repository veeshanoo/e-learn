package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (mc *MongoClient) GetTeacher(filter bson.M) (*Teacher, error) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Teachers)

	teacher := &Teacher{}
	if err := collection.FindOne(ctx, filter).Decode(teacher); err != nil {
		return nil, err
	}

	return teacher, nil
}

func (mc *MongoClient) InsertTeacher(teacher *Teacher) error {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Teachers)

	if _, err := collection.InsertOne(ctx, teacher); err != nil {
		return err
	}

	return nil
}
