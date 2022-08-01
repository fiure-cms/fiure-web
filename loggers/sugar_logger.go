package loggers

import "go.uber.org/zap"

var Sugar *zap.SugaredLogger

func SetupSugarLogger(serviceName string, serviceVersion string) {
	// Zap Logger Init
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	Sugar = logger.Sugar()
	Sugar.Named(serviceName + ":" + serviceVersion)
}
