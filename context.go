package graphcool

import (
	"net/http"
	"sync"
)

type Context struct {
	mu      sync.Mutex
	Request *http.Request
	Keys    map[string]interface{}
	Errors  []*QueryError
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.Lock()
	value, exists = c.Keys[key]
	c.mu.Unlock()
	return
}

func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.mu.Lock()
	c.Keys[key] = value
	c.mu.Unlock()
}
