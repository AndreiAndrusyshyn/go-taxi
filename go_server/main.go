package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

type single_response struct {
	Request string `json:"request"`
}

func main() {
	flag.Parse()

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	switch string(ctx.Path()) {
	case "/request":
		val2, err := client.HKeys("order").Result()
		if err != nil {
			panic(err)
		}
		var pair = val2[rand.Intn(len(val2))]
		client.HIncrBy("order", pair, 1 )

		var resp = &single_response{Request:pair}
		a, _ := json.Marshal(resp)

		fmt.Fprintf(ctx, string(a))
	case "/admin/requests":
		val2 := client.HGetAll("order").Val()

		m := make(map[string]string)

		for i, a := range val2{
			if a != "0" {
				m[i] = a
			}
		}
		q, _ := json.Marshal(m)
		fmt.Fprintf(ctx, string(q))
	default:
		fmt.Fprintf(ctx, "no route")
	}

	ctx.Response.Header.Set( "Content-Type", "application/json")
}
