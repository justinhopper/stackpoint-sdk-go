package stackpointio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var lastHTTPResponse http.Response

// ViewResp breaks down last HTTP call's response for debugging
func ViewResp() {
	fmt.Println("response Status:", lastHTTPResponse.Status)
	fmt.Println("response StatusCode:", lastHTTPResponse.StatusCode)
	fmt.Println("response Headers:", lastHTTPResponse.Header)
	body, _ := ioutil.ReadAll(lastHTTPResponse.Body)
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

// NewClientFromEnv creates a new client from environment variables
func NewClientFromEnv() (*APIClient, error) {
	token := os.Getenv("SPC_API_TOKEN")
	if token == "" {
		return nil, errors.New("Missing token env in SPC_API_TOKEN")
	}
	endpoint := os.Getenv("SPC_BASE_API_URL")
	if endpoint == "" {
		return nil, errors.New("Missing endpoint env in SPC_BASE_API_URL")
	}
	return NewClient(token, endpoint), nil
}

// runRequest performs the actual HTTP request, returns bytes of body response
func (client *APIClient) runRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+client.token)
	req.Header.Set("User-Agent", "Stackpoint Go SDK")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)

	// Store response in case it's needed later
	lastHTTPResponse = *resp

	if err == nil && resp.StatusCode >= 400 {
		err = fmt.Errorf("Status code %d", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

// get sets up the GET request for runRequest to perform, returns bytes of body response
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

// delete sets up the DELETE request for runRequest to perform, returns bytes of body response (empty on success)
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

// post sets up the POST request for runRequest to perform, returns bytes of body response
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
	var profile UserProfile
	err = json.Unmarshal(content, &profile)
	if err != nil {
		return UserProfile{}, err
	}
	return profile, nil
}

// GetKeysets gets list of keysets for Org ID
func (client *APIClient) GetKeysets(orgid int) ([]Keyset, error) {
	resp, err := client.get(fmt.Sprintf("/orgs/%d/keysets", orgid))
	if err != nil {
		return nil, err
	}
	keysets := []Keyset{}
	err = json.Unmarshal(resp, &keysets)
	return keysets, err
}

// GetKeyset returns keyset for Org ID and Keyset ID
func (client *APIClient) GetKeyset(orgid, keysetid int) (Keyset, error) {
	var keyset Keyset
	resp, err := client.get(fmt.Sprintf("/orgs/%d/keysets/%d", orgid, keysetid))
	if err != nil {
		return keyset, err
	}
	err = json.Unmarshal(resp, &keyset)
	return keyset, err
}

// CreateKeyset creates keyset
func (client *APIClient) CreateKeyset(orgid int, dObj interface{}) ([]byte, error) {
	return client.post(fmt.Sprintf("/orgs/%d/keysets", orgid), dObj)
}

// DeleteKeyset deletes keyset
func (client *APIClient) DeleteKeyset(orgid, keysetid int) ([]byte, error) {
	return client.delete(fmt.Sprintf("/orgs/%d/keysets/%d", orgid, keysetid))
}

// GetMachSpecs returns list of machine types for cloud provider type
func (client *APIClient) GetMachSpecs(prov string) ([]string, error) {
	// TODO: this should be later replaced with a call to the API
	url := fmt.Sprintf("https://stackpointcloud-196003.appspot.com/specs/%s/", prov)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	content, err := client.runRequest(req)
	if err != nil {
		return nil, err
	}
	var specs ProviderSpecs
	err = json.Unmarshal(content, &specs)
	if err != nil {
		return nil, err
	}
	return specs.Machines, err
}

// GetClusters gets all clusters associated with an organization, returns list of Cluster objects
func (client *APIClient) GetClusters(organizationID int) ([]Cluster, error) {
	path := fmt.Sprintf("/orgs/%d/clusters", organizationID)
	content, err := client.get(path)
	if err != nil {
		return nil, err
	}
	var clusters []Cluster
	err = json.Unmarshal(content, &clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

// GetCluster gets a single cluster by primary ID and organization, returns Cluster object
func (client *APIClient) GetCluster(organizationID, clusterID int) (Cluster, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d", organizationID, clusterID)
	content, err := client.get(path)
	if err != nil {
		return Cluster{}, err
	}
	var cluster Cluster
	err = json.Unmarshal(content, &cluster)
	if err != nil {
		return Cluster{}, err
	}
	return cluster, nil
}

// CreateCluster requests cluster creation, returns Cluster object of newly created cluster
func (client *APIClient) CreateCluster(organizationID int, cluster Cluster) (Cluster, error) {
	path := fmt.Sprintf("/orgs/%d/clusters", organizationID)
	content, err := client.post(path, cluster)
	if err != nil {
		return Cluster{}, err
	}
	var newCluster Cluster
	err = json.Unmarshal(content, &newCluster)
	if err != nil {
		return Cluster{}, err
	}
	return newCluster, nil
}

// DeleteCluster deletes cluster, should return nothing on success
func (client *APIClient) DeleteCluster(orgid, clusterid int) ([]byte, error) {
	return client.delete(fmt.Sprintf("/orgs/%d/clusters/%d", orgid, clusterid))
}

// GetNodes gets the nodes associated with a cluster and organization, returns list of Node objects
func (client *APIClient) GetNodes(organizationID, clusterID int) ([]Node, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodes", organizationID, clusterID)
	content, err := client.get(path)
	if err != nil {
		return nil, err
	}
	var nodes []Node
	err = json.Unmarshal(content, &nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// GetNode retrieves data for a single node, returns Node object
func (client *APIClient) GetNode(organizationID, clusterID, nodeID int) (Node, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodes/%d", organizationID, clusterID, nodeID)
	content, err := client.get(path)
	if err != nil {
		return Node{}, err
	}
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
}

// AddNodes sends a request to add master nodes to a cluster, returns list of Node objects created
func (client *APIClient) AddNodes(organizationID, clusterID int, nodeAdd NodeAdd) ([]Node, error) {
	var newNodes []Node
	invalid := Validate(nodeAdd)
	if invalid != nil {
		return newNodes, invalid
	}
	path := fmt.Sprintf("/orgs/%d/clusters/%d/add_node", organizationID, clusterID)
	content, err := client.post(path, nodeAdd)
	if err != nil {
		return newNodes, err
	}
	err = json.Unmarshal(content, &newNodes)
	if err != nil {
		return newNodes, err
	}
	return newNodes, nil
}

// AddNodesToNodePool sends a request to add worker nodes to a nodepool, returns list of Node objects created
func (client *APIClient) AddNodesToNodePool(organizationID, clusterID, nodepoolID int, nodeAdd NodeAddToPool) ([]Node, error) {
	var newNodes []Node
	invalid := Validate(nodeAdd)
	if invalid != nil {
		return newNodes, invalid
	}
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodepools/%d/add", organizationID, clusterID, nodepoolID)
	content, err := client.post(path, nodeAdd)
	if err != nil {
		return newNodes, err
	}
	err = json.Unmarshal(content, &newNodes)
	if err != nil {
		return newNodes, err
	}
	return newNodes, nil
}

// GetNodePools gets the NodePools for a cluster, returns list of NodePool objects
func (client *APIClient) GetNodePools(organizationID, clusterID int) ([]NodePool, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodepools", organizationID, clusterID)
	content, err := client.get(path)
	if err != nil {
		return nil, err
	}
	var pools []NodePool
	err = json.Unmarshal(content, &pools)
	if err != nil {
		return nil, err
	}
	return pools, nil
}

// GetNodePool gets a NodePool for a cluster, returns NodePool object
func (client *APIClient) GetNodePool(organizationID, clusterID, nodepoolID int) (NodePool, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodepools/%d", organizationID, clusterID, nodepoolID)
	content, err := client.get(path)
	if err != nil {
		return NodePool{}, err
	}
	var pool NodePool
	err = json.Unmarshal(content, &pool)
	if err != nil {
		return NodePool{}, err
	}
	return pool, nil
}

// CreateNodePool creates a new nodepool for a cluster, returns NodePool object
func (client *APIClient) CreateNodePool(organizationID, clusterID int, pool NodePool) (NodePool, error) {
	var newNodePool NodePool
	invalid := Validate(pool)
	if invalid != nil {
		return newNodePool, invalid
	}
	path := fmt.Sprintf("/orgs/%d/clusters/%d/nodepools", organizationID, clusterID)
	content, err := client.post(path, pool)
	if err != nil {
		return newNodePool, err
	}
	err = json.Unmarshal(content, &newNodePool)
	if err != nil {
		return newNodePool, err
	}
	return newNodePool, nil
}

// DeleteNodePool deletes nodepool, returns nothing on success
func (client *APIClient) DeleteNodePool(orgid, clusterid, nodepoolid int) ([]byte, error) {
	return client.delete(fmt.Sprintf("/orgs/%d/clusters/%d/nodepools/%d", orgid, clusterid, nodepoolid))
}

// GetVolumes gets the Persistent Volumes attached to a cluster (NOT TESTED!!)
func (client *APIClient) GetVolumes(organizationID, clusterID int) ([]PersistentVolume, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/volumes", organizationID, clusterID)
	content, err := client.get(path)
	if err != nil {
		return nil, err
	}
	var volumes []PersistentVolume
	err = json.Unmarshal(content, &volumes)
	if err != nil {
		return nil, err
	}
	return volumes, nil
}

// GetLogs gets the BuildEventLogs for a cluster (NOT TESTED!!)
func (client *APIClient) GetLogs(organizationID, clusterID int) ([]BuildLogEntry, error) {
	path := fmt.Sprintf("/orgs/%d/clusters/%d/logs", organizationID, clusterID)
	content, err := client.get(path)
	if err != nil {
		return nil, err
	}
	var logs []BuildLogEntry
	err = json.Unmarshal(content, &logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// PostBuildLog adds a build log to the cluster (NOT TESTED!!)
func (client *APIClient) PostBuildLog(organizationID, clusterID int, log BuildLogEntry) (BuildLogEntry, error) {
	invalid := Validate(log)
	if invalid != nil {
		return BuildLogEntry{}, invalid
	}
	path := fmt.Sprintf("/orgs/%d/clusters/%d/buildlogs", organizationID, clusterID)
	content, err := client.post(path, log)
	if err != nil {
		return BuildLogEntry{}, err
	}
	var responseLog BuildLogEntry
	err = json.Unmarshal(content, &responseLog)
	if err != nil {
		return BuildLogEntry{}, err
	}
	return responseLog, nil
}

// PostAlert adds a alert message to the cluster as a build log (NOT TESTED!!)
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
