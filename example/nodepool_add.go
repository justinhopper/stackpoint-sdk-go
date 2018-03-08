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

    // Get cluster ID selection from user
    var clusterid int
    fmt.Printf("Enter cluster ID to add nodepool to: ")
    fmt.Scanf("%d", &clusterid)

    // Get cluster provider from selection
    var provider_name string
    for i := 0; i < len(clusters); i++ {
        if clusters[i].PrimaryKey == clusterid {
            provider_name = clusters[i].Provider
            break;
        }
    }

    // Get number of worker nodes to add
    var node_count int
    fmt.Printf("Enter number of worker nodes to add into pool: ")
    fmt.Scanf("%v", &node_count)

    new_nodepool := spio.NodePool{Name: "Test NodePool",
                                  ClusterID: clusterid,
                                  NodeCount: node_count,
                                  Platform: "coreos",
                                  State: "draft",
                                  Role: "worker" }
    if node_count > 0 {
        // Get machine types allowed for this provider
        m_options, err := client.GetMachSpecs(provider_name)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        fmt.Printf("Node size options for provider %s:\n", provider_name)
        for _, opt := range m_options {
            fmt.Println(opt)
        }
        var node_size string
        fmt.Printf("Enter node size: ")
        fmt.Scanf("%s", &node_size)
    
        // Validate machine type selection
        found := false
        for _, opt := range m_options {
            if node_size == opt {
                found = true
                break
            }
        }
        if found != true {
            fmt.Printf("Invalid option: %s\n", node_size)
            os.Exit(1)
        }
        new_nodepool.Size = node_size
    }

    // Create new nodepool
    _, err2 := client.CreateNodePool(orgid, clusterid, new_nodepool)
    if err2 != nil {
        fmt.Printf("Error: %v\n", err2)
        spio.ViewResp()
        os.Exit(1)
    }
    fmt.Printf("NodePool creation sent, building...\n")
}
