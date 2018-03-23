package main

import (
	"fmt"
	spio "github.com/justinhopper/stackpoint-sdk-go/pkg/stackpointio"
	"log"
)

const (
	orgid      = 111
	keysetName = "Test AWS Keyset"
	provider   = "aws"
)

func main() {
	// Set up HTTP client with with environment variables for API token and URL
	client, err := spio.NewClientFromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Gather access key for AWS
	var awsAccessKey string
	fmt.Printf("Enter AWS Access Key ID: ")
	fmt.Scanf("%s", &awsAccessKey)

	// Gather secret key for AWS
	var awsSecretKey string
	fmt.Printf("Enter AWS Secret Access Key: ")
	fmt.Scanf("%s", &awsSecretKey)

	pubKey := spio.Key{Type: "pub",
		Value: awsAccessKey}
	pvtKey := spio.Key{Type: "pvt",
		Value: awsSecretKey}
	newKeyset := spio.Keyset{Name: keysetName,
		Category:   "provider",
		Entity:     provider,
		Workspaces: []int{},
		Keys:       []spio.Key{pubKey, pvtKey}}

	_, err2 := client.CreateKeyset(orgid, newKeyset)
	if err2 != nil {
		spio.ViewResp()
		log.Fatal(err2)
	}
	fmt.Printf("CreateKeyset created\n")
}
