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

    keysets, err := client.GetKeysets(orgid)
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    for i := 0; i < len(keysets); i++ {
        fmt.Printf("Keysets(%d): %v\n", keysets[i].PrimaryKey, keysets[i].Name)
    }

    var keysetid int
    fmt.Printf("Enter keyset ID to delete: ")
    fmt.Scanf("%d", &keysetid)

    client.DeleteKeyset(orgid, keysetid)
    fmt.Printf("Keyset should delete shortly\n")
}
