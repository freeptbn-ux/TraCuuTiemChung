package main

import (
	"fmt"
	"log"
	"vercel-backend/pkg/portal"
)

func main() {
	username := "bn_dv_tcdvquevo"
	password := "Tinh@2027"
	
	client := portal.NewPortalClient(username, password, nil)
	
	phone := "0999999999"
	fmt.Printf("Searching for %s (DIRECTLY)...\n", phone)
	
	patients, err := client.LookupPatients(phone)
	if err != nil {
		log.Fatalf("Lookup failed: %v", err)
	}
	
	fmt.Printf("Found %d patients\n", len(patients))
	for _, p := range patients {
		fmt.Printf("- %s (%s)\n", p.Name, p.ID)
	}
}
