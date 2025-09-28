package cache

import "github.com/aidapedia/gdk/cache/engine"

type Cache struct {
	engine.Interface
}

func NewCache(cli engine.Interface) *Cache {
	return &Cache{
		Interface: cli,
	}
}
