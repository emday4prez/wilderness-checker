package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type CampgroundResponse struct {
	Campsites map[string]CampsiteNode `json:"campsites"`
}

type CampsiteNode struct {
	Loop           string            `json:"loop"`
	CampsiteType   string            `json:"campsite_type"`
	Availabilities map[string]string `json:"availabilities"` // Key is Date, Value is Status
}

func main() {

	campgroundID := "232445"

	baseURL := fmt.Sprintf("https://www.recreation.gov/api/camps/availability/campground/%s/month", campgroundID)

	startDate := "2026-06-01T00:00:00.000Z"

	params := url.Values{}
	params.Add("start_date", startDate)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	fmt.Printf("Checking: %s\n", fullURL)

	req, _ := http.NewRequest("GET", fullURL, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Referer", "https://www.recreation.gov/")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Failed! Status: %d\nBody: %s", resp.StatusCode, string(body))
	}

	var data CampgroundResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal("JSON Decode Error:", err)
	}

	fmt.Println("------------------------------------------------")
	fmt.Printf("Scanning %d campsites at Watchman...\n", len(data.Campsites))

	found := 0
	for id, site := range data.Campsites {
		for date, status := range site.Availabilities {
			if status == "Available" {
				// Clean up the date string for printing
				cleanDate := date[:10]
				fmt.Printf("[OPEN] Site ID %s (%s) is open on %s\n", id, site.Loop, cleanDate)
				found++
			}
		}
	}

	if found == 0 {
		fmt.Println("No open spots found (but the API worked!).")
	} else {
		fmt.Printf("Total open spots found: %d\n", found)
	}
}
