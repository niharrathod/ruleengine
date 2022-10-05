package datastore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/niharrathod/ruleengine/app/config"
	"github.com/niharrathod/ruleengine/app/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.uber.org/zap"
)

const (
	ID         = "_id"
	name       = "name"
	tag        = "tag"
	database   = "ruleengineDB"
	collection = "ruleengine"
)

type RuleEngine struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
	Tag  string             `bson:"tag"`
}

var client *mongo.Client
var ruleEngineCollection *mongo.Collection

func Initialize() {
	mongoUrl := config.Datastore.Mongo.Url
	username := config.Datastore.Mongo.Username
	password := config.Datastore.Mongo.Password

	options := []*options.ClientOptions{
		options.Client().SetAuth(options.Credential{Username: username, Password: password}),
		options.Client().ApplyURI(mongoUrl),
		options.Client().SetReadConcern(readconcern.Majority()),
		options.Client().SetWriteConcern(writeconcern.New(writeconcern.J(true), writeconcern.W(1)))}

	mongoClient, err := mongo.Connect(context.TODO(), options...)
	if err != nil {
		log.Logger.Error("Client connect failed. error: %v", zap.String("error", err.Error()))
	}

	client = mongoClient

	pingCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		log.Logger.Error("Client ping failed. error: %v", zap.String("error", err.Error()))
	}

	ruleEngineCollection = client.Database(database).Collection(collection)

	// index creation for text search
	model := mongo.IndexModel{Keys: bson.D{{Key: name, Value: "text"}}}
	name, err := ruleEngineCollection.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		log.Logger.Error("Index create failed error: %v", zap.String("error", err.Error()))
	}
	fmt.Println("Index Created by name: " + name)

	// rule engine name and tag index
	model = mongo.IndexModel{Keys: bson.D{{Key: name, Value: 1}, {Key: tag, Value: 2}}}
	name, err = ruleEngineCollection.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		log.Logger.Error("Index create failed error: %v", zap.String("error", err.Error()))
	}
	fmt.Println("Index Created by name: " + name)
}

func Close(ctx context.Context) error {
	log.Logger.Info("Mongo client closing")
	return client.Disconnect(ctx)
}

func CreateRuleEngine(ctx context.Context, ruleEngine RuleEngine) error {
	_, err := ruleEngineCollection.InsertOne(ctx, ruleEngine)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRuleEngine(ctx context.Context, ruleEngineName, ruleEngineTag string) error {
	filter := bson.D{{Key: name, Value: ruleEngineName}, {Key: tag, Value: ruleEngineTag}}
	_, err := ruleEngineCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func GetRuleEngine(ctx context.Context, ruleEngineName, ruleEngineTag string) ([]*RuleEngine, error) {
	var filter primitive.D
	if len(ruleEngineTag) != 0 {
		filter = bson.D{{Key: name, Value: ruleEngineName}, {Key: tag, Value: ruleEngineTag}}
	} else {
		filter = bson.D{{Key: name, Value: ruleEngineName}}
	}

	cursor, err := ruleEngineCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := []*RuleEngine{}
	if err = cursor.All(ctx, &result); err != nil {
		return nil, errors.New("Failed to unmarshal result for, ruleEngineName:" + ruleEngineName + " ruleEngineTag:" + ruleEngineTag)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

func SearchByName(ctx context.Context, nameQuery string) ([]*RuleEngine, error) {
	filter := bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: nameQuery}}}}
	options := options.Find().SetSort(bson.D{{Key: "name", Value: -1}})

	cursor, err := ruleEngineCollection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}

	result := []*RuleEngine{}
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}
