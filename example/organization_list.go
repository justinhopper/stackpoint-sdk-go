package main

import (
	"fmt"
	spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
	"log"
)

func main() {
	// Set up HTTP client with with environment variables for API token and URL
	client, err := spio.NewClientFromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get list of configured organizations
	orgs, err := client.GetOrganizations()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Print list of organizations
	for i := 0; i < len(orgs); i++ {
		fmt.Printf("Org(%d): %v\n", orgs[i].PrimaryKey, orgs[i].Name)
	}
	if len(orgs) == 0 {
		fmt.Println("Sorry, no organizations defined yet")
		return
	}
	// Get organization ID from user to inspect
	var orgid int
	fmt.Printf("Enter org ID to inspect: ")
	fmt.Scanf("%d", &orgid)

	org, err := client.GetOrganization(orgid)
	if err != nil {
		spio.ViewResp()
		log.Fatal(err)
	}
	spio.PrettyPrint(org)
}
