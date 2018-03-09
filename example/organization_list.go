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

func main() {
    // Set up HTTP client with API token and URL
    token := os.Getenv("SPC_API_TOKEN")
    endpoint := os.Getenv("SPC_BASE_API_URL")
    client := spio.NewClient(token, endpoint)

    orgs, err := client.GetOrganizations()
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    for i := 0; i < len(orgs); i++ {
        fmt.Printf("Org(%d): %v\n", orgs[i].PrimaryKey, orgs[i].Name)
    }

    var orgid int
    fmt.Printf("Enter org ID to inspect: ")
    fmt.Scanf("%d", &orgid)

    org, err := client.GetOrganization(orgid)
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    PrettyPrint(org)
}
