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
)

// Setting the environment variables
var sysdigURL = os.Getenv("SYSDIG_URL")
var token = os.Getenv("SYSDIG_TOKEN")

func main() {
	args := parseArguments()

	clusters, err := getClusterData(args.Limit, args.Filter, args.Connected)
	if err != nil {
		log.Fatal("error getting cluster data: ", err)
		return
	}

	fmt.Println("Number of clusters:", len(clusters))
	// Write to CSV
	err = writeToCSV(args.Output, clusters)
	if err != nil {
		fmt.Println("Failed to write to CSV:", err)
	}
}

// Function to write to CSV
func writeToCSV(fileName string, clusters []model.ClusterInfo) error {
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
	for _, cluster := range clusters {
		name := cluster.Name
		nodeCount := fmt.Sprintf("%d", cluster.NodeCount)
		agentConnected := fmt.Sprintf("%v", cluster.AgentConnected)
		provider := cluster.Provider

		environment := getEnvironment(string(name[3]))

		var nodesConnected, agentStatus, agentVersion string

		if agentConnected == "true" {
			agentData, _ := getAgentData(name)
			if agentData.AgentStats != (model.AgentStats{}) {
				agentStats := agentData.AgentStats
				nodesConnected = fmt.Sprintf("%v", agentStats.TotalCount)
			} else {
				nodesConnected = "0"
			}

			agentDetails := agentData.Details

			filtered := filterAgentDetails(agentDetails, []string{"Almost out of date", "Out of date", "Up to date"})
			if len(filtered) > 0 {
				agentStatus = fmt.Sprintf("%v", filtered[0].AgentStatus)
				agentVersion = fmt.Sprintf("%v", filtered[0].AgentVersion)
			} else {
				agentStatus = "N/A"
				agentVersion = "N/A"
			}
		} else {
			nodesConnected = "0"
			agentStatus = "N/A"
			agentVersion = "N/A"
		}

		writer.Write([]string{name, nodeCount, agentConnected, nodesConnected, agentStatus, agentVersion, provider, environment})
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

// Function to filter agent details based on status
func filterAgentDetails(details []model.AgentDetail, statuses []string) []model.AgentDetail {
	var filteredDetails []model.AgentDetail

	statusMap := make(map[string]bool)
	for _, status := range statuses {
		statusMap[status] = true
	}

	for _, detail := range details {
		if _, ok := statusMap[detail.AgentStatus]; ok {
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

func parseArguments() CommandLineArgs {
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
