package main

import (
    "os"
    "fmt"
    spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
)

const orgid = 111

func main() {
    token := os.Getenv("CLUSTER_API_TOKEN")
    endpoint := os.Getenv("SPC_BASE_API_URL")
    client := spio.NewClient(token, endpoint)

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
    _, err = fmt.Scanf("%d", &clusterid)

    // Get cluster provider from selection
    var provider_name string
    for i := 0; i < len(clusters); i++ {
        if clusters[i].PrimaryKey == clusterid {
            provider_name = clusters[i].Provider
            break;
        }
    }

    // Get machine types allowed for this provider
    m_options, err := client.GetMachSpecs(provider_name)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
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

    // Get list of nodepools to select from
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
    fmt.Printf("Enter nodepool ID to add node to: ")
    fmt.Scanf("%d", &nodepoolid)

    // Get number of nodes to add
    var node_count int
    fmt.Printf("Enter number of nodes to add: ")
    fmt.Scanf("%d", &node_count)

    new_node := spio.NodeAdd{Size: node_size,
                             Count: node_count,
                             Role: "worker",
                             NodePoolID: nodepoolid}

    _, err2 := client.AddNodesToNodePool(orgid, clusterid, nodepoolid, new_node)
    if err2 != nil {
        fmt.Printf("Error: %v\n", err2)
        spio.ViewResp()
        os.Exit(1)
    }
    fmt.Printf("Node creation sent, building...\n")
}
