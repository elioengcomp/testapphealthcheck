package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var maxTicks = 10
var ticks = 0

func healthcheck(w http.ResponseWriter, req *http.Request) {
	log.Printf("Healthcheck requested. Ticks: %d/%d", ticks, maxTicks)
	fmt.Fprint(w, "OK")
}

func main() {
	log.Print("Starting testappheathcheck...")

	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Print("Increasing ticks...")
				ticks++
				if ticks >= maxTicks {
					panic("App has run out of ticks")
				}
			}
		}
	}()

	http.HandleFunc("/healthcheck", healthcheck)
	http.ListenAndServe(":80", nil)
}
