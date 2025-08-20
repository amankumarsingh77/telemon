package main

import (
	"flag"
	"github.com/amankumarsinghy77/telemon/server/storage"
	"log"
)

func main() {
	udpAddr := flag.String("udp", ":8080", "UDP listen address")
	flag.Parse()

	inmemStorage := storage.NewInMemoryStorage()
	server := NewCollectorServer(*udpAddr, inmemStorage)

	log.Printf("Starting collector server on %s", *udpAddr)

	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
