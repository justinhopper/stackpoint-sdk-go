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

	// Fetch list of configured clusters
	clusters, err := client.GetClusters(orgid)
	if err != nil {
		log.Fatal(err.Error())
	}

	// List clusters
	for i := 0; i < len(clusters); i++ {
		fmt.Printf("Clusters(%d): %v\n", clusters[i].PrimaryKey, clusters[i].Name)
	}
	if len(clusters) == 0 {
		fmt.Println("Sorry, no clusters defined yet")
		return
	}

	// Get cluster ID to delete from user
	var clusterid int
	fmt.Printf("Enter cluster ID to delete: ")
	fmt.Scanf("%d", &clusterid)

	_, err2 := client.DeleteCluster(orgid, clusterid)
	if err2 != nil {
		spio.ViewResp()
		log.Fatal(err2.Error())
	}

	fmt.Printf("Cluster should delete shortly\n")
}
