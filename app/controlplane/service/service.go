package service

import (
	"context"

	ruleenginecore "github.com/niharrathod/ruleengine-core"
	"github.com/niharrathod/ruleengine/app/entities"
	"github.com/niharrathod/ruleengine/app/ext/datastore"
	"github.com/niharrathod/ruleengine/app/validator"
)

func CreateRuleEngine(ctx context.Context, ruleEngineName string, tag string, config *ruleenginecore.RuleEngineConfig) *entities.Error {
	if !validator.IsAlphanumericMax30(ruleEngineName) {
		return entities.NewError(entities.ErrCodeInvalidRuleEngineName)
	}
	if !validator.IsAlphanumericMax30(tag) {
		return entities.NewError(entities.ErrCodeInvalidTagName)
	}
	if err := config.Validate(); err != nil {
		return entities.NewErrorWithMsg(entities.ErrCodeInvalidRuleEngineConfig, err.Error())
	}

	if err := datastore.CreateRuleEngine(ctx, ruleEngineName, tag, config); err != nil {
		return err
	}
	return nil
}

func DeleteRuleEngine(ctx context.Context, ruleEngineName string) *entities.Error {
	if !validator.IsAlphanumericMax30(ruleEngineName) {
		return entities.NewError(entities.ErrCodeInvalidRuleEngineName)
	}

	if err := datastore.DeleteRuleEngine(ctx, ruleEngineName); err != nil {
		return err
	}
	return nil
}

func DeleteRuleEngineConfig(ctx context.Context, ruleEngineName string, tag string) *entities.Error {
	if !validator.IsAlphanumericMax30(ruleEngineName) {
		return entities.NewError(entities.ErrCodeInvalidRuleEngineName)
	}
	if !validator.IsAlphanumericMax30(tag) {
		return entities.NewError(entities.ErrCodeInvalidTagName)
	}

	if err := datastore.DeleteRuleEngineConfig(ctx, ruleEngineName, tag); err != nil {
		return err
	}
	return nil
}

func GetCompleteRuleEngine(ctx context.Context, ruleEngineName string) (*entities.CompleteRuleEngine, *entities.Error) {
	if !validator.IsAlphanumericMax30(ruleEngineName) {
		return nil, &entities.Error{ErrCode: entities.ErrCodeInvalidRuleEngineName}
	}

	if result, err := datastore.GetCompleteRuleEngine(ctx, ruleEngineName); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func SetDefaultTag(ctx context.Context, ruleEngineName string, tag string) *entities.Error {
	if !validator.IsAlphanumericMax30(ruleEngineName) {
		return entities.NewError(entities.ErrCodeInvalidRuleEngineName)
	}
	if !validator.IsAlphanumericMax30(tag) {
		return entities.NewError(entities.ErrCodeInvalidTagName)
	}

	if err := datastore.SetDefaultTag(ctx, ruleEngineName, tag); err != nil {
		return err
	}
	return nil
}

func RemoveDefaultTag(ctx context.Context, ruleEngineName string) *entities.Error {
	if !validator.IsAlphanumericMax30(ruleEngineName) {
		return entities.NewError(entities.ErrCodeInvalidRuleEngineName)
	}

	if err := datastore.RemoveDefaultTag(ctx, ruleEngineName); err != nil {
		return err
	}
	return nil
}

func EnableRuleEngine(ctx context.Context, ruleEngineName string, tag string) *entities.Error {
	if !validator.IsAlphanumericMax30(ruleEngineName) {
		return entities.NewError(entities.ErrCodeInvalidRuleEngineName)
	}
	if !validator.IsAlphanumericMax30(tag) {
		return entities.NewError(entities.ErrCodeInvalidTagName)
	}

	if err := datastore.EnableTag(ctx, ruleEngineName, tag); err != nil {
		return err
	}
	return nil
}

func DisableRuleEngine(ctx context.Context, ruleEngineName string, tag string) *entities.Error {
	if !validator.IsAlphanumericMax30(ruleEngineName) {
		return entities.NewError(entities.ErrCodeInvalidRuleEngineName)
	}
	if !validator.IsAlphanumericMax30(tag) {
		return entities.NewError(entities.ErrCodeInvalidTagName)
	}

	if err := datastore.DisableTag(ctx, ruleEngineName, tag); err != nil {
		return err
	}
	return nil
}
