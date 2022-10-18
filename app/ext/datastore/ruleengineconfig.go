package datastore

import (
	"context"

	ruleenginecore "github.com/niharrathod/ruleengine-core"
	"github.com/niharrathod/ruleengine/app/entities"
	"github.com/niharrathod/ruleengine/app/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func DeleteRuleEngineConfig(ctx context.Context, ruleEngineName string, tag string) *entities.Error {

	deleteRuleEngineConfigTxnFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		existingEngine, err := getRuleEngine(sessCtx, ruleEngineName)
		if err != nil {
			log.Logger.Error("Get RuleEngine failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		if existingEngine == nil {
			return nil, entities.NewError(entities.ErrCodeRuleEngineNotFound)
		}

		t, ok := existingEngine.Tags[tag]
		if !ok {
			return nil, entities.NewError(entities.ErrCodeTagNotFound)
		}
		if t.IsEnable {
			return nil, entities.NewError(entities.ErrCodeTagDeleteNotAllowed)
		}

		if _, err := engineConfigCollection.DeleteOne(sessCtx, bson.M{"_id": t.EngineConfigID}); err != nil {
			log.Logger.Error("Delete RuleEngineConfig failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		delete(existingEngine.Tags, tag)

		if err := upsertRuleEngine(sessCtx, existingEngine); err != nil {
			log.Logger.Error("Upsert RuleEngine failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		return nil, nil
	}

	session, err := client.StartSession()
	if err != nil {
		log.Logger.Error("DeleteRuleEngineConfig StartSession() failed", zap.String("Error", err.Error()))
		return entities.NewError(entities.ErrCodeDatastoreFailed)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, deleteRuleEngineConfigTxnFunc)
	if err != nil {
		if txnErr, ok := err.(*entities.Error); ok {
			log.Logger.Error("DeleteRuleEngineConfig WithTransaction() failed", zap.String("Error", txnErr.Error()))
			return txnErr
		} else {
			log.Logger.Error("DeleteRuleEngineConfig WithTransaction() failed, assert txnError failed", zap.String("Error", err.Error()))
			return entities.NewError(entities.ErrCodeDatastoreFailed)
		}
	}

	return nil
}

func getRuleEngineConfig(ctx context.Context, id primitive.ObjectID) (*ruleenginecore.RuleEngineConfig, error) {
	var config entities.EngineConfig
	err := engineConfigCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&config)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return config.EngineCoreConfig, nil
}
