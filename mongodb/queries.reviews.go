package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (mc *MongoClient) AddReview(review *Review, courseId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Courses)

	course, err := mc.GetCourse(courseId)
	if err != nil {
		return err
	}

	student, err := mc.GetStudent(bson.M{"_id": review.StudentId})
	if err != nil {
		return err
	}
	review.Name = student.Name
	course.Reviews = append(course.Reviews, review)

	filter := bson.M{"_id": courseId}
	update := bson.M{"$set": bson.M{"reviews": course.Reviews}}

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}
