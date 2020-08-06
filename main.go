package main

/*
$ curl http://localhost:9999/_distcache/scores/Tom
630
$ curl http://localhost:9999/_distcache/scores/kkk
kkk not exist
*/

import (
	"fmt"
	"distcache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	distcache.NewGroup("scores", 2<<10, distcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := distcache.NewHTTPPool(addr)
	log.Println("distcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
