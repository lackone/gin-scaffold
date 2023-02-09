package services

import (
	"context"
	"errors"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/internal/provider/redis"
	redis9 "github.com/redis/go-redis/v9"
	"strings"
	"sync"
	"time"
)

type RedisService struct {
	container core.IContainer
	client    *redis9.Client
	lock      sync.RWMutex
}

func NewRedisService(params ...interface{}) (interface{}, error) {
	container := params[0].(core.IContainer)

	if !container.IsBind(contract.RedisKey) {
		err := container.Bind(&redis.RedisProvider{})
		if err != nil {
			return nil, err
		}
	}

	redisService := container.MustMake(contract.RedisKey).(contract.IRedis)
	client, err := redisService.GetClient(redis.WithConfigPath("cache"))
	if err != nil {
		return nil, err
	}

	obj := &RedisService{
		container: container,
		client:    client,
		lock:      sync.RWMutex{},
	}
	return obj, nil
}

func (r *RedisService) Get(ctx context.Context, key string) (string, error) {
	cmd := r.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redis9.Nil) {
		return "", ErrKeyNotFound
	}
	return cmd.Result()
}

func (r *RedisService) GetObj(ctx context.Context, key string, model interface{}) error {
	cmd := r.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redis9.Nil) {
		return ErrKeyNotFound
	}
	err := cmd.Scan(model)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisService) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	pipeline := r.client.Pipeline()
	vals := make(map[string]string)
	cmds := make([]*redis9.StringCmd, 0, len(keys))

	for _, key := range keys {
		cmds = append(cmds, pipeline.Get(ctx, key))
	}
	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	errs := make([]string, 0, len(keys))
	for _, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		key := cmd.Args()[1].(string)
		vals[key] = val
	}
	if len(errs) == 0 {
		return vals, nil
	}
	return vals, errors.New(strings.Join(errs, "||"))
}

func (r *RedisService) Set(ctx context.Context, key string, val string, timeout time.Duration) error {
	return r.client.Set(ctx, key, val, timeout).Err()
}

func (r *RedisService) SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return r.client.Set(ctx, key, val, timeout).Err()
}

func (r *RedisService) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	pipline := r.client.Pipeline()
	cmds := make([]*redis9.StatusCmd, 0, len(data))
	for k, v := range data {
		cmds = append(cmds, pipline.Set(ctx, k, v, timeout))
	}
	_, err := pipline.Exec(ctx)
	return err
}

func (r *RedisService) SetForever(ctx context.Context, key string, val string) error {
	return r.client.Set(ctx, key, val, NoneDuration).Err()
}

func (r *RedisService) SetForeverObj(ctx context.Context, key string, val interface{}) error {
	return r.client.Set(ctx, key, val, NoneDuration).Err()
}

func (r *RedisService) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	return r.client.Expire(ctx, key, timeout).Err()
}

func (r *RedisService) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *RedisService) Remember(ctx context.Context, key string, timeout time.Duration, rememberFunc contract.RememberFunc, obj interface{}) error {
	err := r.GetObj(ctx, key, obj)
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return err
	}

	objNew, err := rememberFunc(ctx, r.container)
	if err != nil {
		return err
	}

	if err := r.SetObj(ctx, key, objNew, timeout); err != nil {
		return err
	}
	if err := r.GetObj(ctx, key, obj); err != nil {
		return err
	}
	return nil
}

func (r *RedisService) Calc(ctx context.Context, key string, step int64) (int64, error) {
	return r.client.IncrBy(ctx, key, step).Result()
}

func (r *RedisService) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.IncrBy(ctx, key, 1).Result()
}

func (r *RedisService) Decrement(ctx context.Context, key string) (int64, error) {
	return r.client.IncrBy(ctx, key, -1).Result()
}

func (r *RedisService) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisService) DelMany(ctx context.Context, keys []string) error {
	pipline := r.client.Pipeline()
	cmds := make([]*redis9.IntCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipline.Del(ctx, key))
	}
	_, err := pipline.Exec(ctx)
	return err
}
