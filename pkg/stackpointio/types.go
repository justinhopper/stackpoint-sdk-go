package stackpointio

import (
	"time"
)

// Organization is the top level of the hierarchy
type Organization struct {
	Name       string `json:"name"`
	PrimaryKey int    `json:"pk"`
}

// User is a stackpoint user
type User struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Provider struct for machine types call
type ProviderSpecs struct {
    Name     string   `json:"Name"`
    Label    string   `json:"Label"`
    Machines []string `json:"Machines"`
}

// Key struct
// Provider Type: aws{pub,pvt},azure{subscription,tenant,pub,pvt},do{token},gce{other},gke{other}
// Solution Type: sysdig{token},turbonomic{url,username,password,scope}
// SSH Type: pub
type Key struct {
    PrimaryKey  int    `json:"pk"`
    Type        string `json:"key_type"`
    Value       string `json:"key"`
    Fingerprint string `json:"fingerprint"`
    User        int    `json:"user"`
}

// Keyset struct
type Keyset struct {
    Name       string `json:"name"`
    PrimaryKey int    `json:"pk"`
    Category   string `json:"category"` // provider, solution, user_ssh
    Entity     string `json:"entity"` // provider{aws,azure,do,gce,gke},solution{sysdig,turbonomic}
    Org        int    `json:"org"`
    Workspaces []int  `json:"workspaces"`
    User       int    `json:"user"`
    IsDefault  bool   `json:"is_default"`
    Keys       []Key  `json:"keys"`
    Created    string `json:"created"`
}

// UserProfile includes detailed information about a StackPointCloud user
type UserProfile struct {
	User
	PrimaryKey     int            `json:"pk"`
	OrgMemberships []Organization `json:"org_memberships"`
	Keysets        []Keyset       `json:"keysets"`
}

// Solution is a application or process running with or on a kubernetes cluster,
// including "deis", "tectonic", "prometheus" and others
type Solution struct {
	PrimaryKey int    `json:"pk"`
	Solution   string `json:"solution"`
	Keyset     int    `json:"keyset,omitempty"` // only for turbonomic and sysdig
        MaxNodes   int    `json:"max_nodes"` // only for autoscaler
}

// Cluster describes a Kubernetes cluster in the StackPointCloud system
type Cluster struct {
	PrimaryKey         int        `json:"pk"`
	Name               string     `json:"name"`
	OrganizationKey    int        `json:"org"`
	InstanceID         string     `json:"instance_id"`
	Provider           string     `json:"provider"`
	ProviderKey        int        `json:"provider_keyset"`
	ProviderKeyName    string     `json:"provider_keyset_name"`
	ProviderNetworkID  string     `json:"provider_network_id"`
	ProviderNetworkCdr string     `json:"provider_network_cidr"`
	ProviderSubnetID   string     `json:"provider_subnet_id"`
	ProviderSubnetCidr string     `json:"provider_subnet_cidr"`
	ProviderBalancerID string     `json:"provider_balancer_id"`
	Region             string     `json:"region"`
	Zone               string     `json:"zone,omitempty"`
	State              string     `json:"state"`
	ProjectID          string     `json:"project_id,omitempty"`
	Owner              int        `json:"owner"`
	Notified           bool       `json:"notified,omitempty"`
	KubernetesVersion  string     `json:"k8s_version"`
	Created            time.Time  `json:"created"`
	Updated            time.Time  `json:"updated,omitempty"`
	DashboardEnabled   bool       `json:"k8s_dashboard_enabled"`
	DashboardInstalled bool       `json:"k8s_dashboard_installed"`
	KubeconfigPath     string     `json:"kubeconfig_path"`
	RbacEnabled        bool       `json:"k8s_rbac_enabled"`
	MasterCount        int        `json:"master_count"`
	WorkerCount        int        `json:"worker_count"`
	MasterSize         string     `json:"master_size"`
	WorkerSize         string     `json:"worker_size"`
	NodeCount          int        `json:"node_count"`
	EtcdType           string     `json:"etcd_type"`
	Platform           string     `json:"platform"`
	Image              string     `json:"image"`
	Channel            string     `json:"channel"`
        SSHKeySet          int        `json:"user_ssh_keyset"`
	Solutions          []Solution `json:"solutions"`
}

// Node describes a node in a cluster.  The string field Size is provider-specific
type Node struct {
	PrimaryKey int       `json:"pk"`
	Name       string    `json:"name"`
	ClusterID  int       `json:"cluster"`
	InstanceID string    `json:"instance_id"`
	Role       string    `json:"role"`
	Group      string    `json:"group_name,omitempty"`
	PrivateIP  string    `json:"private_ip"`
	PublicIP   string    `json:"public_ip"`
	Platform   string    `json:"platform"`
	Image      string    `json:"image"`
	Location   string    `json:"location"`
	Size       string    `json:"size"`
	State      string    `json:"state"` // draft, building, provisioned, running, deleting, deleted
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated,omitempty"`
}

// NodePool defines the characteristics of a grouping of nodes
type NodePool struct {
	PrimaryKey         int       `json:"pk"`
	Name               string    `json:"name"`
	ClusterID          int       `json:"cluster"`
	InstanceID         string    `json:"instance_id"`
	Size               string    `json:"instance_size"`
	CPU                string    `json:"cpu"`
	Memory             string    `json:"memory"`
	Labels             string    `json:"labels"`
	Autoscaled         bool      `json:"autoscaled"`
	MinCount           int       `json:"min_count"`
	MaxCount           int       `json:"max_count"`
	Zone               string    `json:"zone,omitempty"`
	ProviderSubnetID   string    `json:"provider_subnet_id"`
	ProviderSubnetCidr string    `json:"provider_subnet_cidr"`
	NodeCount          int       `json:"node_count"`
	Platform           string    `json:"platform"`
	Channel            string    `json:"channel"`
	Role               string    `json:"role"`  // []string{"master", "worker"}
	State              string    `json:"state"` // []string{"draft","active","failed","deleting","deleted"}
	Default            bool      `json:"is_default"`
	Created            time.Time `json:"created"`
	Updated            time.Time `json:"updated,omitempty"`
	Deleted            time.Time `json:"deleted,omitempty"`
}

// NodeAdd encapsulates the details of a call to add nodes to a cluster.
type NodeAdd struct {
	// Size is a cloudprovider-dependent string that describes the type of node to add
	Size string `json:"size"`
	// Count is the number of nodes to add
	Count int `json:"node_count"`
	// Group is the name of a stackpointcloud nodepool in the cluster
	Group string `json:"group,omitempty"`
	// NodePool is the id of the stackpointcloud nodepool.
	NodePoolID int `json:"node_pool"`
	// Role describes the role of this node - ["worker", "master"]
	Role string `json:"role,omitempty"`
	// Zone is a cloudprovider-dependent string for the node location
	Zone string `json:"zone,omitempty"`
	// ProviderSubnet... are cloudprovider-dependent network restrictions
	ProviderSubnetID   string `json:"provider_subnet_id,omitempty"`
	ProviderSubnetCidr string `json:"provider_subnet_cidr,omitempty"`
}

// PersistentVolume is the representation of a Kubernetes PersistentVolume in
// StackPointCloud, and includes details of the PersistentVolumeClaim
type PersistentVolume struct {
	PrimaryKey      int       `json:"pk"`
	Name            string    `json:"name"`
	ClusterID       int       `json:"cluster"`
	VolumeID        string    `json:"volume_id"`
	VolumeType      string    `json:"volume_type"`
	Claim           string    `json:"claim_name"`
	NameSpace       string    `json:"namespace"`
	SizeGB          int       `json:"size"`
	AccessMode      string    `json:"access_mode"`
	State           string    `json:"state"`
	RecyclingPolicy string    `json:"recycling_policy"`
	Owner           int       `json:"owner"`
	Deleted         bool      `json:"deleted,omitempty"`
	Backend         string    `json:"backend"`
	Created         time.Time `json:"created"`
}

// BuildLogEntry is an event log for the cluster build process
type BuildLogEntry struct {
	ClusterID     int       `json:"cluster"`
	EventCategory string    `json:"event_category"`
	EventType     string    `json:"event_type"`
	EventState    string    `json:"event_state"`
	Message       string    `json:"message"`
	Details       string    `json:"details,omitempty"`
	Reference     string    `json:"reference,omitempty"`
	Created       time.Time `json:"created,omitempty"`
}
