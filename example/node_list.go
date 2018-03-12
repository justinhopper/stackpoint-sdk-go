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

        // Get list of nodes configured
        nodes, err := client.GetNodes(orgid, clusterid)
        if err != nil { log.Fatal(err.Error()) }

        // List nodes
        for i := 0; i < len(nodes); i++ {
                fmt.Printf("Nodes(%d): %v\n", nodes[i].PrimaryKey, nodes[i].Name)
        }
        if len(nodes) == 0 {
                fmt.Printf("Sorry, no nodes found\n")
                return
        }
	// Get node ID from user to inspect
	var nodeid int
	fmt.Printf("Enter node ID to inspect: ")
	fmt.Scanf("%d", &nodeid)

	node, err := client.GetNode(orgid, clusterid, nodeid)
	if err != nil {
                spio.ViewResp()
                log.Fatal(err)
	}
	spio.PrettyPrint(node)
}
