package main

import (
    "os"
    "fmt"
    spio "github.com/StackPointCloud/stackpoint-sdk-go/pkg/stackpointio"
)

const mytoken = `98addc5550b98ab74499dd4cd0dc4dd03e0e0d0be82acb0213d1ec7ef2c79457`
const myurl = `https://api-staging.stackpoint.io/`
const orgid = 111
const ssh_keysetid = 3524
const provider_name = "do"
const do_keysetid = 3556

func main() {
    client := spio.NewClient(mytoken, myurl)

    // Get list of machine types for provider
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

    new_solution := spio.Solution{Solution: "helm_tiller"}
    new_cluster := spio.Cluster{Name: "Test DigitalOcean Cluster",
                                Provider: "do",
                                ProviderKey: do_keysetid,
                                MasterCount: 1,
                                MasterSize: node_size,
                                WorkerCount: 2,
                                WorkerSize: node_size,
                                Region: "nyc1",
                                State: "draft",
                                KubernetesVersion: "v1.8.3",
                                RbacEnabled: true,
                                DashboardEnabled: true,
                                EtcdType: "self_hosted",
                                Platform: "coreos",
                                Channel: "stable",
                                SSHKeySet: ssh_keysetid,
                                Solutions: []spio.Solution{new_solution} }

    resp, err := client.CreateCluster(orgid, new_cluster)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        fmt.Printf("Cluster: %+v\n", resp)
        spio.ViewResp()
        os.Exit(1)
    }
    fmt.Printf("Cluster created, building...\n")
}
