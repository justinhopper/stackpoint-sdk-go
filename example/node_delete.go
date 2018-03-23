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

	// Print list of clusters
	for i := 0; i < len(clusters); i++ {
		fmt.Printf("Cluster(%d): %v\n", clusters[i].PrimaryKey, clusters[i].Name)
	}
	if len(clusters) == 0 {
		fmt.Println("Sorry, no clusters defined yet")
		return
	}
	// Get cluster ID from user to delete node from
	var clusterid int
	fmt.Printf("Enter cluster ID to delete node from: ")
	fmt.Scanf("%d", &clusterid)

	// Get list of nodes configured
	nodes, err := client.GetNodes(orgid, clusterid)
	if err != nil {
		log.Fatal(err.Error())
	}

	// List nodes
	for i := 0; i < len(nodes); i++ {
		fmt.Printf("Nodes(%d): %s node is %s\n", nodes[i].PrimaryKey, nodes[i].Role, nodes[i].State)
	}
	if len(nodes) == 0 {
		fmt.Printf("Sorry, no nodes found\n")
		return
	}
	// Get node ID to delete from user
	var nodeid int
	fmt.Printf("Enter node ID to delete: ")
	fmt.Scanf("%d", &nodeid)

	_, err2 := client.DeleteNode(orgid, clusterid, nodeid)
	if err2 != nil {
		spio.ViewResp()
		log.Fatal(err2)
	}
	fmt.Println("Node should delete shortly")
}
