package db

import (
	"api/pkg/models"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrDuplicateExternalId = errors.New("externalId already exists")

type Repository interface {
	Ping() error
	AddOne(collectionName string, data interface{}) error
	GetAllTemplates(filter primitive.D, options *options.FindOptions) ([]models.Template, error)
	GetAllInstances(filter primitive.D, options *options.FindOptions) ([]models.Instance, error)
	GetTemplate(filter primitive.D) (*models.Template, error)
	GetTypeDropdownValues(collection string) ([]models.Dropdown, error)
	ReplaceTemplate(filter primitive.D, data interface{}) error
}

type repository struct {
	client *mongo.Client
}

func NewRepository(client *mongo.Client) Repository {
	return &repository{
		client: client,
	}
}

func (r *repository) ReplaceTemplate(filter primitive.D, data interface{}) error {
	collection := r.client.Database("buildifyy").Collection("templates")
	_, err := collection.ReplaceOne(context.Background(), filter, data)
	if err != nil {
		log.Println("error replacing data in database")
		return err
	}

	return nil
}

func (r *repository) GetTypeDropdownValues(collection string) ([]models.Dropdown, error) {
	c := r.client.Database("buildifyy").Collection(collection)
	opts := options.Find().SetSort(bson.D{{Key: "label", Value: 1}})
	cursor, err := c.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		log.Println("error finding dropdown values in database: ", err)
		return nil, err
	}

	var results []models.Dropdown
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Println("error parsing all data from database: ", err)
		return nil, err
	}

	return results, nil
}

func (r *repository) Ping() error {
	if err := r.client.Database("admin").RunCommand(context.Background(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Println("error pinging database: ", err)
		return err
	}

	return nil
}

func (r *repository) GetAllTemplates(filter primitive.D, options *options.FindOptions) ([]models.Template, error) {
	collection := r.client.Database("buildifyy").Collection("templates")
	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		log.Println("error finding data in database: ", err)
		return nil, err
	}

	var results []models.Template
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Println("error parsing all data from database: ", err)
		return nil, err
	}

	return results, nil
}

func (r *repository) GetAllInstances(filter primitive.D, options *options.FindOptions) ([]models.Instance, error) {
	collection := r.client.Database("buildifyy").Collection("instances")
	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		log.Println("error finding data in database: ", err)
		return nil, err
	}

	var results []models.Instance
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Println("error parsing all data from database: ", err)
		return nil, err
	}

	return results, nil
}

func (r *repository) GetTemplate(filter primitive.D) (*models.Template, error) {
	collection := r.client.Database("buildifyy").Collection("templates")

	var result models.Template
	if err := collection.FindOne(context.Background(), filter).Decode(&result); err != nil {
		log.Println("error finding data in database: ", err)
		return nil, err
	}

	return &result, nil
}

func (r *repository) AddOne(collectionName string, data interface{}) error {
	collection := r.client.Database("buildifyy").Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Println("error inserting data to database")
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateExternalId
		}
		return err
	}

	return nil
}
