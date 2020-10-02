package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strconv"
	"time"
)

var maxTicks = 12
var ticks = 0

func healthcheck(w http.ResponseWriter, req *http.Request) {
	log.Printf("Healthcheck requested. Ticks: %d/%d", ticks, maxTicks)
	fmt.Fprint(w, "OK")
}

func main() {
	logOutput := os.Args[2]
	if logOutput == "syslog" {
		logwriter, e := syslog.New(syslog.LOG_NOTICE, "testappheathcheck")
		if e != nil {
			panic(fmt.Sprint("Syslog not available: %w", e))
		} else {
			log.SetOutput(logwriter)
		}
	}

	log.Print("Starting testappheathcheck...")

	var err error
	maxTicks, err = strconv.Atoi(os.Args[1])
	log.Printf("Max ticks: %d", maxTicks)

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Print("Tick...")
				ticks++
				resp, err := http.Get("https://search.gocenter.io/api/ui/search?name_fragment=elio")
				if err != nil {
					log.Printf("Failed to send request: Response: %d Error: %s", resp.StatusCode, err)
				}
				defer resp.Body.Close()
				if ticks >= maxTicks {
					log.Print("App has run out of ticks. Shuting down...")
					panic("App has run out of ticks")
				}
			}
		}
	}()

	http.HandleFunc("/healthcheck", healthcheck)
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}
