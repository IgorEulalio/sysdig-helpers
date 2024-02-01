package client

import (
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/logging"
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/model"
	"fmt"
	"net/http"
	"strconv"
)

const (
	runtimeInformationPath = "%s/api/scanning/runtime/v2/workflows/results"
	clusterInformationPath = "%s/api/cloud/v2/dataSources/clusters"
	agentInformationPath   = "%s/api/cloud/v2/dataSources/agents"
	inventoryPath          = "%s/api/cspm/v1/inventory/resources"
	resourcePath           = "%s/api/cspm/v1/cloud/resource"
	connectedAgentsPath    = "%s/api/agents/connected"
)

func (c *Client) GetRuntimeData(clusterName string) (model.RuntimeCluster, error) {

	c.ServiceName = "RUNTIME_DATA"

	filter := fmt.Sprintf("kubernetes.cluster.name = \"%s\"", clusterName)
	pathParams := map[string]string{"filter": filter}
	urlFormat, err := c.CreateUrl(runtimeInformationPath, pathParams)

	if err != nil {
		logging.Log.Errorf("Error creating URL for service %s. Error: %v", c.ServiceName, err)
		return model.RuntimeCluster{}, err
	}

	req, err := c.NewRequest(http.MethodGet, urlFormat, nil)
	if err != nil {
		logging.Log.Errorf("Error creating request for service %s. Error: %v", c.ServiceName, err)
		return model.RuntimeCluster{}, err
	}

	var runtimeData model.RuntimeData

	err = c.Do(req, &runtimeData)
	if err != nil {
		return model.RuntimeCluster{}, err
	}
	runtimeCluster := model.RuntimeCluster{
		ClusterName: clusterName,
		IsEnabled:   isRuntimeEnabled(runtimeData),
	}
	return runtimeCluster, nil
}

func (c *Client) GetClusterData(limit int, filter, connected string) ([]model.ClusterInfo, error) {

	c.ServiceName = "CLUSTER_DATA" // Set the service name for logging

	// Construct path parameters
	pathParams := map[string]string{
		"limit":     strconv.Itoa(limit),
		"filter":    filter,
		"connected": connected,
	}

	// Create URL using the client's method
	urlFormat, err := c.CreateUrl(clusterInformationPath, pathParams)
	if err != nil {
		logging.Log.Errorf("Error creating URL for service %s. Error: %v", c.ServiceName, err)
		return nil, err
	}

	// Create a new request
	req, err := c.NewRequest(http.MethodGet, urlFormat, nil)
	if err != nil {
		logging.Log.Errorf("Error creating request for service %s. Error: %v", c.ServiceName, err)
		return nil, err
	}

	// Declare a slice to hold the cluster data
	var clusters []model.ClusterInfo

	// Make the request and decode response into clusters
	err = c.Do(req, &clusters)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func (c *Client) GetAgentData(clusterName string) (model.AgentData, error) {

	c.ServiceName = "AGENT_DATA" // Set the service name for logging

	pathParams := map[string]string{"filter": clusterName,
		// If cluster has multiple nodes, limit 1 might not bring clusters with status we're looking for
		// such as Up to date, Almost out of date, Out of date
		"limit":  "500",
		"offset": "0",
	}

	// Create URL using the client's method
	urlFormat, err := c.CreateUrl(agentInformationPath, pathParams)
	if err != nil {
		logging.Log.Errorf("Error creating URL for service %s. Error: %v", c.ServiceName, err)
		return model.AgentData{}, err
	}

	// Create a new request
	req, err := c.NewRequest(http.MethodGet, urlFormat, nil)
	if err != nil {
		logging.Log.Errorf("Error creating request for service %s. Error: %v", c.ServiceName, err)
		return model.AgentData{}, err
	}

	// Declare a variable to hold the agent data
	var agentData model.AgentData

	// Make the request and decode the response into agentData
	err = c.Do(req, &agentData)
	if err != nil {
		return model.AgentData{}, err
	}

	return agentData, nil
}

func isRuntimeEnabled(runtimeData model.RuntimeData) bool {
	isEnabled := runtimeData.Page.Matched > 0

	return isEnabled
}

func (c *Client) GetInventoryData(pageSize int) (model.InventoryWrapper, error) {

	c.ServiceName = "INVENTORY_SERVICE" // Set the service name for logging
	pageNumber := 1

	pathParams := map[string]string{
		"pageNumber": strconv.Itoa(pageNumber),
		"pageSize":   strconv.Itoa(pageSize),
		"filter":     "type = \"EC2 Instance\"",
		"fields":     "id,hash,name",
	}

	// Create URL using the client's method
	urlFormat, err := c.CreateUrl(inventoryPath, pathParams)
	if err != nil {
		logging.Log.Errorf("Error creating URL for service %s. Error: %v", c.ServiceName, err)
		return model.InventoryWrapper{}, err
	}

	// Create a new request
	req, err := c.NewRequest(http.MethodGet, urlFormat, nil)
	if err != nil {
		logging.Log.Errorf("Error creating request for service %s. Error: %v", c.ServiceName, err)
		return model.InventoryWrapper{}, err
	}

	// Declare a slice to hold the cluster data
	var inventoryResponse model.InventoryWrapper

	// Make the request and decode response into clusters
	err = c.Do(req, &inventoryResponse)
	if err != nil {
		return model.InventoryWrapper{}, err
	}

	return inventoryResponse, nil
}

func (c *Client) GetCloudResourceFromHash(hash string) (model.CloudResource, error) {
	c.ServiceName = "RESOURCE_SERVICE" // Set the service name for logging

	pathParams := map[string]string{
		"resourceHash": hash,
		"fields":       "id,hash,name,platform,type,configuration,keyvalueconfigs,labels,lastseen,metadata,zones,posturepolicysummary,posturecontrolsummary,resourceorigin,category",
	}

	// Create URL using the client's method
	urlFormat, err := c.CreateUrl(resourcePath, pathParams)
	if err != nil {
		logging.Log.Errorf("Error creating URL for service %s. Error: %v", c.ServiceName, err)
		return model.CloudResource{}, err
	}

	// Create a new request
	req, err := c.NewRequest(http.MethodGet, urlFormat, nil)
	if err != nil {
		logging.Log.Errorf("Error creating request for service %s. Error: %v", c.ServiceName, err)
		return model.CloudResource{}, err
	}

	// Declare a slice to hold the cluster data
	var cloudResourceWrapper model.CloudResourceWrapper

	// Make the request and decode response into clusters
	err = c.Do(req, &cloudResourceWrapper)
	if err != nil {
		return model.CloudResource{}, err
	}

	return cloudResourceWrapper.CloudResource, nil
}

func (c *Client) GetConnectedAgents() (model.AgentsConnectedWrapper, error) {
	c.ServiceName = "CONNECTED_AGENTS" // Set the service name for logging

	// Create URL using the client's method
	pathParams := map[string]string{}
	urlFormat, err := c.CreateUrl(connectedAgentsPath, pathParams)
	if err != nil {
		logging.Log.Errorf("Error creating URL for service %s. Error: %v", c.ServiceName, err)
		return model.AgentsConnectedWrapper{}, err
	}

	// Create a new request
	req, err := c.NewRequest(http.MethodGet, urlFormat, nil)
	if err != nil {
		logging.Log.Errorf("Error creating request for service %s. Error: %v", c.ServiceName, err)
		return model.AgentsConnectedWrapper{}, err
	}

	// Declare a slice to hold the cluster data
	var agentsConnectedWrapper model.AgentsConnectedWrapper

	// Make the request and decode response into clusters
	err = c.Do(req, &agentsConnectedWrapper)
	if err != nil {
		return model.AgentsConnectedWrapper{}, err
	}

	return agentsConnectedWrapper, nil
}
