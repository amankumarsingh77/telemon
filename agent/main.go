package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	serverAddr := flag.String("server", "localhost:8080", "Collector server address")
	pollInterval := flag.Int("interval", 5, "Poll interval in seconds")
	flag.Parse()
	agent := NewTelemetry(*serverAddr, time.Duration(*pollInterval)*time.Second)
	if err := agent.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	log.Printf("Starting telemetry agent, sending to %s every %v",
		*serverAddr, time.Duration(*pollInterval)*time.Second)

	agent.Run()
}
