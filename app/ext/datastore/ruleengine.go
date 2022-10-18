package datastore

import (
	"context"
	"fmt"
	"time"

	ruleenginecore "github.com/niharrathod/ruleengine-core"
	"github.com/niharrathod/ruleengine/app/entities"
	"github.com/niharrathod/ruleengine/app/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func CreateRuleEngine(ctx context.Context, ruleEngineName string, tag string, config *ruleenginecore.RuleEngineConfig) *entities.Error {

	createRuleEngineTxnFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		var ruleEngine *entities.RuleEngine
		var err error

		ruleEngine, err = getRuleEngine(sessCtx, ruleEngineName)
		if err != nil {
			log.Logger.Error("Get RuleEngine failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		if ruleEngine != nil {
			// tag already exist check
			if _, ok := ruleEngine.Tags[tag]; ok {
				return nil, entities.NewError(entities.ErrCodeTagAlreadyExist)
			}
		}

		engineConfig := entities.EngineConfig{
			ID:               primitive.NewObjectID(),
			EngineCoreConfig: config,
		}
		if _, err := engineConfigCollection.InsertOne(sessCtx, engineConfig); err != nil {
			log.Logger.Error("Insert RuleEngineConfig failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		if ruleEngine == nil {
			ruleEngine = &entities.RuleEngine{
				Name:           ruleEngineName,
				DefaultTag:     "",
				Tags:           map[string]*entities.Tag{},
				LastUpdateTime: time.Now().Unix(),
			}
		}

		if ruleEngine.Tags == nil {
			ruleEngine.Tags = map[string]*entities.Tag{}
		}

		ruleEngine.Tags[tag] = &entities.Tag{
			Name:           tag,
			EngineConfigID: engineConfig.ID,
			IsEnable:       false,
		}

		if err := upsertRuleEngine(sessCtx, ruleEngine); err != nil {
			log.Logger.Error("Upsert RuleEngine failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}
		return nil, nil
	}
	session, err := client.StartSession()
	if err != nil {
		log.Logger.Error("CreateRuleEngine StartSession() failed", zap.String("Error", err.Error()))
		return entities.NewError(entities.ErrCodeDatastoreFailed)
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, createRuleEngineTxnFunc)
	if err != nil {
		if txnErr, ok := err.(*entities.Error); ok {
			log.Logger.Error("CreateRuleEngine WithTransaction() failed", zap.String("Error", err.Error()))
			return txnErr
		} else {
			log.Logger.Error("CreateRuleEngine WithTransaction() failed, assert txnError failed", zap.String("Error", err.Error()))
			return entities.NewError(entities.ErrCodeDatastoreFailed)
		}
	}

	return nil
}

func DeleteRuleEngine(ctx context.Context, ruleEngineName string) *entities.Error {

	deleteRuleEngineTxnFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		existingEngine, err := getRuleEngine(sessCtx, ruleEngineName)
		if err != nil {
			log.Logger.Error("Get RuleEngine failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		if existingEngine == nil {
			return nil, entities.NewError(entities.ErrCodeRuleEngineNotFound)
		}

		objectIds := []primitive.ObjectID{}
		for _, tag := range existingEngine.Tags {
			objectIds = append(objectIds, tag.EngineConfigID)
		}

		filter := bson.M{"_id": bson.M{"$in": objectIds}}
		if _, err := engineConfigCollection.DeleteMany(sessCtx, filter); err != nil {
			log.Logger.Error("Delete RuleEngineConfig failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		if _, err := ruleEngineCollection.DeleteOne(sessCtx, bson.D{{Key: "name", Value: ruleEngineName}}); err != nil {
			log.Logger.Error("Delete RuleEngine failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		return nil, nil
	}

	session, err := client.StartSession()
	if err != nil {
		log.Logger.Error("DeleteRuleEngine StartSession() failed", zap.String("Error", err.Error()))
		return entities.NewError(entities.ErrCodeDatastoreFailed)
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, deleteRuleEngineTxnFunc)
	if err != nil {
		if txnErr, ok := err.(*entities.Error); ok {
			log.Logger.Error("DeleteRuleEngine WithTransaction() failed", zap.String("Error", err.Error()))
			return txnErr
		} else {
			log.Logger.Error("DeleteRuleEngine WithTransaction() failed, assert txnError failed", zap.String("Error", err.Error()))
			return entities.NewError(entities.ErrCodeDatastoreFailed)
		}
	}

	return nil
}

func SetDefaultTag(ctx context.Context, ruleEngineName string, tag string) *entities.Error {

	setDefaultTagTxnFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		existingEngine, err := getRuleEngine(sessCtx, ruleEngineName)
		if err != nil {
			log.Logger.Error("Get RuleEngine failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		if existingEngine == nil {
			return nil, entities.NewError(entities.ErrCodeRuleEngineNotFound)
		}

		t, ok := existingEngine.Tags[tag]
		if !ok || !t.IsEnable {
			return nil, entities.NewError(entities.ErrCodeDefaultTagExistAndMustBeEnabled)
		}

		existingEngine.LastUpdateTime = time.Now().Unix()
		existingEngine.DefaultTag = tag

		if err := upsertRuleEngine(sessCtx, existingEngine); err != nil {
			log.Logger.Error("Upsert RuleEngine failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		return nil, nil
	}

	session, err := client.StartSession()
	if err != nil {
		log.Logger.Error("SetDefaultTag StartSession() failed", zap.String("Error", err.Error()))
		return entities.NewError(entities.ErrCodeDatastoreFailed)
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, setDefaultTagTxnFunc)
	if err != nil {
		if txnErr, ok := err.(*entities.Error); ok {
			log.Logger.Error("SetDefaultTag WithTransaction() failed", zap.String("Error", txnErr.Error()))
			return txnErr
		} else {
			log.Logger.Error("SetDefaultTag WithTransaction() failed, assert txnError", zap.String("Error", err.Error()))
			return entities.NewError(entities.ErrCodeDatastoreFailed)
		}
	}
	return nil
}

func RemoveDefaultTag(ctx context.Context, ruleEngineName string) *entities.Error {

	removeDefaultTagTxnFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		existingEngine, err := getRuleEngine(sessCtx, ruleEngineName)
		if err != nil {
			log.Logger.Error("Get RuleEngine failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		if existingEngine == nil {
			return nil, entities.NewError(entities.ErrCodeRuleEngineNotFound)
		}

		existingEngine.LastUpdateTime = time.Now().Unix()
		existingEngine.DefaultTag = ""

		if err := upsertRuleEngine(sessCtx, existingEngine); err != nil {
			log.Logger.Error("Upsert RuleEngine failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		return nil, nil
	}

	session, err := client.StartSession()
	if err != nil {
		log.Logger.Error("RemoveDefaultTag StartSession() failed", zap.String("Error", err.Error()))
		return entities.NewError(entities.ErrCodeDatastoreFailed)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, removeDefaultTagTxnFunc)
	if err != nil {
		if txnErr, ok := err.(*entities.Error); ok {
			log.Logger.Error("RemoveDefaultTag WithTransaction() failed", zap.String("Error", txnErr.Error()))
			return txnErr
		} else {
			log.Logger.Error("RemoveDefaultTag WithTransaction() failed, assert txnError failed", zap.String("Error", err.Error()))
			return entities.NewError(entities.ErrCodeDatastoreFailed)
		}
	}
	return nil
}

func EnableTag(ctx context.Context, ruleEngineName string, tag string) *entities.Error {

	enableTagTxnFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
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
			return nil, nil
		}

		existingEngine.LastUpdateTime = time.Now().Unix()
		t.IsEnable = true

		if err := upsertRuleEngine(sessCtx, existingEngine); err != nil {
			log.Logger.Error("EnableTag failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}
		return nil, nil
	}

	session, err := client.StartSession()
	if err != nil {
		log.Logger.Error("EnableTag StartSession() failed", zap.String("Error", err.Error()))
		return entities.NewError(entities.ErrCodeDatastoreFailed)
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, enableTagTxnFunc)
	if err != nil {
		if txnErr, ok := err.(*entities.Error); ok {
			log.Logger.Error("EnableTag WithTransaction() failed", zap.String("Error", txnErr.Error()))
			return txnErr
		} else {
			log.Logger.Error("EnableTag WithTransaction() failed, assert txnError failed", zap.String("Error", err.Error()))
			return entities.NewError(entities.ErrCodeDatastoreFailed)
		}
	}
	return nil
}

func DisableTag(ctx context.Context, ruleEngineName string, tag string) *entities.Error {

	disableTagTxnFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
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
		if !t.IsEnable {
			return nil, nil
		}

		if existingEngine.DefaultTag == tag {
			return nil, entities.NewError(entities.ErrCodeTagDisableNotAllowed)
		}

		existingEngine.LastUpdateTime = time.Now().Unix()
		t.IsEnable = false

		if err := upsertRuleEngine(sessCtx, existingEngine); err != nil {
			log.Logger.Error("DisableTag failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}
		return nil, nil
	}

	session, err := client.StartSession()
	if err != nil {
		log.Logger.Error("DisableTag StartSession() failed", zap.String("Error", err.Error()))
		return entities.NewError(entities.ErrCodeDatastoreFailed)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, disableTagTxnFunc)
	if err != nil {
		if txnErr, ok := err.(*entities.Error); ok {
			return txnErr
		} else {
			log.Logger.Error("assert txn error failed", zap.String("Error", fmt.Sprintf("%+v", err)))
			return entities.NewError(entities.ErrCodeDatastoreFailed)
		}
	}

	return nil
}

func GetCompleteRuleEngine(ctx context.Context, ruleEngineName string) (*entities.CompleteRuleEngine, *entities.Error) {

	getCompleteRuleEngineTxnFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		existingEngine, err := getRuleEngine(sessCtx, ruleEngineName)
		if err != nil {
			log.Logger.Error("Get RuleEngine failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}

		if existingEngine == nil {
			return nil, entities.NewError(entities.ErrCodeRuleEngineNotFound)
		}

		result := &entities.CompleteRuleEngine{
			Name:       existingEngine.Name,
			DefaultTag: existingEngine.DefaultTag,
			Tags:       map[string]*entities.TagResponse{},
		}

		for tag, t := range existingEngine.Tags {
			if config, err := getRuleEngineConfig(sessCtx, t.EngineConfigID); err != nil {
				log.Logger.Error("Get RuleEngineConfig failed", zap.String("Error", err.Error()))
				return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
			} else {
				result.Tags[tag] = &entities.TagResponse{
					IsEnable: t.IsEnable,
					Config:   config,
				}
			}
		}

		return result, nil
	}

	session, err := client.StartSession()
	if err != nil {
		log.Logger.Error("GetCompleteRuleEngine StartSession() failed", zap.String("Error", err.Error()))
		return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
	}
	defer session.EndSession(ctx)

	txnResult, err := session.WithTransaction(ctx, getCompleteRuleEngineTxnFunc)
	if err != nil {
		if txnErr, ok := err.(*entities.Error); ok {
			log.Logger.Error("GetCompleteRuleEngine WithTransaction() failed", zap.String("Error", txnErr.Error()))
			return nil, txnErr
		} else {
			log.Logger.Error("GetCompleteRuleEngine WithTransaction() failed, assert txnError failed", zap.String("Error", err.Error()))
			return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
		}
	}
	if txnResult == nil {
		return nil, nil
	}

	if result, ok := txnResult.(*entities.CompleteRuleEngine); ok {
		return result, nil
	} else {
		log.Logger.Error("GetCompleteRuleEngine assert txnResult failed", zap.String("txnResult", fmt.Sprintf("%+v", txnResult)))
		return nil, entities.NewError(entities.ErrCodeDatastoreFailed)
	}
}

func upsertRuleEngine(ctx context.Context, ruleEngine *entities.RuleEngine) error {
	filter := bson.D{{Key: "name", Value: ruleEngine.Name}}
	opts := options.Replace().SetUpsert(true)
	_, err := ruleEngineCollection.ReplaceOne(ctx, filter, ruleEngine, opts)
	if err != nil {
		return err
	}

	return nil
}

func getRuleEngine(ctx context.Context, ruleEngineName string) (*entities.RuleEngine, error) {
	var ruleEngine entities.RuleEngine
	err := ruleEngineCollection.FindOne(ctx, bson.D{{Key: "name", Value: ruleEngineName}}).Decode(&ruleEngine)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &ruleEngine, nil
}
