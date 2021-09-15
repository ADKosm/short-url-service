package main

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	collectionName = "urlsByKey"
	databaseName   = "urlShortener"
)

func newMongoStorage(mongoURL string) *mongoStorage {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		panic(fmt.Errorf("failed to connect to MongoDB, cause: %w", err))
	}
	return &mongoStorage{
		urlsCollection: client.Database(databaseName).Collection(collectionName),
	}
}

type mongoStorage struct {
	urlsCollection *mongo.Collection
}

var _ Storage = (*mongoStorage)(nil)

func (s *mongoStorage) PutURL(ctx context.Context, url RedirectURL) (Key, error) {
	for retriesCount := 0; retriesCount < 5; retriesCount++ {
		key := generateKey()
		result, err := s.urlsCollection.InsertOne(ctx, urlEntry{
			Key: key,
			URL: url,
		})
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				log.Printf("Duplicate key %s on insertion. Retry insertion with new key", key)
				continue // retry
			}
			return "", err
		}
		return Key(result.InsertedID.(string)), nil
	}
	return "", fmt.Errorf("%w: failed to generate unique keys multiple times", ErrInsertionCollision)
}

func (s *mongoStorage) GetURL(ctx context.Context, key Key) (RedirectURL, error) {
	var entry urlEntry
	err := s.urlsCollection.FindOne(ctx, bson.M{"_id": string(key)}).Decode(&entry)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", fmt.Errorf("%w: no url found for key %s, cause: %s", ErrNotFound, key, err.Error())
		}
		return "", fmt.Errorf("%w, cause: %s", StorageError, err.Error())
	}
	return entry.URL, nil
}

type urlEntry struct {
	Key Key         `bson:"_id"` // there is an implicit unique index over _id field, so we utilize it
	URL RedirectURL `bson:"url"`
}
