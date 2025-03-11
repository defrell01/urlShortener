package db

import (
	"context"
	"fmt"
	"time"
	"urlshortener/configs"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultTTL = 3600
)

type Redis struct {
	RedisClient *redis.Client
	ctx         context.Context
}

func NewRedisClient(conf *configs.Config) *Redis {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       int(conf.Redis.DB),
	})

	ctx := context.Background()
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("❌ Redis connection error: %v\n", err)
		return nil
	}

	fmt.Println("✅ Redis connected:", pong)

	return &Redis{
		RedisClient: redisClient,
		ctx:         ctx,
	}
}

func (r *Redis) SetCache(key string, value string, ttl time.Duration) error {
	fmt.Println("setting cache")
	return r.RedisClient.Set(r.ctx, key, value, ttl).Err()
}

func (r *Redis) GetCache(key string) (string, error) {
	return r.RedisClient.Get(r.ctx, key).Result()
}

func (r *Redis) DeleteCache(key string) error {
	res, err := r.GetCache(key)
	if err != nil {
		return err
	}
	if len(res) != 0 {
		_, err := r.RedisClient.Del(r.ctx, key).Result()
		return err
	}
	return nil
}

func (r *Redis) UpdateCache(key string) error {
	res, err := r.GetCache(key)
	if err != nil {
		return err
	}
	if len(res) != 0 {
		return r.SetCache(key, res, time.Second*DefaultTTL)
	}
	return nil
}
