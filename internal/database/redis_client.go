package database

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	redisAddr, found := os.LookupEnv("REDIS_ADDR");
	if !found {
		return nil, errors.New("error loading env")
	}
	redisPassword, found := os.LookupEnv("REDIS_PASS");
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
        Addr: redisAddr,
        Password: redisPassword, 
        DB: redisDB,  
				Protocol: redisProtocol,              
	})
	if rdb == nil {
		return nil, errors.New("error connecting to redis client")
	}
	fmt.Println("connected redis client")
	return rdb, nil
}
