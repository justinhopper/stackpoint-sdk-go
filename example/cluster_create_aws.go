package main

import (
	"fmt"
	spio "github.com/justinhopper/stackpoint-sdk-go/pkg/stackpointio"
	"log"
)

const (
	orgid          = 111
	provider       = "aws"
	clusterName    = "Test AWS Cluster"
	awsRegion      = "us-west-2"
	awsZone        = "us-west-2a"
	awsNetworkID   = "__new__"
	awsNetworkCidr = "172.23.0.0/16"
	awsSubnetID    = "__new__"
	awsSubnetCidr  = "172.23.5.0/24"
)

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

	awsKeysetid, err := spio.GetIDFromEnv("SPC_AWS_KEYSET")
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
		Provider:           provider,
		ProviderKey:        awsKeysetid,
		MasterCount:        1,
		MasterSize:         nodeSize,
		WorkerCount:        2,
		WorkerSize:         nodeSize,
		Region:             awsRegion,
		Zone:               awsZone,
		State:              "draft",
		ProviderNetworkID:  awsNetworkID,
		ProviderNetworkCdr: awsNetworkCidr,
		ProviderSubnetID:   awsSubnetID,
		ProviderSubnetCidr: awsSubnetCidr,
		KubernetesVersion:  "v1.8.7",
		RbacEnabled:        true,
		DashboardEnabled:   true,
		EtcdType:           "self_hosted",
		Platform:           "coreos",
		Channel:            "stable",
		SSHKeySet:          sshKeysetid,
		Solutions:          []spio.Solution{newSolution}}

	_, err2 := client.CreateCluster(orgid, newCluster)
	if err2 != nil {
		spio.ViewResp()
		log.Fatal(err2)
	}
	fmt.Printf("Cluster created, building...\n")

}
