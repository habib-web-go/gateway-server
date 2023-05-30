package redis

import (
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	RateLimit         = 100
	RateLimitDuration = time.Hour * 24
)

var (
	redisClient *redis.Client
)

func RunRedis() {
	redisPort := os.Getenv("REDIS_PORT")
	redisUrl := os.Getenv("REDIS_URL")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisUrl + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})

}

func GetRateLimit(ip string) (int, error) {
	err := redisClient.SetNX(ip, RateLimit, RateLimitDuration).Err()
	if err != nil {
		return 0, nil
	}
	rate, err := redisClient.Get(ip).Result()
	if err != nil {
		return 0, nil
	}
	return strconv.Atoi(rate)
}

func DecreaseRateLimit(ip string) error {
	err := redisClient.DecrBy(ip, 1).Err()
	return err
}
