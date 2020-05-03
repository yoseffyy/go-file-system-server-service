package fileModel

import (
	"context"
	"fmt"
	"github.com/yosef32/go-file-system-server-service/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type File struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Owner    string             `bson:"owner,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Path     string             `bson:"path,omitempty"`
	IsFolder bool               `bson:"isFolder,omitempty"`
}

func New(owner string, name string, path string, isFolder bool) *File {
	return &File{
		ID:       primitive.NewObjectID(),
		Owner:    owner,
		Name:     name,
		Path:     path,
		IsFolder: isFolder,
	}
}

func (f *File) Insert(ctx context.Context, client *mongo.Client) (*File, error) {
	collection := client.Database(config.DBName()).Collection(config.CollectionName())
	res, err := collection.InsertOne(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("Internal error: %v", err)
	}
	f.ID = res.InsertedID.(primitive.ObjectID)
	return f, nil
}

func (f *File) Update(ctx context.Context, client *mongo.Client, filter bson.M) (*File, error) {
	collection := client.Database(config.DBName()).Collection(config.CollectionName())

	updatedData := bson.M{}
	fileBytes, err := bson.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("Internal error: %v", err)
	}
	err = bson.Unmarshal(fileBytes, &updatedData)
	if err != nil {
		return nil, fmt.Errorf("Internal error: %v", err)
	}
	update := bson.D{{Key: "$set", Value: updatedData}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("Internal error: %v", err)
	}

	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("Internal error: %v", err)
	}

	return f, nil
}

func Fined(ctx context.Context, client *mongo.Client, filter bson.M) (*File, error) {
	collection := client.Database(config.DBName()).Collection(config.CollectionName())

	result := collection.FindOne(ctx, filter)

	decoded := File{}
	err := result.Decode(&decoded)

	if err != nil {
		return nil, fmt.Errorf("Could not find file with supplied ID: %v", err)
	}

	return &decoded, nil
}

func Delete(ctx context.Context, client *mongo.Client, filter bson.M) error {
	collection := client.Database(config.DBName()).Collection(config.CollectionName())

	_, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return fmt.Errorf("Could not felede file with supplied ID: %v", err)
	}

	return nil
}

func FindAll(ctx context.Context, client *mongo.Client, filter bson.M) ([]*File, error) {
	files := []*File{}
	collection := client.Database(config.DBName()).Collection(config.CollectionName())

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Unknown internal error: %v", err)
	}

	for cursor.Next(ctx) {
		file := &File{}
		err := cursor.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("Could not decode data: %v", err)
		}
		files = append(files, file)
	}
	return files, nil
}
