package main

import (
	"github.com/go-redis/redis/v7"
	"math/rand"
	"time"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	client.FlushAll()
	rand.Seed(time.Now().Unix())

	var letters = [26]string{"a", "b", "c", "d" ,"e", "f", "g", "h", "i", "j", "k", "l",
		"m","n","o", "p", "q", "r","s", "t", "u", "v", "w", "x", "y", "z"}

	for i := 0; i < 50; i++ {
		var two_random_letters = letters[rand.Intn(len(letters))] + letters[rand.Intn(len(letters))]
		var res , _ = client.HSetNX("order", two_random_letters, 0).Result()
		if !res { i-- }
	}

	for {
		val2, err := client.HKeys("order").Result()
		if err != nil {
			panic(err)
		}
		client.HDel("order", val2[rand.Intn(len(val2))])

		for {
			var two_random_letters = letters[rand.Intn(len(letters))] + letters[rand.Intn(len(letters))]
			var res , _ = client.HSetNX("order", two_random_letters, 0).Result()
			if res { break }
		}

		time.Sleep(200 * time.Millisecond)
	}
}


