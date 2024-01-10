package client

import (
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/config"
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/logging"
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func GetRuntimeData(clusterName string) (model.RuntimeCluster, error) {
	filter := fmt.Sprintf("kubernetes.cluster.name = \"%s\"", clusterName)
	urlFormat := fmt.Sprintf("%s/api/scanning/runtime/v2/workflows/results?filter=%s", config.Config.ApiURL, filter)
	parsedUrl, err := url.Parse(urlFormat)
	if err != nil {
		logging.Log.Errorf("Error parsing URL for runtime endpoint. Cluster name %s, Request URL %s", clusterName, urlFormat)
	}
	q := parsedUrl.Query()
	q.Set("filter", fmt.Sprintf("kubernetes.cluster.name = \"%s\"", clusterName))
	q.Set("limit", strconv.Itoa(1))
	parsedUrl.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", parsedUrl.String(), nil)
	if err != nil {
		return model.RuntimeCluster{}, fmt.Errorf("failed to create request for runtime request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.Config.SecureApiToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.RuntimeCluster{}, fmt.Errorf("failed to perform request for runtime request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.Log.Errorf("Error closing body: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return model.RuntimeCluster{}, fmt.Errorf("failed to runtime data for cluster %s. Status code: %d. Request: %s", clusterName, resp.StatusCode, req)
	}
	var runtimeData model.RuntimeData

	err = json.NewDecoder(resp.Body).Decode(&runtimeData)
	if err != nil {
		return model.RuntimeCluster{}, fmt.Errorf("failed to decode response from runtime endpoint: %v", err)
	}
	runtimeCluster := model.RuntimeCluster{
		ClusterName: clusterName,
		IsEnabled:   isRuntimeEnabled(runtimeData),
	}
	return runtimeCluster, nil
}

func GetClusterData(limit int, filter, connected string) ([]model.ClusterInfo, error) {
	urlFormat := fmt.Sprintf("%s/api/cloud/v2/dataSources/clusters?limit=%d&filter=%s&connected=%s", config.Config.ApiURL, limit, filter, connected)
	req, err := http.NewRequest("GET", urlFormat, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.Config.SecureApiToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.Log.Errorf("Error closing body: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch cluster data. Status code: %d", resp.StatusCode)
	}

	var clusters []model.ClusterInfo
	err = json.NewDecoder(resp.Body).Decode(&clusters)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return clusters, nil
}

func GetAgentData(clusterName string) (model.AgentData, error) {
	urlFormat := fmt.Sprintf("%s/api/cloud/v2/dataSources/agents?filter=%s", config.Config.ApiURL, clusterName)
	req, err := http.NewRequest("GET", urlFormat, nil)
	if err != nil {
		return model.AgentData{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.Config.SecureApiToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.AgentData{}, fmt.Errorf("failed to perform request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.Log.Errorf("Error closing body: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return model.AgentData{}, fmt.Errorf("failed to fetch agent data. Status code: %d", resp.StatusCode)
	}

	var agentData model.AgentData
	err = json.NewDecoder(resp.Body).Decode(&agentData)
	if err != nil {
		return model.AgentData{}, fmt.Errorf("failed to decode response: %v", err)
	}

	return agentData, nil
}

func isRuntimeEnabled(runtimeData model.RuntimeData) bool {
	isEnabled := runtimeData.Page.Matched > 0

	return isEnabled
}
