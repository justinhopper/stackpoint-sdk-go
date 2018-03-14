package main

import (
	"fmt"
	spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
	"log"
)

const orgid = 111
const keysetName = "Test DO Keyset"
const provider = "do"

func main() {
	// Set up HTTP client with with environment variables for API token and URL
	client, err := spio.NewClientFromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Gather access token for DO
	var doToken string
	fmt.Printf("Enter DigitalOcean Token: ")
	fmt.Scanf("%s", &doToken)

	newKey := spio.Key{Type: "token",
		Value: doToken}
	newKeyset := spio.Keyset{Name: keysetName,
		Category:   "provider",
		Entity:     provider,
		Workspaces: []int{},
		Keys:       []spio.Key{newKey}}

	_, err2 := client.CreateKeyset(orgid, newKeyset)
	if err2 != nil {
		spio.ViewResp()
		log.Fatal(err2)
	}
	fmt.Printf("CreateKeyset created\n")
}