package main

import (
	"encoding/json"
	"fmt"
	spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
	"os"
)

// PrettyPrint to break down objects
func PrettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}

const orgid = 111

func main() {
	// Set up HTTP client with API token and URL
	token := os.Getenv("SPC_API_TOKEN")
	endpoint := os.Getenv("SPC_BASE_API_URL")
	client := spio.NewClient(token, endpoint)

	clusters, err := client.GetClusters(orgid)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	for i := 0; i < len(clusters); i++ {
		fmt.Printf("Cluster(%d): %v\n", clusters[i].PrimaryKey, clusters[i].Name)
	}
	if len(clusters) == 0 {
		fmt.Printf("Sorry, no clusters defined yet\n")
		os.Exit(0)
	}

	var clusterid int
	fmt.Printf("Enter cluster ID to list nodes from: ")
	fmt.Scanf("%d", &clusterid)

	nodes, err := client.GetNodes(orgid, clusterid)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	for i := 0; i < len(nodes); i++ {
		fmt.Printf("Node(%d): %v (%v)\n", nodes[i].PrimaryKey, nodes[i].Name, nodes[i].Role)
	}
	if len(nodes) == 0 {
		fmt.Printf("Sorry, no nodes found\n")
		os.Exit(0)
	}

	var nodeid int
	fmt.Printf("Enter node ID to inspect: ")
	_, err = fmt.Scanf("%d", &nodeid)

	node, err := client.GetNode(orgid, clusterid, nodeid)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	PrettyPrint(node)
}
