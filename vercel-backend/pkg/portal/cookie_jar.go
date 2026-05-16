package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCookieJar implements http.CookieJar interface with Redis persistence
type RedisCookieJar struct {
	redisClient *redis.Client
	key         string
	ttl         time.Duration
}

func NewRedisCookieJar(client *redis.Client, username string) *RedisCookieJar {
	return &RedisCookieJar{
		redisClient: client,
		key:         fmt.Sprintf("portal:session:v2:%s", username),
		ttl:         24 * time.Hour,
	}
}

// SetCookies handles saving cookies to Redis
func (jar *RedisCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if jar.redisClient == nil {
		return
	}


	// Load existing cookies to merge
	existing := jar.loadFromRedis()
	cookieMap := make(map[string]*http.Cookie)
	for _, c := range existing {
		cookieMap[c.Name] = c
	}
	
	for _, c := range cookies {
		// If MaxAge < 0, it means delete the cookie
		if c.MaxAge < 0 {
			delete(cookieMap, c.Name)
		} else {
			cookieMap[c.Name] = c
		}
	}

	merged := make([]*http.Cookie, 0, len(cookieMap))
	for _, c := range cookieMap {
		merged = append(merged, c)
	}

	jar.saveToRedis(merged)
}

// Cookies handles retrieving cookies from Redis
func (jar *RedisCookieJar) Cookies(u *url.URL) []*http.Cookie {
	if jar.redisClient == nil {
		return nil
	}

	// Only return cookies for portal domain
	if !jar.isPortalDomain(u.Host) {
		return nil
	}

	cookies := jar.loadFromRedis()
	
	// Ensure Domain and Path are set if empty to make net/http client send them
	for _, c := range cookies {
		if c.Domain == "" {
			c.Domain = u.Host
		}
		if c.Path == "" {
			c.Path = "/"
		}
	}

	return cookies
}

func (jar *RedisCookieJar) isPortalDomain(host string) bool {
	return strings.Contains(host, "vncdc.gov.vn")
}

func (jar *RedisCookieJar) saveToRedis(cookies []*http.Cookie) {
	if len(cookies) == 0 {
		return
	}

	data, err := json.Marshal(cookies)
	if err != nil {
		return
	}

	jar.redisClient.Set(context.Background(), jar.key, data, jar.ttl)
}

func (jar *RedisCookieJar) loadFromRedis() []*http.Cookie {
	val, err := jar.redisClient.Get(context.Background(), jar.key).Result()
	if err != nil {
		return nil
	}

	var cookies []*http.Cookie
	if err := json.Unmarshal([]byte(val), &cookies); err != nil {
		return nil
	}

	// Refresh TTL
	jar.redisClient.Expire(context.Background(), jar.key, jar.ttl)

	return cookies
}
