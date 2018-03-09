package main

import (
    "os"
    "fmt"
    spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
)

const orgid = 111
const provider = "gce"
const cluster_name = "Test GCE Cluster"

func main() {
    // Set up HTTP client with API token and URL
    token := os.Getenv("SPC_API_TOKEN")
    endpoint := os.Getenv("SPC_BASE_API_URL")
    client := spio.NewClient(token, endpoint)

    var ssh_keysetid int
    fmt.Sscanf(os.Getenv("SPC_SSH_KEYSET"), "%d", &ssh_keysetid)
    var gce_keysetid int
    fmt.Sscanf(os.Getenv("SPC_GCE_KEYSET"), "%d", &gce_keysetid)

    // Get list of machine types for provider
    m_options, err := client.GetMachSpecs(provider)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Printf("Node size options for provider %s:\n", provider)
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

    new_solution := spio.Solution{Solution: "helm_tiller"}
    new_cluster := spio.Cluster{Name: cluster_name,
                                Provider: provider,
                                ProviderKey: gce_keysetid,
                                MasterCount: 1,
                                MasterSize: node_size,
                                WorkerCount: 2,
                                WorkerSize: node_size,
                                Region: "us-west1-a",
                                State: "draft",
                                KubernetesVersion: "v1.8.7",
                                RbacEnabled: true,
                                DashboardEnabled: true,
                                EtcdType: "self_hosted",
                                Platform: "coreos",
                                Channel: "stable",
                                SSHKeySet: ssh_keysetid,
                                Solutions: []spio.Solution{new_solution} }

    _, err2 := client.CreateCluster(orgid, new_cluster)
    if err2 != nil {
        fmt.Printf("Error: %v\n", err2)
        spio.ViewResp()
        os.Exit(1)
    }
    fmt.Printf("Cluster created, building...\n")

}
