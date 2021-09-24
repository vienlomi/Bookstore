package redis

import (
	"fmt"
	"project_v3/config"
	"strconv"

	"github.com/go-redis/redis"
)

var Client *redis.Client

//create client to connect redis
func NewClient() {
	dbRedis, _ := strconv.Atoi(config.DataConfig["DB_redis"])
	Client = redis.NewClient(&redis.Options{
		Addr:     config.DataConfig["address_redis"],
		Password: config.DataConfig["password_redis"],
		DB:       dbRedis,
	})

	pong, err := Client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

//func close Client redis
func CloseClient() {
	Client.Close()
}
