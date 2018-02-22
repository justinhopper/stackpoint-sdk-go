package main

import (
    "os"
    "fmt"
    spio "github.com/stackpoint-sdk-go/stackpointio"
)

const mytoken string = `98addc5550b98ab74499dd4cd0dc4dd03e0e0d0be82acb0213d1ec7ef2c79457`
const myurl string = `https://api-staging.stackpoint.io/`

func main() {
    client := spio.NewClient(mytoken, myurl)

    orgs, err := client.GetOrgs()
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    //fmt.Printf("orgs: %+v\n", orgs)
    orgid := orgs[0].PrimaryKey;
    fmt.Printf("org ID: %v\n", orgid)

    org, err := client.GetOrg(orgid)
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    fmt.Printf("org: %+v\n", org)

    keysets, err := client.GetKeysets(orgid)
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    fmt.Printf("keysets: %+v\n", keysets)
    keysetid := keysets[1].PrimaryKey
    fmt.Printf("keyset ID: %v\n", keysetid)

    keyset, err := client.GetKeyset(orgid, keysetid)
    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
    fmt.Printf("keyset: %+v\n", keyset)

//    err2 := client.DeleteKeyset(orgid, keysetid)
//    if err2 != nil { fmt.Printf("Error: %v\n", err2); os.Exit(1) }
//    fmt.Printf("keyset ID %d is deleted\n", keysetid)

//    keysets2, err := client.GetKeysets(orgid)
//    if err != nil { fmt.Printf("Error: %v\n", err); os.Exit(1) }
//    fmt.Printf("keysets2: %+v\n", keysets2)

//    new_key := spio.Key{Type: "pub",
//                        Value: "FfgssHsds9u3rfsjOHDsssy8w3ehdsokFDDrR"}
//    new_keyset := spio.Keyset{Name: "Test SSH Keyset",
//                              Category: "user_ssh",
//                              Entity: "",
//                              Workspaces: []int{},
//                              Keys: []spio.Key{new_key}}
//    var jsonStr = string(`
//{
//  "name": "Test SSH Keyset",
//  "category": "provider",
//  "entity": "aws",
//  "workspaces": [],
//  "keys": [
//    {
//      "key_type": "pub",
//      "key": "AKIAIXH7J9KGB56VWZGA"
//    },
//    {
//      "key_type": "pvt",
//      "key": "FfgssHsds9u3rfsjOHDsssy8w3ehdsokFDDrR"
//    }
//  ]
//}
//`)
//    resp, err := client.CreateKeyset(orgid, nil, jsonStr)
//    spio.ViewResp()
//    fmt.Printf("CreateKeyset on string probably failed: %+v\n", resp)

//    resp2, err := client.CreateKeyset(orgid, new_keyset, "")
//    spio.ViewResp()
//    fmt.Printf("CreateKeyset on object maybe worked: %+v\n", resp2)

    os.Exit(0)

    clusters, err := client.GetClusters(orgid)
    if err != nil { fmt.Printf("Error: %v\n"); os.Exit(1) }
    fmt.Printf("clusters: %+v\n", clusters)
    clusterid := clusters[0].PrimaryKey;
    fmt.Printf("cluster ID: %v\n", clusterid)
}

