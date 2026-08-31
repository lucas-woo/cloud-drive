package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context) (*redis.Client, error) {
	redisAddr, found := os.LookupEnv("REDIS_ADDR")
	if !found {
		return nil, errors.New("error loading env")
	}

	redisPassword, found := os.LookupEnv("REDIS_PASS")
	if !found {
		return nil, errors.New("error loading env")
	}

	redisDBString, found := os.LookupEnv("REDIS_DB")
	if !found {
		return nil, errors.New("error loading env")
	}

	redisDB, err := strconv.Atoi(redisDBString)
	if err != nil {
		return nil, err
	}

	redisProtocolString, found := os.LookupEnv("REDIS_PROTOCOL")
	if !found {
		return nil, errors.New("error loading env")
	}

	redisProtocol, err := strconv.Atoi(redisProtocolString)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
		Protocol: redisProtocol,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("error pinging redis: %w", err)
	}

	fmt.Println("connected to redis successfully")

	return rdb, nil
}