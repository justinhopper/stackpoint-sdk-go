package stackpointio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

var lastHttpResponse http.Response

// ViewResp breaks down last HTTP call's response for debugging
func ViewResp() {
    fmt.Println("response Status:", lastHttpResponse.Status)
    fmt.Println("response StatusCode:", lastHttpResponse.StatusCode)
    fmt.Println("response Headers:", lastHttpResponse.Header)
    body, _ := ioutil.ReadAll(lastHttpResponse.Body)
    fmt.Println("response Body:", string(body))
}

// APIClient references an api token and an http endpoint
type APIClient struct {
	token      string
	endpoint   string
	httpClient *http.Client
}

// NewClient returns a new api client
func NewClient(token, endpoint string, client ...*http.Client) *APIClient {
	c := &APIClient{
		token:    token,
		endpoint: strings.TrimRight(endpoint, "/"),
	}
	if len(client) != 0 {
		c.httpClient = client[0]
	} else {
		c.httpClient = http.DefaultClient
	}
	return c
}

func (client *APIClient) runRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+client.token)
	req.Header.Set("User-Agent", "Stackpoint Go SDK")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)

        // Store response in case it's needed later
        lastHttpResponse = *resp

	if err == nil && resp.StatusCode >= 400 {
		err = fmt.Errorf("Status code %d", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func (client *APIClient) get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", client.endpoint+path, nil)
	if err != nil {
		return nil, err
	}
	content, err := client.runRequest(req)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (client *APIClient) delete(path string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", client.endpoint+path, nil)
	if err != nil {
		return nil, err
	}
	content, err := client.runRequest(req)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (client *APIClient) post(path string, dataObject interface{}) ([]byte, error) {
	data, err := json.Marshal(dataObject)
	if err != nil {
		return nil, err
        }
	req, err := http.NewRequest("POST", client.endpoint+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	content, err := client.runRequest(req)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// GetOrganizations retrieves data organizations that the client can access
func (client *APIClient) GetOrganizations() ([]Organization, error) {
	content, err := client.get("/orgs")
	if err != nil {
		return nil, err
	}
	var organizations []Organization
	err = json.Unmarshal(content, &organizations)
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

// GetOrganization retrieves data for a single organization
func (client *APIClient) GetOrganization(organizationID int) (Organization, error) {
	path := fmt.Sprintf("/orgs/%d", organizationID)
	content, err := client.get(path)
	if err != nil {
		return Organization{}, err
	}
	var organization Organization
	err = json.Unmarshal(content, &organization)
	if err != nil {
		return Organization{}, err
	}
	return organization, nil
}

// GetUser gets the StackPointCloud user
func (client *APIClient) GetUser() (User, error) {
	content, err := client.get("/rest-auth/user/")
	if err != nil {
		return User{}, err
	}
	glog.V(8).Info(string(content))
	var user User
	err = json.Unmarshal(content, &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// GetUserProfile gets details of StackPointCloud user profile
func (client *APIClient) GetUserProfile(username string) (UserProfile, error) {
	path := fmt.Sprintf("/userprofile/%s", username)
	content, err := client.get(path)
	if err != nil {
		return UserProfile{}, err
	}
	glog.V(8).Info(string(content))
	var profile UserProfile
	err = json.Unmarshal(content, &profile)
	if err != nil {
		return UserProfile{}, err
	}
	return profile, nil
}

// GetKeysets gets list of keysets for Org ID
func (c *APIClient) GetKeysets(orgid int) ([]Keyset, error) {
    resp, err := c.get(fmt.Sprintf("/orgs/%d/keysets", orgid))
    if err != nil { return nil, err }
    keysets := []Keyset{}
    err = json.Unmarshal(resp, &keysets)
    return keysets, err
}

// GetKeyset returns keyset for Org ID and Keyset ID
func (c *APIClient) GetKeyset(orgid, keysetid int) (Keyset, error) {
    var keyset Keyset
    resp, err := c.get(fmt.Sprintf("/orgs/%d/keysets/%d", orgid, keysetid))
    if err != nil { return keyset, err }
    err = json.Unmarshal(resp, &keyset)
    return keyset, err
}

// CreateKeyset creates keyset 
func (c *APIClient) CreateKeyset(orgid int, d_obj interface{}) ([]byte, error) {
    return c.post(fmt.Sprintf("/orgs/%d/keysets", orgid), d_obj)
}

// DeleteKeyset deletes keyset
func (c *APIClient) DeleteKeyset(orgid, keysetid int) ([]byte, error) {
    // Might return Error: 204 No Content, but Delete seems to happen fine
    return c.delete(fmt.Sprintf("/orgs/%d/keysets/%d", orgid, keysetid))
}

// GetMachSpecs returns list of machine types for cloud provider type
func (c *APIClient) GetMachSpecs(prov string) ([]string, error) {
    url := fmt.Sprintf("https://stackpointcloud-196003.appspot.com/specs/%s/", prov)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil { return nil, err }
    content, err := c.runRequest(req)
    if err != nil { return nil, err }
    var specs ProviderSpecs
    err = json.Unmarshal(content, &specs)
    if err != nil { return nil, err }
    return specs.Machines, err
}

// GetClusters gets all clusters associated with an organization
func (client *APIClient) GetClusters(organizationID int) ([]Cluster, error) {
	path := fmt.Sprintf("/orgs/%d/clusters", organizationID)
	content, err := client.get(path)
	if err != nil {
		return nil, err
	}
	glog.V(8).Info(string(content))
	var clusters []Cluster
	err = json.Unmarshal(content, &clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

// GetCluster gets a single cluster by primary ID and organization
func (client *APIClient) GetCluster(organizationID, clusterID int) (Cluster, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d", organizationID, clusterID)
	content, err := client.get(path)
	if err != nil {
		return Cluster{}, err
	}
	glog.V(8).Info(string(content))
	var cluster Cluster
	err = json.Unmarshal(content, &cluster)
	if err != nil {
		return Cluster{}, err
	}
	return cluster, nil
}

// CreateCluster requests cluster creation, returns immediately
func (client *APIClient) CreateCluster(organizationID int, cluster Cluster) (Cluster, error) {
	path := fmt.Sprintf("/orgs/%d/clusters", organizationID)
	content, err := client.post(path, cluster)
	if err != nil {
		return Cluster{}, err
	}
	glog.V(8).Info(string(content))
	var newCluster Cluster
	err = json.Unmarshal(content, &newCluster)
	if err != nil {
		return Cluster{}, err
	}
	return newCluster, nil
}

// DeleteCluster deletes cluster
func (c *APIClient) DeleteCluster(orgid, clusterid int) ([]byte, error) {
    // Might return Error: 204 No Content, but Delete seems to happen fine
    return c.delete(fmt.Sprintf("/orgs/%d/clusters/%d", orgid, clusterid))
}

// GetNodes gets the nodes associated with a cluster and organization
func (client *APIClient) GetNodes(organizationID, clusterID int) ([]Node, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodes", organizationID, clusterID)
	content, err := client.get(path)
	if err != nil {
		return nil, err
	}
	glog.V(8).Info(string(content))
	var nodes []Node
	err = json.Unmarshal(content, &nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// GetNode retrieves data for a single node
func (client *APIClient) GetNode(organizationID, clusterID, nodeID int) (Node, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodes/%d", organizationID, clusterID, nodeID)
	content, err := client.get(path)
	if err != nil {
		return Node{}, err
	}
	glog.V(8).Info(string(content))
	var node Node
	err = json.Unmarshal(content, &node)
	if err != nil {
		return Node{}, err
	}
	return node, nil
}

// DeleteNode makes an API call to begin deleting a node, and returns the contents of the web response
func (client *APIClient) DeleteNode(organizationID, clusterID, nodeID int) ([]byte, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodes/%d", organizationID, clusterID, nodeID)
	content, err := client.delete(path)
	return content, err
	// if err != nil {
	// 	return Node{}, err
	// }
	// glog.V(8).Info(string(content))
	// var node Node
	// err = json.Unmarshal(content, &node)
	// if err != nil {
	// 	return Node{}, err
	// }
	// return node, nil
}

// AddNodes sends a request to add nodes to a cluster, returns immediately with the Nodes under construction
func (client *APIClient) AddNodes(organizationID, clusterID int, nodeAdd NodeAdd) ([]Node, error) {
	invalid := Validate(nodeAdd)
	if invalid != nil {
		return []Node{}, invalid
	}
	path := fmt.Sprintf("/orgs/%d/clusters/%d/add_node", organizationID, clusterID)
	content, err := client.post(path, nodeAdd)
	if err != nil {
		return []Node{}, err
	}
	glog.V(8).Info("add node response: " + string(content))
	var response []Node
	err = json.Unmarshal(content, &response)
	if err != nil {
		return []Node{}, err
	}
	return response, nil
}

// AddNodesToNodePool sends a request to add nodes to a nodepool, returns immediately with the Nodes under construction
func (client *APIClient) AddNodesToNodePool(organizationID, clusterID, nodepoolID int, nodeAdd NodeAdd) ([]Node, error) {
	invalid := Validate(nodeAdd)
	if invalid != nil {
		return []Node{}, invalid
	}
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodepools/%d/add", organizationID, clusterID, nodepoolID)
	content, err := client.post(path, nodeAdd)
	if err != nil {
		return []Node{}, err
	}
	glog.V(8).Info("add node response: " + string(content))
	var response []Node
	err = json.Unmarshal(content, &response)
	if err != nil {
		return []Node{}, err
	}
	return response, nil
}

// GetNodePools gets the NodePools for a cluster
func (client *APIClient) GetNodePools(organizationID, clusterID int) ([]NodePool, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodepools", organizationID, clusterID)
	content, err := client.get(path)
	if err != nil {
		return nil, err
	}
	glog.V(8).Info(string(content))
	var pools []NodePool
	err = json.Unmarshal(content, &pools)
	if err != nil {
		return nil, err
	}
	for i, pool := range pools {
		pools[i] = setNodePoolSpecs(pool)
	}
	return pools, nil
}

// GetNodePool gets a NodePool for a cluster
func (client *APIClient) GetNodePool(organizationID, clusterID, nodepoolID int) (NodePool, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodepools/%d", organizationID, clusterID, nodepoolID)
	content, err := client.get(path)
	if err != nil {
		return NodePool{}, err
	}
	glog.V(8).Info(string(content))
	var pool NodePool
	err = json.Unmarshal(content, &pool)
	if err != nil {
		return NodePool{}, err
	}
	return setNodePoolSpecs(pool), nil
}

// CreateNodePool gets a NodePool for a cluster
func (client *APIClient) CreateNodePool(organizationID, clusterID int, pool NodePool) (NodePool, error) {
	invalid := Validate(pool)
	if invalid != nil {
		return NodePool{}, invalid
	}
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodepools", organizationID, clusterID)
	content, err := client.post(path, pool)
	if err != nil {
		return NodePool{}, err
	}
	glog.V(8).Info(string(content))
	var created NodePool
	err = json.Unmarshal(content, &created)
	if err != nil {
		return NodePool{}, err
	}
	return created, nil
}

// DeleteNodePool deletes nodepool
func (c *APIClient) DeleteNodePool(orgid, clusterid, nodepoolid int) ([]byte, error) {
    // Might return Error: 204 No Content, but Delete seems to happen fine
    return c.delete(fmt.Sprintf("/orgs/%d/clusters/%d/nodepools/%d", orgid, clusterid, nodepoolid))
}

// sets the cpu and memory of a pool if they are not configured already,
// based on the type ("size")
func setNodePoolSpecs(pool NodePool) NodePool {
	if pool.CPU == "" {
		cpu := getMachineTypeCPU(pool.Size)
		if cpu != 0 {
			pool.CPU = strconv.Itoa(cpu)
		}
	}
	if pool.Memory == "" {
		memory := getMachineTypeMemory(pool.Size)
		if memory != 0 {
			pool.Memory = strconv.Itoa(memory)
		}
	}
	return pool
}

// GetVolumes gets the Persistent Volumes attached to a cluster
func (client *APIClient) GetVolumes(organizationID, clusterID int) ([]PersistentVolume, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/volumes", organizationID, clusterID)
	content, err := client.get(path)
	if err != nil {
		return nil, err
	}
	glog.V(8).Info(string(content))
	var volumes []PersistentVolume
	err = json.Unmarshal(content, &volumes)
	if err != nil {
		return nil, err
	}
	return volumes, nil
}

// GetLogs gets the BuildEventLogs for a cluster
func (client *APIClient) GetLogs(organizationID, clusterID int) ([]BuildLogEntry, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/logs", organizationID, clusterID)
	content, err := client.get(path)
	if err != nil {
		return nil, err
	}
	glog.V(8).Info(string(content))
	var logs []BuildLogEntry
	err = json.Unmarshal(content, &logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// PostBuildLog adds a build log to the cluster
func (client *APIClient) PostBuildLog(organizationID, clusterID int, log BuildLogEntry) (BuildLogEntry, error) {
	invalid := Validate(log)
	if invalid != nil {
		return BuildLogEntry{}, invalid
	}
	path := fmt.Sprintf("/orgs/%d/clusters/%d/buildlogs", organizationID, clusterID)
	// data, _ := json.Marshal(log)
	// glog.V(8).Infof("PostBuildLog received %s", data)
	content, err := client.post(path, log)
	if err != nil {
		return BuildLogEntry{}, err
	}
	glog.V(8).Infof("PostBuildLog response %s", string(content))
	var responseLog BuildLogEntry
	err = json.Unmarshal(content, &responseLog)
	if err != nil {
		return BuildLogEntry{}, err
	}
	return responseLog, nil
}

// PostAlert adds a alert message to the cluster as a build log
func (client *APIClient) PostAlert(organizationID, clusterID int, message, details, reference string) (BuildLogEntry, error) {
	alert := BuildLogEntry{
		ClusterID:     clusterID,
		EventCategory: "kubernetes",
		EventType:     "provider_communication",
		EventState:    "failure",
		Message:       message,
		Details:       details,
		Reference:     reference,
	}
	return client.PostBuildLog(organizationID, clusterID, alert)
}
