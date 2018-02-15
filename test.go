package main

import (
    "fmt"
    "net/http"
    "strings"
    "io/ioutil"
    "encoding/json"
)

const mytoken string = `98addc5550b98ab74499dd4cd0dc4dd03e0e0d0be82acb0213d1ec7ef2c79457`
const myurl string = `https://api-staging.stackpoint.io/`

// Struct to hold orgs
type Org struct {
    PrimaryKey int `json:"pk"`
    Name string `json:"name"`
    Slug string `json:"slug"`
    Logo string `json:"logo"`
}

// APIClient references an api token and an http endpoint
type APIClient struct {
    token      string
    endpoint   string
    httpClient *http.Client
}

// NewClient returns a new api client
func NewClient(token, endpoint string) *APIClient {
    c := APIClient {
        token:      token,
        endpoint:   strings.TrimRight(endpoint, "/"),
        httpClient: &http.Client{},
    }
    return &c
}

// RunGet runs the basic GET call to API
func (c *APIClient) RunGet(path string) ([]byte, error) {
    req, _ := http.NewRequest("GET", c.endpoint + path, nil)
    req.Header.Set("Authorization", "Bearer " + c.token)
    req.Header.Set("User-Agent", "Stackpoint Go SDK")
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient.Do(req)
    if err != nil { return nil, err }
    return ioutil.ReadAll(resp.Body)
}

// GetOrgs gets organizations list
func (c *APIClient) GetOrgs() ([]Org, error) {
    data, err := c.RunGet("/orgs")
    if err != nil { return nil, err }
    orgs := []Org{}
    err = json.Unmarshal(data, &orgs)
    return orgs, err
}

// GetOrg gets organization information
func (c *APIClient) GetOrg(id int) ([]Org, error) {
    var org []Org
    data, err := c.RunGet("/orgs/" + string(id))
    if err != nil { return org, err }
    err = json.Unmarshal(data, &org)
    return org, err
}

// ----------------MAIN---------------------
func main() {
    client := NewClient(mytoken, myurl)
    orgs, err := client.GetOrgs()
    fmt.Printf("orgs: %+v\nerr: %v\n", orgs, err)
    fmt.Printf("org ID: %v\n", orgs[0].PrimaryKey)
    org, err := client.GetOrg(orgs[0].PrimaryKey)
    fmt.Printf("org: %+v\nerr: %v\n", org, err)
}
