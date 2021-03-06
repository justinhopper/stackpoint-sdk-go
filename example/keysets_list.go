package main

import (
	"encoding/json"
	"fmt"
	spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
	"os"
)

// PrettyPrint to break down objects
func PrettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}

const orgid = 111

func main() {
	// Set up HTTP client with API token and URL
	token := os.Getenv("SPC_API_TOKEN")
	endpoint := os.Getenv("SPC_BASE_API_URL")
	client := spio.NewClient(token, endpoint)

	keysets, err := client.GetKeysets(orgid)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	for i := 0; i < len(keysets); i++ {
		fmt.Printf("Keyset(%d): %v\n", keysets[i].PrimaryKey, keysets[i].Name)
	}

	var keysetid int
	fmt.Printf("Enter keyset ID to inspect: ")
	fmt.Scanf("%d", &keysetid)

	keyset, err := client.GetKeyset(orgid, keysetid)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	PrettyPrint(keyset)
}
