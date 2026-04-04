package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

const introspectCacheMaxEntries = 4096

type introspectionResult struct {
	Active   bool   `json:"active"`
	Sub      string `json:"sub"`
	Username string `json:"username"`
	Scope    string `json:"scope"`
}

type introspectCache struct {
	mu   sync.RWMutex
	ttl  time.Duration
	data map[string]introCacheEntry
}

type introCacheEntry struct {
	res introspectionResult
	exp time.Time
}

func tokenCacheKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func newIntrospectCache(ttl time.Duration) *introspectCache {
	if ttl <= 0 {
		return nil
	}
	return &introspectCache{ttl: ttl, data: make(map[string]introCacheEntry)}
}

func (c *introspectCache) get(key string) (introspectionResult, bool) {
	if c == nil {
		return introspectionResult{}, false
	}
	now := time.Now()
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()
	if !ok || now.After(e.exp) {
		return introspectionResult{}, false
	}
	return e.res, true
}

func (c *introspectCache) set(key string, res introspectionResult) {
	if c == nil {
		return
	}
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.data) >= introspectCacheMaxEntries {
		c.sweepExpiredLocked(now)
	}
	for len(c.data) >= introspectCacheMaxEntries {
		for k := range c.data {
			delete(c.data, k)
			break
		}
	}
	c.data[key] = introCacheEntry{res: res, exp: now.Add(c.ttl)}
}

func (c *introspectCache) sweepExpiredLocked(now time.Time) {
	for k, v := range c.data {
		if now.After(v.exp) {
			delete(c.data, k)
		}
	}
}
