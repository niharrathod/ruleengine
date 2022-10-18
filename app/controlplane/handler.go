package controlplane

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ruleenginecore "github.com/niharrathod/ruleengine-core"
	"github.com/niharrathod/ruleengine/app/controlplane/service"
	"github.com/niharrathod/ruleengine/app/entities"
	"github.com/niharrathod/ruleengine/app/log"
	"go.uber.org/zap/zapcore"
)

func CreateRuleEngine() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ruleEngineName := ctx.Param("ruleengine")
		tag := ctx.Param("tag")
		var config ruleenginecore.RuleEngineConfig
		if err := ctx.BindJSON(&config); err != nil {
			log.Logger.Error("Could not unmarshal the patient as body", zapcore.Field{Key: "Error", String: err.Error()})
			setResponse(ctx, entities.NewError(entities.ErrCodeParsingFailed))
			return
		}
		if err := service.CreateRuleEngine(ctx, ruleEngineName, tag, &config); err != nil {
			setResponse(ctx, err)
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func GetRuleEngine() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ruleEngineName := ctx.Param("ruleengine")

		ruleEngine, err := service.GetCompleteRuleEngine(ctx, ruleEngineName)
		if err != nil {
			setResponse(ctx, err)
		}

		ctx.JSON(http.StatusOK, ruleEngine)
	}
}

func DeleteRuleEngine() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ruleEngineName := ctx.Param("ruleengine")
		if err := service.DeleteRuleEngine(ctx, ruleEngineName); err != nil {
			setResponse(ctx, err)
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func DeleteRuleEngineConfig() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ruleEngineName := ctx.Param("ruleengine")
		tag := ctx.Param("tag")
		if err := service.DeleteRuleEngineConfig(ctx, ruleEngineName, tag); err != nil {
			setResponse(ctx, err)
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func SetDefaultTag() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ruleEngineName := ctx.Param("ruleengine")
		tag := ctx.Param("tag")
		if err := service.SetDefaultTag(ctx, ruleEngineName, tag); err != nil {
			setResponse(ctx, err)
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func RemoveDefaultTag() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ruleEngineName := ctx.Param("ruleengine")
		if err := service.RemoveDefaultTag(ctx, ruleEngineName); err != nil {
			setResponse(ctx, err)
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func EnableRuleEngine() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ruleEngineName := ctx.Param("ruleengine")
		tag := ctx.Param("tag")
		if err := service.EnableRuleEngine(ctx, ruleEngineName, tag); err != nil {
			setResponse(ctx, err)
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func DisableRuleEngine() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ruleEngineName := ctx.Param("ruleengine")
		tag := ctx.Param("tag")
		if err := service.DisableRuleEngine(ctx, ruleEngineName, tag); err != nil {
			setResponse(ctx, err)
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func setResponse(ctx *gin.Context, err *entities.Error) {
	switch err.ErrCode {
	case entities.ErrCodeRuleEngineNotFound,
		entities.ErrCodeTagNotFound:
		ctx.JSON(http.StatusNotFound, err)
		return
	case entities.ErrCodeParsingFailed,
		entities.ErrCodeInvalidRuleEngineName,
		entities.ErrCodeInvalidTagName,
		entities.ErrCodeInvalidRuleEngineConfig,
		entities.ErrCodeTagDeleteNotAllowed,
		entities.ErrCodeTagDisableNotAllowed,
		entities.ErrCodeDefaultTagExistAndMustBeEnabled,
		entities.ErrCodeTagAlreadyExist:
		ctx.JSON(http.StatusBadRequest, err)
		return
	case entities.ErrCodeDatastoreFailed:
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusInternalServerError, err)
}
