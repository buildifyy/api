package db

import (
	"api/pkg/models"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Repository interface {
	Ping() error
	AddOne(data interface{}) error
	GetAllTemplates(filter primitive.D) ([][]byte, error)
}

type repository struct {
	client *mongo.Client
}

func NewRepository(client *mongo.Client) Repository {
	return &repository{
		client: client,
	}
}

func (r *repository) Ping() error {
	if err := r.client.Database("admin").RunCommand(context.Background(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Println("error pinging database: ", err)
		return err
	}

	return nil
}

func (r *repository) GetAllTemplates(filter primitive.D) ([][]byte, error) {
	collection := r.client.Database("buildifyy").Collection("templates")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Println("error finding data in database: ", err)
		return nil, err
	}

	var results []models.Template
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Println("error parsing all data from database: ", err)
		return nil, err
	}

	ret := make([][]byte, 0)
	for _, data := range results {
		bytesData, err := json.Marshal(data)
		if err != nil {
			log.Println("error marshalling data: ", err)
			return nil, err
		}
		ret = append(ret, bytesData)
	}

	return ret, nil
}

func (r *repository) AddOne(data interface{}) error {
	collection := r.client.Database("buildifyy").Collection("templates")
	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Println("error inserting data to database")
		return err
	}

	return nil
}
