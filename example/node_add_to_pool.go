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
	fmt.Printf("Enter nodepool ID to add node to: ")
	fmt.Scanf("%d", &nodepoolid)

	// Get number of nodes to add from user
	var nodeCount int
	fmt.Printf("Enter number of nodes to add: ")
	fmt.Scanf("%d", &nodeCount)

	// Use NodeAddToPool struct for this endpoint
	newNode := spio.NodeAddToPool{Count: nodeCount,
		Role:       "worker",
		NodePoolID: nodepoolid}

	nodes, err := client.AddNodesToNodePool(orgid, clusterid, nodepoolid, newNode)
	if err != nil {
		spio.ViewResp()
		log.Fatal(err)
	}
	fmt.Println("Node creation sent,")
	for _, n := range nodes {
		fmt.Printf("ID: %d, InstanceID: %s", n.PrimaryKey, n.InstanceID)
	}
	fmt.Println("...building...")
}
