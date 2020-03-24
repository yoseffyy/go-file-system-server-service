package fileModel

import (
	"fmt"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	dbName string = "test"
	collectionName string = "files"
)

type File struct {
	ID       primitive.ObjectID `bson:"_id, omitempty"`
	Owner    string             `bson:"owner"`
	Name     string             `bson:"name"`
	Path     string             `bson:"path"`
	IsFolder bool               `bson:"isFolder"`
}

func New(owner, name, path string, isFolder bool) File {
	return File {
		ID:    primitive.NewObjectID(),
		Owner: owner,
		Name:  name,
		Path:  path,
		IsFolder: isFolder,
	};
}

func (f *File) Insert(client *mongo.Client) (*File, error) {
    collection := client.Database(dbName).Collection(collectionName)
    res, err := collection.InsertOne(context.TODO(), f)
	if err != nil {
        return nil, fmt.Errorf(
			"Internal",
			fmt.Sprintf("Internal error: %v", err),
		)
    }
	f.ID = res.InsertedID.(primitive.ObjectID);
	return f, nil;
}

func (f *File) Update(client *mongo.Client, filter bson.M) (*File, error) {
	collection := client.Database(dbName).Collection(collectionName)
	
	updatedData := bson.M{};
	fileBytes, _ := bson.Marshal(f)
	bson.Unmarshal(fileBytes, &updatedData)

    atualizacao := bson.D{ {Key: "$set",Value: updatedData} }
	result, err := collection.UpdateOne(context.TODO(), filter, atualizacao)
	if err != nil {
        return nil, fmt.Errorf(
			"Internal",
			fmt.Sprintf("Internal error: %v", err),
		)
    }
	
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf(
			"Internal",
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	return f, nil;
}

func FinedOne(client *mongo.Client, filter bson.M) (*File, error) {
	collection := client.Database(dbName).Collection(collectionName)

	result := collection.FindOne(context.TODO(), filter)

	decoded := File{}
	err := result.Decode(&decoded)

	if err != nil {
		return nil, fmt.Errorf(
			"NotFound",
			fmt.Sprintf("Could not find blog with supplied ID: %v", err),
		)
	}
	
    return &decoded, nil
}

func FinedAll(client *mongo.Client, filter bson.M) ([]*File, error) {
    files := []*File{};
	collection := client.Database(dbName).Collection(collectionName)

	cursor, err := collection.Find(context.TODO(), filter)
    if err != nil {
        return nil, fmt.Errorf("Internal", fmt.Sprintf("Unknown internal error: %v", err))
	}
	
    for cursor.Next(context.TODO()) {
    	file := File{};
		err := cursor.Decode(&file)
		if err != nil {
			return nil, fmt.Errorf("Unavailable", fmt.Sprintf("Could not decode data: %v", err))
		}
        files = append(files, &file)
	}
    return files, nil
}
