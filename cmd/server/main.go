package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

func main() {
	src := rand.NewSource(time.Now().Unix())
	rnd := rand.New(src)
	rand.Seed(time.Now().Unix())

	runtime.GOMAXPROCS(1)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep((100 + (time.Duration)(rnd.Intn(300))) * time.Millisecond)
		w.Write(([]byte)("Hello world"))
	})

	fmt.Println("Server listening on port 8080")
	fmt.Println("HTTP Error :", http.ListenAndServe(":8080", nil))
}
