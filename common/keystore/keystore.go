package keystore

import (
	"common/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	Client *redis.Client
)

func Connect() {
	cacheAddr := fmt.Sprintf("%s:%s", commonconfig.KeyStoreHost, commonconfig.KeyStorePort)

	Client = redis.NewClient(&redis.Options{
		Addr:     cacheAddr,
		Password: commonconfig.KeyStorePassword,
		DB:       commonconfig.KeyStoreDb,
	})

	_, err := Client.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("Error connecting to cache client: %v", err)
	} else {
		log.Printf("Connected to cache")
	}
}

func Close() {
	log.Printf("Disconnected from cache")
	Client.Close()
}

func Get(key string, group string) (string, error) {
	val, err := Client.Get(context.Background(), fmt.Sprintf("%s:%s", group, key)).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

const (
	GlobalGroup = "global"
)
