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
    var clusterid int
    fmt.Printf("Enter cluster ID to add node to: ")
    fmt.Scanf("%d", &clusterid)

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

    // Get type of node from user
    var node_role string
    fmt.Printf("Enter node role (worker or master) to use: ")
    fmt.Scanf("%v", &node_role)

    var new_node spio.NodeAdd
    if node_role == "master" {
        new_node.Count = 1
        new_node.Role = "master"
        new_node.Size = node_size
    } else if node_role == "worker" {
        // Get number of worker nodes to add
        var node_count int
        fmt.Printf("Enter number of worker nodes to add: ")
        fmt.Scanf("%v", &node_count)

        new_node.Count = node_count
        new_node.Role = "master"
        new_node.Size = node_size
    } else {
        fmt.Printf("Invalid node role: %s\n", node_role)
        os.Exit(1)
    }

    // Add new node
    _, err2 := client.AddNodes(orgid, clusterid, new_node)
    if err2 != nil {
        fmt.Printf("Error: %v\n", err2)
        spio.ViewResp()
        os.Exit(1)
    }
    fmt.Printf("Node creation sent, building...\n")
}
