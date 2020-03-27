package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"log"
	"math/rand"
	"runtime"
)

var (

	client = redis.NewClient(&redis.Options{
		Addr:     "unix://localhost:6379",
		Password: "",
		DB:       0,
	})
)

type single_response struct {
	Request string `json:"request"`
}

func main() {
	runtime.GOMAXPROCS(1)
	flag.Parse()

	ln, err := reuseport.Listen("tcp4", "localhost:8080")
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}
	if err = fasthttp.Serve(ln, requestHandler); err != nil {
		log.Fatalf("error in fasthttp Server: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {

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
		m := make(map[string]string)
		m["err"] = "no route"
		q, _ := json.Marshal(m)
		fmt.Fprintf(ctx, string(q))
	}

	ctx.Response.Header.Set( "Content-Type", "application/json")
}
