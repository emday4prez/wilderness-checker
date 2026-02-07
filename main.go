package main

import (
	"fmt"
	"log"
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
	fmt.Println("   Checking Recreation.gov (Simulated)...")

	return nil
}
