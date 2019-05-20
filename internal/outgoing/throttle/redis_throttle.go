package throttles

import "time"

// RedisThrottle - Throttle by redis
type RedisThrottle struct {
}

// NewRedisThrottle - Create redis throttle
func NewRedisThrottle(store interface{}, period time.Duration, trailing bool) *RedisThrottle {
	return &RedisThrottle{}
}
