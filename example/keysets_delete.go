package main

import (
	"fmt"
	spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
	"os"
)

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
		fmt.Printf("Keysets(%d): %v\n", keysets[i].PrimaryKey, keysets[i].Name)
	}

	var keysetid int
	fmt.Printf("Enter keyset ID to delete: ")
	fmt.Scanf("%d", &keysetid)

	client.DeleteKeyset(orgid, keysetid)
	fmt.Printf("Keyset should delete shortly\n")
}
