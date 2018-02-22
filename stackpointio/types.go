package stackpointio

import (
    "time"
)

// Struct to hold orgs
type Org struct {
    PrimaryKey int    `json:"pk"`
    Name       string `json:"name"`
    Slug       string `json:"slug"`
    Logo       string `json:"logo"`
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

// Solution struct
type Solution struct {
    PrimaryKey int       `json:"pk"`
    Solution   string    `json:"solution"`
    URL        string    `json:"url"`
    Username   string    `json:"username,omitempty"`
    Password   string    `json:"password,omitempty"`
    GitRepo    string    `json:"git_repo,omitempty"`
    GitPath    string    `json:"git_path,omitempty"`
    Created    time.Time `json:"created"`
    Updated    time.Time `json:"updated,omitempty"`
}

// Cluster struct
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
        Solutions          []Solution `json:"solutions"`
}

// Node struct
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

// NodePool struct
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
        Role               string    `json:"role"`  // {"master", "worker"}
        State              string    `json:"state"` // {"draft","active","failed","deleting","deleted"}
        Default            bool      `json:"is_default"`
        Created            time.Time `json:"created"`
        Updated            time.Time `json:"updated,omitempty"`
        Deleted            time.Time `json:"deleted,omitempty"`
}
