package main

import (
	"fmt"
	spio "github.com/justinhopper/stackpoint-sdk-go/pkg/stackpointio"
	"log"
)

const orgid = 111
const nodepoolName = "Test Nodepool"

func main() {
	// Set up HTTP client with with environment variables for API token and URL
	client, err := spio.NewClientFromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get list of configured clusters
	clusters, err := client.GetClusters(orgid)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Print list of clusters, saving map of providers for later use
	providers := make(map[int]string)
	for i := 0; i < len(clusters); i++ {
		fmt.Printf("Cluster(%d): %v\n", clusters[i].PrimaryKey, clusters[i].Name)
		providers[clusters[i].PrimaryKey] = clusters[i].Provider
	}
	if len(clusters) == 0 {
		fmt.Println("Sorry, no clusters defined yet")
		return
	}
	// Get cluster ID from user to add node to
	var clusterid int
	fmt.Printf("Enter cluster ID to add node to: ")
	fmt.Scanf("%d", &clusterid)

	// Get number of worker nodes to add to nodepool
	var nodeCount int
	fmt.Printf("Enter number of worker nodes to add into pool: ")
	fmt.Scanf("%v", &nodeCount)

	if nodeCount < 1 {
		log.Fatal("You must add at least one node to the new nodepool")
	}

	newNodepool := spio.NodePool{Name: nodepoolName,
		ClusterID: clusterid,
		NodeCount: nodeCount,
		Platform:  "coreos",
		State:     "draft",
		Role:      "worker"}

	// Get list of machine types for provider
	mOptions, err := client.GetMachSpecs(providers[clusterid])
	if err != nil {
		log.Fatal(err.Error())
	}

	// List machine types
	fmt.Printf("Node size options for provider %s:\n", providers[clusterid])
	for _, opt := range mOptions {
		fmt.Println(opt)
	}
	// Get node size selection from user
	var nodeSize string
	fmt.Printf("Enter node size: ")
	fmt.Scanf("%s", &nodeSize)

	newNodepool.Size = nodeSize

	// Create new nodepool
	pool, err := client.CreateNodePool(orgid, clusterid, newNodepool)
	if err != nil {
		spio.ViewResp()
		log.Fatal(err)
	}
	fmt.Printf("NodePool creation sent (ID: %d), building...\n", pool.PrimaryKey)
}
