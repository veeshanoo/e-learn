package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (mc *MongoClient) GetCourses(workspaceId, catId string) ([]*Course, error) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Categories)

	filter := bson.M{}
	filter["workspace_id"] = workspaceId
	if catId != "" {
		filter["cat_id"] = catId
	}

	iter, err := collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	var courses []*Course
	for iter.Next(ctx) {
		course := &Course{}
		if err := iter.Decode(course); err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (mc *MongoClient) InsertCourse(course *Course) error {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Categories)

	course.Id = primitive.NewObjectID().Hex()

	if _, err := collection.InsertOne(ctx, course); err != nil {
		return err
	}

	return nil
}
