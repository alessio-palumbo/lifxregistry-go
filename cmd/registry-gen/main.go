package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/alessio-palumbo/lifxregistry-go/cmd/registry-gen/decode"
	"github.com/alessio-palumbo/lifxregistry-go/cmd/registry-gen/generate"
)

const (
	registryGenerateDir = "gen/registry"
)

//go:embed src/products.json
var productsJSON []byte

//go:embed src/products_commit.txt
var productsCommit []byte

func main() {
	sourceCommit := strings.TrimSpace(string(productsCommit))
	if sourceCommit == "" {
		log.Fatal("Source commit not found")
	}

	productsSpec, err := decode.DecodeProductsRegistry(productsJSON)
	if err != nil {
		log.Fatalf("Failed to parse embedded products.json: %v", err)
	}

	if err := generate.GenerateProductsRegistry(productsSpec, registryGenerateDir, sourceCommit); err != nil {
		log.Fatalf("Failed to generate Products registry: %v", err)
	}

	fmt.Println("Code generation completed successfully.")
}
