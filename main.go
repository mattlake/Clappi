package main

import (
	"fmt"
	"log"

	"github.com/mattlake/clappi/OpenApi"
)

func main() {
	fmt.Println("Loading specs")
	yamlPath := "OpenApi/Specs/simple-api-yaml.yaml"
	jsonPath := "OpenApi/Specs/simple-api-json.json"

	fmt.Println("Attempting to load yaml")
	if err := OpenApi.PrintTitle(yamlPath); err != nil {
		log.Fatalf("Failed to load yaml: %v", err)
	}

	fmt.Println("Attempting to load json")
	if err := OpenApi.PrintTitle(jsonPath); err != nil {
		log.Fatalf("Failed to load json: %v", err)
	}
}
