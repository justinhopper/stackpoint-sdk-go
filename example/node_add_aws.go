package main

import (
	"fmt"
	spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
	"log"
)

const orgid = 111
const awsZone = "us-east-2a"
const awsSubnetID = "subnet-a4b630cc"
const awsSubnetCidr = "172.31.0.0/24"

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

	// Validate machine type selection
	if !spio.StringInSlice(nodeSize, mOptions) {
		log.Fatalf("Invalid option: %s\n", nodeSize)
	}

	// Get type of node from user
	var nodeRole string
	fmt.Printf("Enter node role (worker or master) to use: ")
	fmt.Scanf("%v", &nodeRole)

	// Set up new node based on role
	newNode := spio.NodeAdd{
		Zone: awsZone,
                ProviderSubnetID:   awsSubnetID,
                ProviderSubnetCidr: awsSubnetCidr }
	if nodeRole == "master" {
		newNode.Count = 1
		newNode.Role = "master"
		newNode.Size = nodeSize
	} else if nodeRole == "worker" {
		// Get number of worker nodes to add
		var nodeCount int
		fmt.Printf("Enter number of worker nodes to add: ")
		fmt.Scanf("%v", &nodeCount)

		newNode.Count = nodeCount
		newNode.Role = "master"
		newNode.Size = nodeSize
	} else {
		log.Fatalf("Invalid node role: %s\n", nodeRole)
	}
	// Add new node
	_, err2 := client.AddNodes(orgid, clusterid, newNode)
	if err2 != nil {
		spio.ViewResp()
		log.Fatal(err2)
	}
	fmt.Println("Node creation sent, building...")
}
