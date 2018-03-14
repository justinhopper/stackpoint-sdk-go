package main

import (
	"fmt"
	spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
	"log"
)

const orgid = 111
const provider = "do"
const clusterName = "Test DigitalOcean Cluster"

func main() {
	// Set up HTTP client with with environment variables for API token and URL
	client, err := spio.NewClientFromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	sshKeysetid, err := spio.GetIDFromEnv("SPC_SSH_KEYSET")
	if err != nil {
		log.Fatal(err.Error())
	}

	doKeysetid, err := spio.GetIDFromEnv("SPC_DO_KEYSET")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get list of machine types for provider
	mOptions, err := client.GetMachSpecs(provider)
	if err != nil {
		log.Fatal(err.Error())
	}

	// List machine types
	fmt.Printf("Node size options for provider %s:\n", provider)
	for _, opt := range mOptions {
		fmt.Println(opt)
	}
	// Get node size selection from user
	var nodeSize string
	fmt.Printf("Enter node size: ")
	fmt.Scanf("%s", &nodeSize)

	// Validate machine type selection
	if !spio.StringInSlice(nodeSize, mOptions) {
		log.Fatalf("Invalid option: %s\n", nodeSize)
	}

	newSolution := spio.Solution{Solution: "helm_tiller"}
	newCluster := spio.Cluster{Name: clusterName,
		Provider:          provider,
		ProviderKey:       doKeysetid,
		MasterCount:       1,
		MasterSize:        nodeSize,
		WorkerCount:       2,
		WorkerSize:        nodeSize,
		Region:            "nyc1",
		State:             "draft",
		KubernetesVersion: "v1.8.3",
		RbacEnabled:       true,
		DashboardEnabled:  true,
		EtcdType:          "self_hosted",
		Platform:          "coreos",
		Channel:           "stable",
		SSHKeySet:         sshKeysetid,
		Solutions:         []spio.Solution{newSolution}}

	_, err2 := client.CreateCluster(orgid, newCluster)
	if err2 != nil {
		spio.ViewResp()
		log.Fatal(err2)
	}
	fmt.Printf("Cluster created, building...\n")
}
