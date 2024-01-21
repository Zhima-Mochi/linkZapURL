package mongodb

import (
	"context"

	"github.com/Zhima-Mochi/linkZapURL/config"
	"github.com/Zhima-Mochi/linkZapURL/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type impl struct {
	client *mongo.Client
	config *config.Mongodb
}

func NewMongodb(config *config.Mongodb) (database.Database, error) {
	client, err := mongo.Connect(context.Background(), config.GetClientOptions()...)
	if err != nil {
		return nil, err
	}

	// Ping the primary
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &impl{
		client: client,
		config: config,
	}, nil
}

func (im *impl) getCollection(collectionName string) *mongo.Collection {
	return im.client.Database(im.config.Database).Collection(collectionName)
}

func (im *impl) Get(ctx context.Context, table string, key int64, result interface{}) error {
	collection := im.getCollection(table)

	filter := bson.M{"_id": key}

	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return err
	}

	return nil
}

func (im *impl) Set(ctx context.Context, table string, key int64, value interface{}) error {
	collection := im.getCollection(table)

	_, err := collection.InsertOne(ctx, value)
	if err != nil {
		return err
	}

	return nil
}
