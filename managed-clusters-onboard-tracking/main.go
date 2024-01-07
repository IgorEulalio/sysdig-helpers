package main

import (
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/internal/model"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Setting the environment variables
var sysdigURL = os.Getenv("API_URL")
var token = os.Getenv("SECURE_API_TOKEN")

func main() {
	args := parseArguments()

	clusters, err := getClusterData(args.Limit, args.Filter, args.Connected)
	if err != nil {
		log.Fatal("error getting cluster data: ", err)
		return
	}

	clusters_with_agents_info, err := addAgentMetadata(clusters)
	if err != nil {
		log.Fatal("error enriching cluster data: ", err)
		return
	}

	getMetricsData(clusters_with_agents_info)

	// Write to CSV
	err = writeToCSV(args.Output, clusters_with_agents_info)
	if err != nil {
		fmt.Println("Failed to write to CSV:", err)
	}
}

// Function to retrieve following metrics from clusterWithAgentMetadata object
// Based on NodesConnected and ClusterInfo.NodeCount, extract total percentage of nodes connected
func getMetricsData(clusterWithAgentMetadata []model.ClusterWithAgentMetadata) {

	totalNodesConnected := 0
	totalNodes := 0

	for _, cluster := range clusterWithAgentMetadata {
		nodesConnected, err := strconv.Atoi(cluster.NodesConnected)
		if err != nil {
			log.Fatal("error converting NodesConnected to int: ", err)
			return
		}
		totalNodesConnected += nodesConnected
		totalNodes += cluster.NodeCount
	}

	fmt.Println("Total Nodes Connected: ", totalNodesConnected)
	fmt.Println("Total Nodes: ", totalNodes)
	fmt.Println("Percentage of Nodes Connected: ", float64(totalNodesConnected)/float64(totalNodes)*100)
}

func addAgentMetadata(clusters []model.ClusterInfo) ([]model.ClusterWithAgentMetadata, error) {
	clusterWithAgentMetadata := make([]model.ClusterWithAgentMetadata, len(clusters))

	// opportunity to refactor this loop to use goroutines
	for i, cluster := range clusters {
		cluster_with_agent_metadata := model.ClusterWithAgentMetadata{ClusterInfo: cluster}
		if cluster.AgentConnected {
			agentData, err := getAgentData(cluster.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to get agent data: %v", err)
			}

			if agentData.AgentStats != (model.AgentStats{}) {
				cluster_with_agent_metadata.NodesConnected = fmt.Sprintf("%v", agentData.AgentStats.TotalCount)
			} else {
				cluster_with_agent_metadata.NodesConnected = "0"
			}

			agentDetails := filterAgentDetails(agentData.Details, []model.AgentStatusType{
				model.AgentStatusAlmostOutOfDate,
				model.AgentStatusOutOfDate,
				model.AgentStatusUpToDate,
				model.AgentStatusDisconnected,
			})
			if len(agentDetails) > 0 {
				cluster_with_agent_metadata.AgentStatus = agentDetails[0].AgentStatus
				cluster_with_agent_metadata.AgentVersion = agentDetails[0].AgentVersion
			} else {
				cluster_with_agent_metadata.AgentStatus = "N/A"
				cluster_with_agent_metadata.AgentVersion = "N/A"
			}
		} else {
			cluster_with_agent_metadata.NodesConnected = "0"
			cluster_with_agent_metadata.AgentStatus = "N/A"
			cluster_with_agent_metadata.AgentVersion = "N/A"
		}

		clusterWithAgentMetadata[i] = cluster_with_agent_metadata
	}

	return clusterWithAgentMetadata, nil
}

// Function to write to CSV
func writeToCSV(fileName string, clusterWithAgentMetadata []model.ClusterWithAgentMetadata) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"name", "node_count", "agentConnected", "nodes_connected", "agent_status", "agent_version", "provider", "environment"})

	// Write data
	for _, cluster_with_agent_metadata := range clusterWithAgentMetadata {
		// Use the fields from `cluster_with_agent_metadata`, e.g., cluster_with_agent_metadata.Name, cluster_with_agent_metadata.NodeCount, etc.
		writer.Write([]string{
			cluster_with_agent_metadata.Name,
			fmt.Sprintf("%d", cluster_with_agent_metadata.NodeCount),
			fmt.Sprintf("%v", cluster_with_agent_metadata.AgentConnected),
			cluster_with_agent_metadata.NodesConnected,
			cluster_with_agent_metadata.AgentStatus,
			cluster_with_agent_metadata.AgentVersion,
			cluster_with_agent_metadata.Provider,
			getEnvironment(string(cluster_with_agent_metadata.Name[3])),
		})
	}

	return nil
}

// Function to get the cluster data with named arguments for filter, limit, and connected
func getClusterData(limit int, filter, connected string) ([]model.ClusterInfo, error) {
	url := fmt.Sprintf("%s/api/cloud/v2/dataSources/clusters?limit=%d&filter=%s&connected=%s", sysdigURL, limit, filter, connected)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

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

// Function to get agent data
func getAgentData(clusterName string) (model.AgentData, error) {
	url := fmt.Sprintf("%s/api/cloud/v2/dataSources/agents?filter=%s", sysdigURL, clusterName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return model.AgentData{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.AgentData{}, fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

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

func filterAgentDetails(details []model.AgentDetail, statuses []model.AgentStatusType) []model.AgentDetail {
	var filteredDetails []model.AgentDetail

	statusMap := make(map[model.AgentStatusType]bool)
	for _, status := range statuses {
		statusMap[status] = true
	}

	for _, detail := range details {
		if _, ok := statusMap[model.AgentStatusType(detail.AgentStatus)]; ok {
			filteredDetails = append(filteredDetails, detail)
		}
	}

	return filteredDetails
}

// Function to return environment based on environment single string
func getEnvironment(environment string) string {
	switch environment {
	case "d":
		return "development"
	case "p":
		return "production"
	case "i":
		return "pre-production"
	default:
		return "unknown"
	}
}

type CommandLineArgs struct {
	Limit     int
	Filter    string
	Connected string
	Output    string
}

func customUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Description:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "\tcommand-line tool for tracking cluster onboarding based on CSPM data. It provides functionalities like fetching cluster data, filtering based on specific criteria, and exporting details to a CSV file.\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Requirements:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "\tSet SYSDIG_TOKEN with your Secure API token from Sysdig UI.\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "\tSet SYSDIG_URL with your Sysdig API endpoint URL.\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func parseArguments() CommandLineArgs {
	flag.Usage = customUsage

	limit := flag.Int("limit", 150, "Limit the number of results")
	filter := flag.String("filter", "", "Filter criteria")
	connected := flag.String("connected", "", "Connected status filter")
	output := flag.String("output", "clusters.csv", "Output file name")
	flag.Parse()

	return CommandLineArgs{
		Limit:     *limit,
		Filter:    *filter,
		Connected: *connected,
		Output:    *output,
	}
}
