package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	interval := 5 * time.Second
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	log.Println("starting wilderness permit monitor...")

	for t := range ticker.C {

		fmt.Println("-------------------")
		log.Printf("tick at %v", t)

		if err := checkPermits(); err != nil {
			log.Printf("Error checking permits: %v", err)
		}

	}

}

func checkPermits() error {
	url := "https://www.recreation.gov/api/permitinyo/233260/availability?start_date=2026-08-01&end_date=2026-09-30"

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status: %d", resp.StatusCode)
	}

	var data RootResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}

	countFount := 0

	for id, node := range data.Payload {
		fmt.Printf("checking trailhead id: %s\n", id)

		for date, count := range node.DateAvailability {
			if count > 0 {
				fmt.Printf("  -> FOUND! Date: %s | Slots: %d\n", date, count)
				countFount++
			}
		}
	}

	if countFount == 0 {
		log.Println("no permits found in this batch")
	}

	return nil
}

type AvailabilityNode struct {
	DateAvailability map[string]int `json:"date_availability"`
}

type RootResponse struct {
	Payload map[string]AvailabilityNode `json:"payload"`
}
