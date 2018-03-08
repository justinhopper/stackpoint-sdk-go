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

func main() {
    client := spio.NewClient(mytoken, myurl)

    new_key := spio.Key{Type: "token",
                        Value: "b5fa8ea386ce5a89308b01fd3603a968e493915e54ba20e87cf9abac615834c4"}
    new_keyset := spio.Keyset{Name: "Test DigitalOcean Keyset",
                              Category: "provider",
                              Entity: "do",
                              Workspaces: []int{},
                              Keys: []spio.Key{new_key}}

    resp, err := client.CreateKeyset(orgid, new_keyset)
    if err != nil { 
        fmt.Printf("Error: %v\n", err)
        spio.ViewResp()
        os.Exit(1)
    }
    fmt.Printf("CreateKeyset created\n")
}
