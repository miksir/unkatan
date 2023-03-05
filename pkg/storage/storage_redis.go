package storage

import (
	"context"
	"github.com/miksir/unkatan/pkg/lconfig"
	zlog "github.com/miksir/unkatan/pkg/log"
	goredis "github.com/redis/go-redis"
	"go.uber.org/zap"
	"math"
	"time"
)

type RedisRegistry struct {
	logger zlog.Logger
	cfg    lconfig.Reader
	redis  *goredis.Client
}

func NewRedisRegistry(cfg lconfig.Reader, log zlog.Logger) RedisRegistry {
	cfg.SetDefault("addr", "127.0.0.1:6397")
	return RedisRegistry{
		logger: log,
		cfg:    cfg,
	}
}

func (str RedisRegistry) RestoreKatanState() ([]byte, error) {
	ctx := context.Background()
	redis := str.getRedis(ctx)

	cmd := redis.Get("status")
	data, err := cmd.Bytes()

	if err != nil {
		str.logger.Error(ctx, "redis get status error", zap.Error(err))
	}

	return data, err

}

func (str RedisRegistry) SaveKatanState(data []byte) error {
	ctx := context.Background()
	redis := str.getRedis(ctx)

	cmd := redis.Set("status", data, -1)
	if cmd.Err() != nil {
		str.logger.Error(ctx, "redis set status error", zap.Error(cmd.Err()))
		return cmd.Err()
	}

	str.logger.Debug(
		ctx,
		"SaveKatanState",
		zap.String("status", string(data)),
	)

	return nil
}

func (str RedisRegistry) getRedis(ctx context.Context) *goredis.Client {
	if str.redis != nil {
		return str.redis
	}

	str.redis = goredis.NewClient(&goredis.Options{
		Network: "tcp",
		Addr:    str.cfg.GetString("addr"),
		Dialer:  nil,
		OnConnect: func(conn *goredis.Conn) error {
			str.logger.Info(ctx, "redis: connected")
			return nil
		},
		Password:           str.cfg.GetString("password"),
		DB:                 0,
		MaxRetries:         math.MaxInt64,
		MinRetryBackoff:    1 * time.Millisecond,
		MaxRetryBackoff:    5 * time.Second,
		DialTimeout:        1500 * time.Millisecond,
		ReadTimeout:        1500 * time.Millisecond,
		WriteTimeout:       1500 * time.Millisecond,
		PoolSize:           10,
		MinIdleConns:       1,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        30 * time.Second,
		IdleCheckFrequency: 10 * time.Second,
		TLSConfig:          nil,
	})

	return str.redis
}
