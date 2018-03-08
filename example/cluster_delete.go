package main

import (
    "os"
    "fmt"
    spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
)

const mytoken = `98addc5550b98ab74499dd4cd0dc4dd03e0e0d0be82acb0213d1ec7ef2c79457`
const myurl = `https://api-staging.stackpoint.io/`
const orgid = 111

func main() {
    client := spio.NewClient(mytoken, myurl)

    clusters, err := client.GetClusters(orgid)
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    for i := 0; i < len(clusters); i++ {
        fmt.Printf("Clusters(%d): %v\n", clusters[i].PrimaryKey, clusters[i].Name)
    }
    if len(clusters) == 0 {
        fmt.Printf("Sorry, no clusters defined yet\n")
        os.Exit(0)
    }

    var clusterid int
    fmt.Printf("Enter cluster ID to delete: ")
    fmt.Scanf("%d", &clusterid)

    client.DeleteCluster(orgid, clusterid)
    fmt.Printf("Cluster should delete shortly\n")

}
