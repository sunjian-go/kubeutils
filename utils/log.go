package utils

import (
	"go.uber.org/zap"
)

var Log logg

type logg struct {
}

var Logg *zap.Logger

func (l *logg) LogInit() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	//logger.Info("Development")
	Logg = logger
	return nil
}
