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

const mytoken string = `98addc5550b98ab74499dd4cd0dc4dd03e0e0d0be82acb0213d1ec7ef2c79457`
const myurl string = `https://api-staging.stackpoint.io/`

func main() {
    client := spio.NewClient(mytoken, myurl)

    orgs, err := client.GetOrganizations()
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    for i := 0; i < len(orgs); i++ {
        fmt.Printf("Org(%d): %v\n", orgs[i].PrimaryKey, orgs[i].Name)
    }

    var orgid int
    fmt.Printf("Enter org ID to inspect: ")
    _, err = fmt.Scanf("%d", &orgid)

    org, err := client.GetOrganization(orgid)
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    PrettyPrint(org)
}
