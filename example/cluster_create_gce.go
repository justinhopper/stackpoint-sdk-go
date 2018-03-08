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

const orgid = 111
const ssh_keysetid = "" // ENTER YOUR SSH KEYSET ID
const gce_keysetid = "" // ENTER YOUR DIGITAL OCEAN KEYSET ID

func main() {
    token := os.Getenv("CLUSTER_API_TOKEN")
    endpoint := os.Getenv("SPC_BASE_API_URL")
    client := spio.NewClient(token, endpoint)

    new_solution := spio.Solution{Solution: "helm_tiller"}
    new_cluster := spio.Cluster{Name: "Test GCE Cluster",
                                Provider: "gce",
                                ProviderKey: gce_keysetid,
                                MasterCount: 1,
                                MasterSize: "n1-standard-1",
                                WorkerCount: 2,
                                WorkerSize: "n1-standard-1",
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

    _, err := client.CreateCluster(orgid, new_cluster)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        spio.ViewResp()
        os.Exit(1)
    }
    fmt.Printf("Cluster created, building...\n")

}
