package DB

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var CTX = context.Background()

var RDB = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})
