package main

import (
    "os"
    "fmt"
    spio "github.com/stackpoint-sdk-go/stackpointio"
)

const mytoken string = `98addc5550b98ab74499dd4cd0dc4dd03e0e0d0be82acb0213d1ec7ef2c79457`
const myurl string = `https://api-staging.stackpoint.io/`
const orgid int = 111

func main() {
    client := spio.NewClient(mytoken, myurl)

    new_key := spio.Key{Type: "pub",
                        Value: "FfgssHsds9u3rfsjOHDsssy8w3ehdsokFDDrR"}
    new_keyset := spio.Keyset{Name: "Test SSH Keyset",
                              Category: "user_ssh",
                              Entity: "",
                              Workspaces: []int{},
                              Keys: []spio.Key{new_key}}

    resp, err := client.CreateKeyset(orgid, new_keyset, "")
    fmt.Printf("CreateKeyset created: %+v\n", resp)
}
