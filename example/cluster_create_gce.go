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
const ssh_keysetid = 3524
const gce_keysetid = 3553

func main() {
    client := spio.NewClient(mytoken, myurl)

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
