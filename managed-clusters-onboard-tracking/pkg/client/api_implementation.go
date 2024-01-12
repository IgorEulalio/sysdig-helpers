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

	// Construct path parameters
	pathParams := map[string]string{"filter": clusterName}

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
