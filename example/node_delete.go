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

    // Get list of clusters to select from
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
    fmt.Printf("Enter cluster ID to add node to: ")
    fmt.Scanf("%d", &clusterid)

    nodes, err := client.GetNodes(orgid, clusterid)
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    for i := 0; i < len(nodes); i++ {
        fmt.Printf("Nodes(%d): %v\n", nodes[i].PrimaryKey, nodes[i].Name)
    }
    if len(nodes) == 0 {
        fmt.Printf("Sorry, no nodes found\n")
        os.Exit(0)
    }

    var nodeid int
    fmt.Printf("Enter node ID to inspect: ")
    _, err = fmt.Scanf("%d", &nodeid)

    client.DeleteNode(orgid, clusterid, nodeid)
    fmt.Printf("Node should delete shortly\n")

}
