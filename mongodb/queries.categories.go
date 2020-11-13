package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (mc *MongoClient) GetCategories(workspaceId primitive.ObjectID) ([]*Category, error) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Categories)

	filter := bson.M{"workspace_id": workspaceId}

	iter, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var categories []*Category
	for iter.Next(ctx) {
		cat := &Category{}
		if err := iter.Decode(cat); err != nil {
			return nil, err
		}

		categories = append(categories, cat)
	}

	return categories, nil
}

func (mc *MongoClient) InsertCategory(cat *Category) error {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Categories)

	cat.Id = primitive.NewObjectID()

	if _, err := collection.InsertOne(ctx, cat); err != nil {
		return err
	}

	return nil
}
