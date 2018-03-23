package main

import (
	"fmt"
	spio "github.com/justinhopper/stackpoint-sdk-go/pkg/stackpointio"
	"log"
)

const orgid = 111

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
	for i := 0; i < len(clusters); i++ {
		fmt.Printf("Cluster(%d): %v\n", clusters[i].PrimaryKey, clusters[i].Name)
	}
	if len(clusters) == 0 {
		fmt.Println("Sorry, no clusters defined yet")
		return
	}
	// Get cluster ID from user to delete nodepool from
	var clusterid int
	fmt.Printf("Enter cluster ID to delete nodepool from: ")
	fmt.Scanf("%d", &clusterid)

	// Get list of nodepools to select from
	nps, err := client.GetNodePools(orgid, clusterid)
	if err != nil {
		log.Fatal(err.Error())
	}

	// List nodepools
	for i := 0; i < len(nps); i++ {
		fmt.Printf("Nodepool(%d): %v (node count: %d)\n", nps[i].PrimaryKey, nps[i].Name, nps[i].NodeCount)
	}
	if len(nps) == 0 {
		fmt.Println("Sorry, no nodepools found")
		return
	}
	// Get nodepool ID from user
	var nodepoolid int
	fmt.Printf("Enter nodepool ID to delete: ")
	fmt.Scanf("%d", &nodepoolid)

	_, err2 := client.DeleteNodePool(orgid, clusterid, nodepoolid)
	if err2 != nil {
		spio.ViewResp()
		log.Fatal(err2)
	}
	fmt.Println("NodePool should delete shortly")
}
