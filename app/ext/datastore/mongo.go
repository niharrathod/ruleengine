package datastore

import (
	"context"
	"os"
	"time"

	"github.com/niharrathod/ruleengine/app/config"
	"github.com/niharrathod/ruleengine/app/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.uber.org/zap"
)

const (
	database           = "ruleengineDB"
	ruleEngineCollName = "ruleengine"
	configCollName     = "ruleengineconfig"
)

var client *mongo.Client
var ruleEngineCollection *mongo.Collection
var engineConfigCollection *mongo.Collection

func Initialize() {
	mongoUrl := config.Datastore.Mongo.Url
	username := config.Datastore.Mongo.Username
	password := config.Datastore.Mongo.Password

	var clientOptions []*options.ClientOptions
	if config.IsProduction() {
		clientOptions = []*options.ClientOptions{
			options.Client().SetAuth(options.Credential{Username: username, Password: password}),
			options.Client().ApplyURI(mongoUrl),
			options.Client().SetTimeout(time.Second),
			options.Client().SetReadConcern(readconcern.Majority()),
			options.Client().SetWriteConcern(writeconcern.New(writeconcern.J(true), writeconcern.W(3)))}
	} else {
		clientOptions = []*options.ClientOptions{
			options.Client().SetAuth(options.Credential{Username: username, Password: password}),
			options.Client().ApplyURI(mongoUrl),
			options.Client().SetReadConcern(readconcern.Majority()),
			options.Client().SetWriteConcern(writeconcern.New(writeconcern.J(true), writeconcern.W(1)))}
	}

	mongoClient, err := mongo.Connect(context.TODO(), clientOptions...)
	if err != nil {
		log.Logger.Error("Client connect failed", zap.String("error", err.Error()))
		os.Exit(1)
	}

	client = mongoClient

	pingCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		log.Logger.Error("Client ping failed", zap.String("error", err.Error()))
		os.Exit(1)
	}

	ruleEngineCollection = client.Database(database).Collection(ruleEngineCollName)
	engineConfigCollection = client.Database(database).Collection(configCollName)

	// RuleEngine name index
	model := mongo.IndexModel{Keys: bson.D{{Key: "name", Value: 1}}}
	name, err := ruleEngineCollection.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		log.Logger.Error("MongoDB index creation failed", zap.String("error", err.Error()))
		os.Exit(1)
	} else {
		log.Logger.Info("MongoDB index creation succeed", zap.String("IndexName", name))
	}
}

func Close(ctx context.Context) error {
	log.Logger.Info("Mongo client closing")
	return client.Disconnect(ctx)
}
