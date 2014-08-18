package main

import (
	"flag"
	"fmt"
	"github.com/fzzy/radix/redis"
	"log"
	"net/http"
	"time"
	"./iputils"
)

// HTTP Handler function
type handler func(response http.ResponseWriter, request *http.Request)

// Returns an array with two HTTP handler functions
// The first item is a handler persists a KV pair with the public IP as the key and the local IP (from the query param "local_ip") as the value
// The second item is a handler that tries to retrieve the local IP for a given public IP.
func makeHandlers(redisAddr string, redisDb int, ttl string) [2]handler {
	var handlers [2]handler
	client, err := redis.DialTimeout("tcp", redisAddr, time.Duration(10)*time.Second)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	r := client.Cmd("select", redisDb)
	if r.Err != nil {
		log.Fatalf("Redis can't select db %d: %v", redisDb, r.Err)
	}

	//announce handler
	handlers[0] = func(w http.ResponseWriter, r *http.Request) {
    	localIp := r.URL.Query().Get("local_ip")

		if iputils.IsIpPrivate(localIp) != true {
			log.Printf("Not a valid local IP: %s", localIp)
			w.WriteHeader(400)
			return
		}

		yourIP := iputils.GetRemoteIpFromRequest(r)
		
		result := client.Cmd("setex", yourIP, ttl, localIp)
		if result.Err != nil {
			log.Printf("Failed saving %s for %s: %v", localIp, yourIP, result.Err)
			w.WriteHeader(500)
		} else {
			fmt.Fprintf(w, "OK")
		}
	}

	//retrieve handler
	handlers[1] = func(w http.ResponseWriter, r *http.Request) {
		var yourIP string
		yourIP = iputils.GetRemoteIpFromRequest(r)
		addr, err := client.Cmd("get", yourIP).Str()
		if err != nil {
		} else {
			fmt.Fprintf(w, "%s", addr)
		}
	}

	return handlers
}

func main() {
	addr := flag.String("addr", ":8080", "The address on which to listen")
	db := flag.Int("db", 0, "The redis database to use (default is 0)")
	ttl := flag.String("ttl", "3600", "The expiry time for stored keys in seconds (default is 3600)")
	flag.Parse()

	handlers := makeHandlers("127.0.01:6379", *db, *ttl)
	http.HandleFunc("/announce", handlers[0])
	http.HandleFunc("/get", handlers[1])

	log.Printf("Starting server on address %s", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalf("Error setting up server: %v", err)
	}
}
