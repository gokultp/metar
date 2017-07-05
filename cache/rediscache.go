package cache

import (
	"fmt"
	"time"

	"encoding/json"

	"github.com/go-redis/redis"
)

// RedisCache is caching mechanism using redis db
type RedisCache struct {
	Client  *redis.Client
	GetData func(string) (interface{}, error)
}

// NewRedisCache will create new Redis cache instance
func NewRedisCache(client *redis.Client, getData func(string) (interface{}, error)) Cache {
	return &RedisCache{
		Client:  client,
		GetData: getData,
	}
}

// GetDataFromCache will gets data from redis and if it expierd will fetch from the exact source
func (r *RedisCache) GetDataFromCache(key string) (interface{}, error) {
	var data []byte
	var object Object
	err := r.Client.Get(key).Scan(&data)

	err = json.Unmarshal(data, &object)

	fmt.Println(err, object)
	if err != nil || time.Since(object.LastUpdatedTime) > time.Minute*5 {
		return r.GetDataFromSource(key)
	}
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
	err = r.Client.Set(key, string(serialised), time.Minute*5).Err()
	fmt.Println(err)

	return data, nil

}
