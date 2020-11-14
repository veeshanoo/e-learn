package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (mc *MongoClient) NewAssignment(studentId, lessonId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Assignments)

	assignment := Assignment{
		Id:        primitive.NewObjectID().Hex(),
		StudentId: studentId,
		LessonId:  lessonId,
		Grade:     -1,
	}

	if _, err := collection.InsertOne(ctx, assignment); err != nil {
		return err
	}

	return nil
}

func (mc *MongoClient) UpdateAssignment() {

}
