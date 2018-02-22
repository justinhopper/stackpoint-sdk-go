package stackpointio

import (
    "fmt"
    "bytes"
    "net/http"
    "strings"
    "io/ioutil"
    "encoding/json"
)

const userAgent = `Stackpoint Go SDK`
var lastHttpResponse http.Response

// PrettyPrint to break down objects
func PrettyPrint(v interface{}) {
      b, _ := json.MarshalIndent(v, "", "  ")
      println(string(b))
}

// ViewResp breaks down last HTTP call's response for debugging
func ViewResp() {
    fmt.Println("response Status:", lastHttpResponse.Status)
    fmt.Println("response StatusCode:", lastHttpResponse.StatusCode)
    fmt.Println("response Headers:", lastHttpResponse.Header)
    body, _ := ioutil.ReadAll(lastHttpResponse.Body)
    fmt.Println("response Body:", string(body))
}

// ClientSt holds token and endpoint, and http client reference
type ClientSt struct {
    token      string
    endpoint   string
    httpClient *http.Client
}

// NewClient returns a new api client
func NewClient(token, endpoint string) *ClientSt {
    c := ClientSt {
        token:      token,
        endpoint:   strings.TrimRight(endpoint, "/"),
        httpClient: &http.Client{},
    }
    return &c
}

// RunRequest executes the HTTP request
func (c *ClientSt) RunRequest(req *http.Request) ([]byte, error) {
    req.Header.Set("Authorization", "Bearer " + c.token)
    req.Header.Set("User-Agent", userAgent)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient.Do(req)
    lastHttpResponse = *resp 
    if err != nil { return nil, err }
    if resp.StatusCode != 200 { return nil, fmt.Errorf("RunRequest error: %v", resp.Status) }
    return ioutil.ReadAll(resp.Body)
}

// RunGet runs the basic GET call to API
func (c *ClientSt) RunGet(path string) ([]byte, error) {
    fmt.Printf("RunGet sees: %v\n", path)
    req, _ := http.NewRequest("GET", c.endpoint + path, nil)
    return c.RunRequest(req)
}

// RunPostObj runs POST call to API for object and converts to json
func (c *ClientSt) RunPostObj(path string, d interface{}) ([]byte, error) {
    fmt.Printf("RunPost sees: %v\n", path)
    data, err := json.Marshal(d)
    if err != nil { return nil, err }
    req, _ := http.NewRequest("POST", c.endpoint + path, bytes.NewReader(data))
    return c.RunRequest(req)
}

// RunPostStr runs POST call to API for json string
func (c *ClientSt) RunPostStr(path, data string) ([]byte, error) {
    fmt.Printf("RunPost sees: %v\n", path)
    req, _ := http.NewRequest("POST", c.endpoint + path, bytes.NewBufferString(data))
    resp, err := c.RunRequest(req)
    return resp, err
}

// RunDelete runs DELETE call to API for URI
func (c *ClientSt) RunDelete(path string) error {
    fmt.Printf("RunDelete sees: %v\n", path)
    req, _ := http.NewRequest("DELETE", c.endpoint + path, nil)
    _, err := c.RunRequest(req)
    return err
}

// GetOrgs gets organizations list
func (c *ClientSt) GetOrgs() ([]Org, error) {
    resp, err := c.RunGet("/orgs")
    if err != nil { return nil, err }
    orgs := []Org{}
    err = json.Unmarshal(resp, &orgs)
    return orgs, err
}

// GetOrg gets information for org ID
func (c *ClientSt) GetOrg(orgid int) (Org, error) {
    var org Org
    resp, err := c.RunGet("/orgs/" + fmt.Sprintf("%d", orgid))
    if err != nil { return org, err }
    err = json.Unmarshal(resp, &org)
    return org, err
}

// GetKeysets gets list of keysets for Org ID
func (c *ClientSt) GetKeysets(orgid int) ([]Keyset, error) {
    resp, err := c.RunGet("/orgs/" + fmt.Sprintf("%d", orgid) + "/keysets")
    if err != nil { return nil, err }
    keysets := []Keyset{}
    err = json.Unmarshal(resp, &keysets)
    return keysets, err
}

// GetKeyset returns keyset for Org ID and Keyset ID
func (c *ClientSt) GetKeyset(orgid, keysetid int) (Keyset, error) {
    var keyset Keyset
    resp, err := c.RunGet("/orgs/" + fmt.Sprintf("%d", orgid) + "/keysets/" + fmt.Sprintf("%d", keysetid))
    if err != nil { return keyset, err }
    err = json.Unmarshal(resp, &keyset)
    return keyset, err
}

// CreateKeyset creates keyset 
func (c *ClientSt) CreateKeyset(orgid int, d_obj interface{}, d_str string) ([]byte, error) {
    if d_obj != nil { 
        return c.RunPostObj("/orgs/" + fmt.Sprintf("%d", orgid) + "/keysets", d_obj) 
    } else { 
        return c.RunPostStr("/orgs/" + fmt.Sprintf("%d", orgid) + "/keysets", d_str) 
    }
}

// DeleteKeyset deletes keyset
func (c *ClientSt) DeleteKeyset(orgid, keysetid int) error {
    // Might return Error: 204 No Content, but Delete seems to happen fine
    return c.RunDelete("/orgs/" + fmt.Sprintf("%d", orgid) + "/keysets/" + fmt.Sprintf("%d", keysetid))
}

// GetClusters gets clusters list for org ID
func (c *ClientSt) GetClusters(orgid int) ([]Cluster, error) {
    resp, err := c.RunGet("/orgs/" + fmt.Sprintf("%d", orgid) + "/clusters")
    if err != nil { return nil, err }
    clusters := []Cluster{}
    err = json.Unmarshal(resp, &clusters)
    return clusters, err
}

// GetCluster gets information for cluster ID and org ID
func (c *ClientSt) GetCluster(orgid, clusterid int) (Cluster, error) {
    var cluster Cluster
    resp, err := c.RunGet("/orgs/" + fmt.Sprintf("%d", orgid) + "/clusters/" + fmt.Sprintf("%d", clusterid))
    if err != nil { return cluster, err }
    err = json.Unmarshal(resp, &cluster)
    return cluster, err
}
