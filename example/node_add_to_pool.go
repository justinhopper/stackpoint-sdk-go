package main

import (
	"fmt"
	spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
	"log"
)

const orgid = 111

func main() {
        // Set up HTTP client with with environment variables for API token and URL
        client, err := spio.NewClientFromEnv()
        if err != nil { log.Fatal(err.Error()) }

        // Get list of configured clusters
        clusters, err := client.GetClusters(orgid)
        if err != nil { log.Fatal(err.Error()) }

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
        if err != nil { log.Fatal(err.Error()) }

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
        if !spio.StringInSlice(nodeSize, mOptions) { log.Fatalf("Invalid option: %s\n", nodeSize) }

	// Get list of nodepools to select from
	nps, err := client.GetNodePools(orgid, clusterid)
	if err != nil { log.Fatal(err.Error()) }

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

	newNode := spio.NodeAdd{Size: nodeSize,
		Count:      nodeCount,
		Role:       "worker",
		NodePoolID: nodepoolid}

	_, err2 := client.AddNodesToNodePool(orgid, clusterid, nodepoolid, newNode)
	if err2 != nil {
                spio.ViewResp()
                log.Fatal(err2)
	}
	fmt.Printf("Node creation sent, building...\n")
}
