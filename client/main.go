package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"gin-example/model"
)

const channel = "mychannel1"

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	order := model.Order{
		Kind: "ETC",
	}
	err := rdb.Publish(context.TODO(), channel, order).Err()
	if err != nil {
		fmt.Println(err)
	}
}
