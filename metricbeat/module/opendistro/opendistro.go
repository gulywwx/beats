package opendistro

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/metricbeat/helper"
)

const ModuleName = "opendistro"

type Info struct {
	ClusterName string `json:"cluster_name"`
	ClusterID   string `json:"cluster_uuid"`
	Version     struct {
		Number *common.Version `json:"number"`
	} `json:"version"`
}

func IsMaster(http *helper.HTTP, uri string) (bool, error) {

	node, err := getNodeName(http, uri)
	if err != nil {
		return false, err
	}

	master, err := getMasterName(http, uri)
	if err != nil {
		return false, err
	}

	return master == node, nil
}

func getNodeName(http *helper.HTTP, uri string) (string, error) {
	content, err := fetchPath(http, uri, "/_nodes/_local/nodes", "")
	if err != nil {
		return "", err
	}

	nodesStruct := struct {
		Nodes map[string]interface{} `json:"nodes"`
	}{}

	json.Unmarshal(content, &nodesStruct)

	// _local will only fetch one node info. First entry is node name
	for k := range nodesStruct.Nodes {
		return k, nil
	}
	return "", fmt.Errorf("No local node found")
}

func getMasterName(http *helper.HTTP, uri string) (string, error) {
	// TODO: evaluate on why when run with ?local=true request does not contain master_node field
	content, err := fetchPath(http, uri, "_cluster/state/master_node", "")
	if err != nil {
		return "", err
	}

	clusterStruct := struct {
		MasterNode string `json:"master_node"`
	}{}

	json.Unmarshal(content, &clusterStruct)

	return clusterStruct.MasterNode, nil
}

func fetchPath(http *helper.HTTP, uri, path string, query string) ([]byte, error) {
	defer http.SetURI(uri)

	// Parses the uri to replace the path
	u, _ := url.Parse(uri)
	u.Path = path
	u.RawQuery = query

	// Http helper includes the HostData with username and password
	http.SetURI(u.String())
	return http.FetchContent()
}

func GetInfo(http *helper.HTTP, uri string) (*Info, error) {

	content, err := fetchPath(http, uri, "/", "")
	if err != nil {
		return nil, err
	}

	info := &Info{}
	err = json.Unmarshal(content, &info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// GetMasterNodeID returns the ID of the Elasticsearch cluster's master node
func GetMasterNodeID(http *helper.HTTP, resetURI string) (string, error) {
	content, err := fetchPath(http, resetURI, "_nodes/_master", "filter_path=nodes.*.name")
	if err != nil {
		return "", err
	}

	var response struct {
		Nodes map[string]interface{} `json:"nodes"`
	}

	if err := json.Unmarshal(content, &response); err != nil {
		return "", err
	}

	for nodeID, _ := range response.Nodes {
		return nodeID, nil
	}

	return "", errors.New("could not determine master node ID")
}

// IsMLockAllEnabled returns if the given Elasticsearch node has mlockall enabled
func IsMLockAllEnabled(http *helper.HTTP, resetURI, nodeID string) (bool, error) {
	content, err := fetchPath(http, resetURI, "_nodes/"+nodeID, "filter_path=nodes.*.process.mlockall")
	if err != nil {
		return false, err
	}

	var response map[string]map[string]map[string]map[string]bool
	err = json.Unmarshal(content, &response)
	if err != nil {
		return false, err
	}

	for _, nodeInfo := range response["nodes"] {
		mlockall := nodeInfo["process"]["mlockall"]
		return mlockall, nil
	}

	return false, fmt.Errorf("could not determine if mlockall is enabled on node ID = %v", nodeID)
}
