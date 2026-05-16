package main

import (
	"fmt"
	"log"
	"vercel-backend/pkg/portal"
)

func main() {
	
	username := "bn_dv_tcdvquevo"
	password := "Tinh@2027"
	
	// We don't use Redis for local test to simplify, unless needed
	client := portal.NewPortalClient(username, password, nil)
	
	phone := "0388634123"
	fmt.Printf("Searching for %s using Go code...\n", phone)
	
	patients, err := client.LookupPatients(phone)
	if err != nil {
		log.Fatalf("Lookup failed: %v", err)
	}
	
	fmt.Printf("Found %d patients\n", len(patients))
	for _, p := range patients {
		fmt.Printf("- ID: %s, Name: %s, DOB: %s\n", p.ID, p.Name, p.DOB)
	}
}
