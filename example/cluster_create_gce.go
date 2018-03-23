package main

import (
	"fmt"
	spio "github.com/justinhopper/stackpoint-sdk-go/pkg/stackpointio"
	"log"
)

const (
	orgid       = 111
	provider    = "gce"
	clusterName = "Test GCE Cluster"
	region      = "us-west1-a"
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

	gceKeysetid, err := spio.GetIDFromEnv("SPC_GCE_KEYSET")
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
		ProviderKey:       gceKeysetid,
		MasterCount:       1,
		MasterSize:        nodeSize,
		WorkerCount:       2,
		WorkerSize:        nodeSize,
		Region:            "us-west1-a",
		State:             "draft",
		KubernetesVersion: "v1.8.7",
		RbacEnabled:       true,
		DashboardEnabled:  true,
		EtcdType:          "self_hosted",
		Platform:          "coreos",
		Channel:           "stable",
		SSHKeySet:         sshKeysetid,
		Solutions:         []spio.Solution{newSolution}}

	cluster, err := client.CreateCluster(orgid, newCluster)
	if err != nil {
		spio.ViewResp()
		log.Fatal(err)
	}
	fmt.Printf("Cluster created (ID: %d) (instance name: %s), building...\n", cluster.PrimaryKey, cluster.InstanceID)
}
