package cache

import (
	"time"

	"encoding/json"

	"fmt"

	"gopkg.in/redis.v3"
)

// RedisCache is caching mechanism using redis db
type RedisCache struct {
	Client  *redis.Client
	GetData func(string) (interface{}, error)
}

// NewRedisCache will create new Redis cache instance
func NewRedisCache(getData DataFunc, redisURL, redisPassword string, redisDB int) Cache {
	return &RedisCache{
		Client: redis.NewClient(&redis.Options{
			Addr:     redisURL,
			Password: redisPassword,
			DB:       int64(redisDB),
		}),
		GetData: getData,
	}
}

// GetDataFromCache will gets data from redis and if it expierd will fetch from the exact source
func (r *RedisCache) GetDataFromCache(key string) (interface{}, error) {
	var data []byte
	var object Object
	err := r.Client.Get(key).Scan(&data)
	err = json.Unmarshal(data, &object)
	if err != nil || time.Since(object.LastUpdatedTime) > time.Minute*5 {
		return r.GetDataFromSource(key)
	}
	fmt.Println("From Cache")

	return object.Data, nil
}

// GetDataFromSource will gets from the exact source
func (r *RedisCache) GetDataFromSource(key string) (interface{}, error) {
	data, err := r.GetData(key)
	if err != nil {
		return nil, err
	}
	object := Object{
		Data:            data,
		LastUpdatedTime: time.Now(),
	}
	// execute redis set in parallel
	serialised, err := json.Marshal(object)
	go r.Client.Set(key, string(serialised), time.Minute*5)
	return data, nil

}
