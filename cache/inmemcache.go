package cache

import (
	"time"
)

// InMemCache is caching mechanism using redis db
type InMemCache struct {
	GetData func(string) (interface{}, error)
	Data    map[string]Object
}

// NewInMemCache creates new InMemCache object
func NewInMemCache(getData DataFunc) Cache {
	return &InMemCache{
		GetData: getData,
		Data:    make(map[string]Object),
	}
}

// GetDataFromCache will gets data from redis and if it expierd will fetch from the exact source
func (i *InMemCache) GetDataFromCache(key string) (interface{}, error) {
	object := i.Data[key]
	if time.Since(object.LastUpdatedTime) > time.Minute*5 {
		return i.GetDataFromSource(key)
	}
	return object.Data, nil
}

// GetDataFromSource will gets from the exact source
func (i *InMemCache) GetDataFromSource(key string) (interface{}, error) {
	data, err := i.GetData(key)
	if err != nil {
		return nil, err
	}
	object := Object{
		Data:            data,
		LastUpdatedTime: time.Now(),
	}

	i.Data[key] = object

	return data, nil

}
