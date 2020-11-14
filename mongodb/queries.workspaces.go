package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (mc *MongoClient) GetWorkspace(name string) (*Workspace, error) {
	filter := bson.M{"name": name}

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Workspaces)

	workspace := &Workspace{}
	if err := collection.FindOne(ctx, filter).Decode(workspace); err != nil {
		return nil, err
	}

	return workspace, nil
}

func (mc *MongoClient) GetWorkspaces() ([]*Workspace, error) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Workspaces)

	filter := bson.M{}

	iter, err := collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	var workspaces []*Workspace
	for iter.Next(ctx) {
		workspace := &Workspace{}
		if err := iter.Decode(workspace); err != nil {
			return nil, err
		}

		workspaces = append(workspaces, workspace)
	}

	return workspaces, nil
}

func (mc *MongoClient) InsertWorkspace(workspace *Workspace) error {
	if _, err := mc.GetWorkspace(workspace.Name); err == nil {
		return errors.New("workspace already exists")
	}

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	collection := mc.Client.Database(MyDb.DbName).Collection(MyDb.Workspaces)

	workspace.Id = primitive.NewObjectID().Hex()

	if _, err := collection.InsertOne(ctx, workspace); err != nil {
		return err
	}

	return nil
}
