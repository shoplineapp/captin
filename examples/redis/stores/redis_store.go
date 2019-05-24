package stores

import (
	"fmt"
	"time"

	interfaces "github.com/shoplineapp/captin/interfaces"
	"github.com/shoplineapp/captin/models"

	lock "github.com/bsm/redis-lock"
	"github.com/go-redis/redis"
)

// RedisStore - Redis data store
type RedisStore struct {
	interfaces.StoreInterface
	redisClient *redis.Client
	locker      *lock.Locker
}

// NewRedisStore - Create new RedisStore
func NewRedisStore(addr string) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	client.Set("test", "value", time.Hour).Err()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Create a new locker with default settings
	locker := lock.New(client, "captin.locker.lock", nil)

	return &RedisStore{
		redisClient: client,
		locker:      locker,
	}
}

func (rs *RedisStore) ping() {
	val, _ := rs.redisClient.Ping().Result()
	fmt.Println(val)
}

// Get - Get value from store, return with remaining time
func (rs RedisStore) Get(key string) (string, bool, time.Duration, error) {
	_, redisLockErr := rs.locker.Lock()
	if redisLockErr != nil {
		fmt.Println("[RedisStore] Redis Lock Error: ", redisLockErr)
		return "", false, time.Duration(0), redisLockErr
	}
	defer rs.locker.Unlock()
	fmt.Println("[RedisStore] Get Key: ", key)
	timeRemain, err := rs.redisClient.TTL(key).Result()

	if err != nil {
		if err == redis.Nil {
			// Key does not exist
			fmt.Println("[RedisStore] Key does not exist")
			return "", false, time.Duration(0), nil
		}
		return "", false, time.Duration(0), err
	}

	val, valErr := rs.redisClient.Get(key).Result()
	if valErr != nil {
		if valErr == redis.Nil {
			// Key does not exist
			fmt.Println("[RedisStore] Key does not exist")
			return "", false, time.Duration(0), nil
		}
		return "", false, time.Duration(0), valErr
	}

	return val, true, timeRemain, nil
}

// Set - Set value into store with ttl
func (rs RedisStore) Set(key string, value string, ttl time.Duration) (bool, error) {
	_, redisLockErr := rs.locker.Lock()
	if redisLockErr != nil {
		fmt.Println("[RedisStore] Redis Lock Error: ", redisLockErr)
		return false, redisLockErr
	}
	defer rs.locker.Unlock()

	fmt.Println("[RedisStore] Set Key: ", key)
	err := rs.redisClient.Set(key, value, ttl).Err()
	if err != nil {
		fmt.Println("[RedisStore] Set Key Error: ", err)
		return false, err
	}
	val, _, _, _ := rs.Get(key)
	fmt.Println("[RedisStore] Set Key Value: ", val)
	return true, nil
}

// Update - Update value for key
func (rs RedisStore) Update(key string, value string) (bool, error) {
	_, redisLockErr := rs.locker.Lock()
	if redisLockErr != nil {
		fmt.Println("[RedisStore] Redis Lock Error: ", redisLockErr)
		return false, redisLockErr
	}
	defer rs.locker.Unlock()

	fmt.Println("[RedisStore] Update Key: ", key)
	timeRemain, err := rs.redisClient.TTL(key).Result()

	if err != nil {
		return false, err
	}

	err = rs.redisClient.Set(key, value, timeRemain).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

// Remove - Remove value for key
func (rs RedisStore) Remove(key string) (bool, error) {
	_, redisLockErr := rs.locker.Lock()
	if redisLockErr != nil {
		fmt.Println("[RedisStore] Redis Lock Error: ", redisLockErr)
		return false, redisLockErr
	}
	defer rs.locker.Unlock()

	fmt.Println("[RedisStore] Remove Key: ", key)
	err := rs.redisClient.Del(key).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

// DataKey - Generate data key
func (rs RedisStore) DataKey(e models.IncomingEvent, dest models.Destination, prefix string, suffix string) string {
	return fmt.Sprintf("%s%s:%s:%s%s", prefix, e.Key, dest.Config.Name, e.TargetId, suffix)
}
