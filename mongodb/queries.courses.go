package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (mc *MongoClient) GetCourse(id string) (*Course, error) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Courses)

	filter := bson.M{"_id": id}

	course := &Course{}
	if err := collection.FindOne(ctx, filter).Decode(course); err != nil {
		return nil, err
	}

	return course, nil
}

func (mc *MongoClient) GetCourses(workspaceId, catId string) ([]*Course, error) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Courses)

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
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Courses)

	course.Id = primitive.NewObjectID().Hex()

	if _, err := collection.InsertOne(ctx, course); err != nil {
		return err
	}

	return nil
}

func (mc *MongoClient) JoinCourse(courseId string, token string) error {
	session, err := mc.GetSession(token)
	if err != nil {
		return err
	}

	user, err := mc.GetUser(session.Email, "", false)
	if err != nil {
		return err
	}

	course, err := mc.GetCourse(courseId)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Courses)
	filter := bson.M{"_id": courseId}

	if user.Type == UserType_Student {
		addStudent(&course.Students, user.Id)
	} else {
		addStudent(&course.Teachers, user.Id)
	}

	if _, err := collection.UpdateOne(ctx, filter, bson.M{"$set": course}); err != nil {
		return err
	}

	if err := mc.AddCourse(user.Email, user.Type, courseId); err != nil {
		return err
	}

	return nil
}

func (mc *MongoClient) GetCoursesForUser(email string, userType UserType) ([]*Course, error) {
	var courses []string
	if userType == UserType_Student {
		student, err := mc.GetStudent(bson.M{"email": email})
		if err != nil {
			return nil, err
		}
		courses = student.Courses
	} else {
		teacher, err := mc.GetTeacher(bson.M{"email": email})
		if err != nil {
			return nil, err
		}
		courses = teacher.Courses
	}

	var crs []*Course
	for _, el := range courses {
		course, err := mc.GetCourse(el)
		if err != nil {
			return nil, err
		}

		crs = append(crs, course)
	}

	return crs, nil
}
