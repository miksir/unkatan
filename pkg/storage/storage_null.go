package storage

import (
	zlog "github.com/miksir/unkatan/pkg/log"
	"go.uber.org/zap"
)

type NullRegistry struct {
	logger zlog.Logger
}

func NewNullRegistry(log zlog.Logger) NullRegistry {
	return NullRegistry{
		logger: log,
	}
}

func (str NullRegistry) RestoreKatanState() ([]byte, error) {
	str.logger.Info(nil, "RestoreKatanState")
	return []byte{}, nil
}

func (str NullRegistry) SaveKatanState(data []byte) error {
	str.logger.Info(nil, "SaveKatanState", zap.String("data", string(data)))
	return nil
}
