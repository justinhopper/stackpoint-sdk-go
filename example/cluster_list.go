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

	// Get cluster ID from user to inspect
	var clusterid int
	fmt.Printf("Enter cluster ID to inspect: ")
	fmt.Scanf("%d", &clusterid)

	cluster, err := client.GetCluster(orgid, clusterid)
	if err != nil {
		log.Fatal(err.Error())
	}
	spio.PrettyPrint(cluster)
}
