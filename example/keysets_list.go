package main

import (
	"fmt"
	spio "github.com/justinhopper/stackpoint-sdk-go/pkg/stackpointio"
	"log"
)

func main() {
	// Set up HTTP client with with environment variables for API token and URL
	client, err := spio.NewClientFromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

        orgid, err := spio.GetIDFromEnv("SPC_ORG_ID")
        if err != nil {
                log.Fatal(err.Error())
        }

	// Gather list of keysets
	keysets, err := client.GetKeysets(orgid)
	if err != nil {
		log.Fatal(err.Error())
	}

	// List keysets configured
	for i := 0; i < len(keysets); i++ {
		fmt.Printf("Keysets(%d): %v\n", keysets[i].PrimaryKey, keysets[i].Name)
	}
	// Get keyset ID to inspect from user
	var keysetid int
	fmt.Printf("Enter keyset ID to delete: ")
	fmt.Scanf("%d", &keysetid)

	keyset, err := client.GetKeyset(orgid, keysetid)
	if err != nil {
		log.Fatal(err.Error())
	}
	spio.PrettyPrint(keyset)
}
