package stackpointio

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getSimpleGetMuxDummy(path, responseBody string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if http.MethodGet != r.Method {
			w.WriteHeader(http.StatusNotImplemented)
		}
		fmt.Fprint(w, responseBody)
	}))

	return mux
}

func TestGetOrganizations(t *testing.T) {

	mux := getSimpleGetMuxDummy("/orgs", "[{\"name\":\"Misty Fire\",\"pk\":1}]")
	ts := httptest.NewServer(mux)
	defer ts.Close()

	token := "not used"
	client := NewClient(token, ts.URL)

	orgs, err := client.GetOrganizations()
	require.Nil(t, err)

	assert.Equal(t, 1, len(orgs))
	assert.Equal(t, 1, orgs[0].PrimaryKey)
	assert.Equal(t, "Misty Fire", orgs[0].Name)
}

func TestGetOrganization(t *testing.T) {

	mux := getSimpleGetMuxDummy("/orgs/1", "{\"name\":\"Misty Fire\",\"pk\":1}")
	ts := httptest.NewServer(mux)
	defer ts.Close()

	token := "not used"
	client := NewClient(token, ts.URL)

	org, err := client.GetOrganization(1)
	require.Nil(t, err)

	assert.Equal(t, 1, org.PrimaryKey)
	assert.Equal(t, "Misty Fire", org.Name)
}

func TestAddNode(t *testing.T) {

	organizationKey := 123
	clusterKey := 456
	nodeAdd := NodeAdd{Count: 1, Size: "t2.medium"}

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	responseText := "{\"channel\":\"stable\",\"cluster\":504,\"created\":\"2016-09-27T22:09:57.089819Z\",\"image\":\"ami-06af7f66\", \"instance_id\": \"spcvd7ah21-worker-1\", \"location\": \"us-west-2:us-west-2a\", \"pk\": 1031,\"platform\": \"coreos\",\"private_ip\": \"172.23.1.209\", \"public_ip\": \"54.70.151.25\",\"role\": \"worker\",\"group_name\":\"autoscaling\",\"size\":\"t2.medium\",\"state\":\"draft\",\"updated\":\"2016-09-27T22:09:57.089836Z\"}"

	mux.Handle(fmt.Sprintf("/orgs/%d/clusters/%d/add_node", organizationKey, clusterKey),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			postedData, err := ioutil.ReadAll(r.Body)
			require.Nil(t, err)
			assert.True(t, 0 < len(postedData), "postedData non-zero length")
			fmt.Fprint(w, responseText)
		}))

	ts := httptest.NewServer(mux)
	defer ts.Close()

	token := "not used"
	client := NewClient(token, ts.URL)

	node, err := client.AddNodes(organizationKey, clusterKey, nodeAdd)
	require.Nil(t, err)
	assert.Equal(t, "draft", node.State, "returned Node in state \"draft\"")
	assert.Equal(t, "autoscaling", node.Group, "returned Node in group \"autoscaling\"")

}

func TestGetNodepool(t *testing.T) {

	organizationKey := 6
	clusterKey := 1665
	nodepoolKey := 80

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	responseText := "{\"pk\":80,\"cluster\":1665,\"name\":\"Default Worker Pool\",\"instance_id\":\"spc5u92han-pool-1\",\"instance_size\":\"2gb\",\"platform\":\"coreos\",\"channel\":\"stable\",\"zone\":\"nyc1\",\"provider_subnet_id\":\"\",\"provider_subnet_cidr\":\"\",\"node_count\":4,\"role\":\"worker\",\"state\":\"active\",\"is_default\":true,\"created\":\"2017-05-30T15:59:09.030257Z\",\"updated\":\"2017-05-30T19:17:01.219649Z\"}"

	mux.Handle(fmt.Sprintf("/orgs/%d/clusters/%d/nodepools/%d", organizationKey, clusterKey, nodepoolKey),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			// postedData, err := ioutil.ReadAll(r.Body)
			// require.Nil(t, err)
			// assert.True(t, 0 == len(postedData), "postedData zero length")
			fmt.Fprint(w, responseText)
		}))

	ts := httptest.NewServer(mux)
	defer ts.Close()

	token := "not used"
	client := NewClient(token, ts.URL)

	nodepool, err := client.GetNodepool(organizationKey, clusterKey, nodepoolKey)
	require.Nil(t, err)
	assert.Equal(t, "Default Worker Pool", nodepool.Name, "nodepool.Name")
	assert.Equal(t, "spc5u92han-pool-1", nodepool.Group, "nodepool.Group")
	assert.Equal(t, "", nodepool.SubnetCIDR)
	assert.Equal(t, "", nodepool.SubnetCIDR)
	assert.Equal(t, "", nodepool.SubnetID)
	assert.Equal(t, clusterKey, nodepool.ClusterID)
	assert.Equal(t, 4, nodepool.Count)
	assert.Equal(t, "coreos", nodepool.Platform)
	assert.Equal(t, "2gb", nodepool.Size)
	assert.Equal(t, "worker", nodepool.Role)
	assert.Equal(t, true, nodepool.IsDefault)
	assert.Equal(t, "nyc1", nodepool.Zone)
}

func TestListNodepools(t *testing.T) {

	organizationKey := 6
	clusterKey := 1665

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	responseText := "[{\"pk\":80,\"cluster\":1665,\"name\":\"Default Worker Pool\",\"instance_id\":\"spc5u92han-pool-1\",\"instance_size\":\"2gb\",\"platform\":\"coreos\",\"channel\":\"stable\",\"zone\":\"nyc1\",\"provider_subnet_id\":\"\",\"provider_subnet_cidr\":\"\",\"node_count\":4,\"role\":\"worker\",\"state\":\"active\",\"is_default\":true,\"created\":\"2017-05-30T15:59:09.030257Z\",\"updated\":\"2017-05-30T19:17:01.219649Z\"}]"

	mux.Handle(fmt.Sprintf("/orgs/%d/clusters/%d/nodepools", organizationKey, clusterKey),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			// postedData, err := ioutil.ReadAll(r.Body)
			// require.Nil(t, err)
			// assert.True(t, 0 == len(postedData), "postedData zero length")
			fmt.Fprint(w, responseText)
		}))

	ts := httptest.NewServer(mux)
	defer ts.Close()

	token := "not used"
	client := NewClient(token, ts.URL)

	nodepools, err := client.ListNodepools(organizationKey, clusterKey)
	require.Nil(t, err)

	assert.Equal(t, 1, len(nodepools), "one nodepool listed")

	nodepool := nodepools[0]

	assert.Equal(t, "Default Worker Pool", nodepool.Name, "nodepool.Name")
	assert.Equal(t, "spc5u92han-pool-1", nodepool.Group, "nodepool.Group")
	assert.Equal(t, "", nodepool.SubnetCIDR)
	assert.Equal(t, "", nodepool.SubnetCIDR)
	assert.Equal(t, "", nodepool.SubnetID)
	assert.Equal(t, clusterKey, nodepool.ClusterID)
	assert.Equal(t, 4, nodepool.Count)
	assert.Equal(t, "coreos", nodepool.Platform)
	assert.Equal(t, "2gb", nodepool.Size)
	assert.Equal(t, "worker", nodepool.Role)
	assert.Equal(t, true, nodepool.IsDefault)
	assert.Equal(t, "nyc1", nodepool.Zone)
}
