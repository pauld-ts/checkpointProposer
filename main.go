package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type apiResponseBody struct {
	Height string      `json:"height"`
	Result []Validator `json:"result"`
}

type Validator struct {
	ID int `json:"ID"`
}

type checkPointAPIResponseBody struct {
	Height string     `json:"height"`
	Result checkpoint `json:"result"`
}

type checkpoint struct {
	ID int `json:"id"`
}

func main() {
	heimdallURL := os.Getenv("HEIMDALLURL")

	if heimdallURL == "" {
		heimdallURL = "https://polygon-heimdall-rest.publicnode.com"
	}

	proposerURL := heimdallURL + "/staking/proposer/100"
	latestCheckpointURL := heimdallURL + "/checkpoints/latest"

	// How many minutes in between runs Default is 1 minute
	interval := os.Getenv("INTERVAL")

	if interval == "" {
		interval = "1"
	}

	i, err := strconv.Atoi(interval)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Duration(i) * time.Minute)

	if os.Getenv("LOGFILE") == "true" {
		f, err := os.OpenFile("proposer.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		wrt := io.MultiWriter(os.Stdout, f)
		log.SetOutput(wrt)
	} else {
		log.SetOutput(os.Stdout)
	}

	for {
		select {
		case <-ticker.C:
			var v apiResponseBody
			resp, err := http.Get(proposerURL)
			if err != nil {
				log.Fatal(err)
			}

			err = json.NewDecoder(resp.Body).Decode(&v)
			if err != nil {
				log.Fatal(err)
			}
			resp.Body.Close()

			var c checkPointAPIResponseBody
			resp, err = http.Get(latestCheckpointURL)
			if err != nil {
				log.Fatal(err)
			}

			err = json.NewDecoder(resp.Body).Decode(&c)
			if err != nil {
				log.Fatal(err)
			}
			resp.Body.Close()

			twinStakePosition := 0
			for i, r := range v.Result {
				if r.ID == 148 {
					twinStakePosition = i
					break
				}
			}
			log.Printf("Proposer ID: %v  Twinstake position: %v  Current Checkpoint: %v \n",
				v.Result[0].ID,
				twinStakePosition,
				c.Result.ID,
			)
		}
	}

}
