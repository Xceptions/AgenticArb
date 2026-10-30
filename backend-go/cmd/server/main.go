package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xceptions/agenticarb/backend-go/internal/api"
)

func main() {
	port := ":8080"

	http.HandleFunc("/api/proposals", api.ProposalHandler)

	fmt.Printf("Autonomous AI Governance Agent active and live on port %s...\n", port)
	fmt.Printf("Listening for inbound UI governance requests on http://localhost%s/api/proposals\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to bind or initialize listener pipelines: %v", err)
	}
}
