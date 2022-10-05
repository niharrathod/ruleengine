package log

import (
	"github.com/niharrathod/ruleengine/app/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Initialize() {
	var loggerConfig zap.Config

	if config.EnvironmentMode == config.ProductionMode {
		loggerConfig = zap.NewProductionConfig()

	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	Logger = zap.Must(loggerConfig.Build())
}
