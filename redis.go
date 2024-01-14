package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisProfile struct {
	Address  string
	Password string
	DB       int
	Client   *redis.Client
}

func (profile RedisProfile) SetupRedis() (*redis.Client, error) {
	log.Println("Redis Setup")
	rdb := redis.NewClient(&redis.Options{
		Addr:     profile.Address,
		Password: profile.Password,
		DB:       profile.DB,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func (profile RedisProfile) Subscribe(channel string, handler func(message string)) error {
	ch := profile.Client.Subscribe(context.Background(), channel)
	defer ch.Unsubscribe(context.Background(), channel)

	for msg := range ch.Channel() {
		handler(msg.Payload)
	}

	return nil
}

func (profile RedisProfile) Publish(channel string, message string) error {
	return profile.Client.Publish(context.Background(), channel, message).Err()
}
