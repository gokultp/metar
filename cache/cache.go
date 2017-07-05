package cache

import "time"

// Cache interface will have cache functions
type Cache interface {
	GetDataFromCache(key string) (interface{}, error)
	GetDataFromSource(key string) (interface{}, error)
}

// Object encapsulates payload and lastupdated time
type Object struct {
	LastUpdatedTime time.Time
	Data            interface{}
}
