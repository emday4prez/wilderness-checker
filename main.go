package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ProbeResponse struct {
	Payload map[string]interface{} `json:"payload"`
}

func main() {

	fmt.Println("Starting Scan for Bubbs Creek...")
	fmt.Println("------------------------------------------------")

	// brute force
	// SEKI trailhead ids
	targets := []int{}
	for i := 1; i <= 150; i++ {
		targets = append(targets, i)
	} // 1-150
	for i := 440; i <= 480; i++ {
		targets = append(targets, i)
	} // 440-480

	var wg sync.WaitGroup
	client := &http.Client{Timeout: 3 * time.Second}

	for _, id := range targets {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			checkID(client, id)
		}(id)

		if id%10 == 0 {
			time.Sleep(50 * time.Millisecond)
		}
	}

	wg.Wait()
	fmt.Println("------------------------------------------------")
	fmt.Println("Scan Complete.")
}

func checkID(client *http.Client, id int) {

	url := fmt.Sprintf("https://www.recreation.gov/api/permits/%d/availability/month?start_date=2026-08-01&commercial_acct=false", id)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var data ProbeResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {

			fmt.Printf("\033[32m[HIT] Found valid data at ID: %d\033[0m\n", id)
		}
	}
}
