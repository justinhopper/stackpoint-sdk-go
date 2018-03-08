package main

import (
    "os"
    "fmt"
    "encoding/json"
    spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
)

// PrettyPrint to break down objects
func PrettyPrint(v interface{}) {
      b, _ := json.MarshalIndent(v, "", "  ")
      println(string(b))
}

const mytoken = `98addc5550b98ab74499dd4cd0dc4dd03e0e0d0be82acb0213d1ec7ef2c79457`
const myurl = `https://api-staging.stackpoint.io/`
const orgid = 111

func main() {
    client := spio.NewClient(mytoken, myurl)

    clusters, err := client.GetClusters(orgid)
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    for i := 0; i < len(clusters); i++ {
        fmt.Printf("Cluster(%d): %v\n", clusters[i].PrimaryKey, clusters[i].Name)
    }
    if len(clusters) == 0 {
        fmt.Printf("Sorry, no clusters defined yet\n")
        os.Exit(0)
    }

    var clusterid int
    fmt.Printf("Enter cluster ID to list nodepools from: ")
    _, err = fmt.Scanf("%d", &clusterid)

    nps, err := client.GetNodePools(orgid, clusterid)
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    for i := 0; i < len(nps); i++ {
        fmt.Printf("Nodepool(%d): %v (node count: %d)\n", nps[i].PrimaryKey, nps[i].Name, nps[i].NodeCount)
    }
    if len(nps) == 0 {
        fmt.Printf("Sorry, no nodepools found\n")
        os.Exit(0)
    }

    var nodepoolid int
    fmt.Printf("Enter nodepool ID to inspect: ")
    _, err = fmt.Scanf("%d", &nodepoolid)

    nodepool, err := client.GetNodePool(orgid, clusterid, nodepoolid)
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    PrettyPrint(nodepool)
}
