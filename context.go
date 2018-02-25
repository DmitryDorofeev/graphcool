package graphcool

import (
	"net/http"
	"net/url"
	"sync"
)

type Context struct {
	mu      sync.Mutex
	Request *http.Request
	Writer  http.ResponseWriter
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

func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}
