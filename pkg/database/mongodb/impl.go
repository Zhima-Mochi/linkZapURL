package mongodb

import (
	"context"

	"github.com/Zhima-Mochi/linkZapURL/config"
	"github.com/Zhima-Mochi/linkZapURL/pkg/database"
	"go.mongodb.org/mongo-driver/mongo"
)

type Impl struct {
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

	return &Impl{
		client: client,
		config: config,
	}, nil
}

func (im *Impl) getCollection(collectionName string) *mongo.Collection {
	return im.client.Database(im.config.Database).Collection(collectionName)
}

func (im *Impl) Get(ctx context.Context, table, key string) (interface{}, error) {
	collection := im.getCollection(table)

	data, err := collection.Find(ctx, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (im *Impl) Set(ctx context.Context, table, key string, value interface{}) error {
	collection := im.getCollection(table)

	_, err := collection.InsertOne(ctx, value)
	if err != nil {
		return err
	}

	return nil
}
