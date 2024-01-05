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

	enrichedClusters, err := enrichClusterData(clusters)
	if err != nil {
		log.Fatal("error enriching cluster data: ", err)
		return
	}

	getMetricsData(enrichedClusters)

	// Write to CSV
	err = writeToCSV(args.Output, enrichedClusters)
	if err != nil {
		fmt.Println("Failed to write to CSV:", err)
	}
}

// Function to retrieve following metrics from enrichedClusters object
// Based on NodesConnected and ClusterInfo.NodeCount, extract total percentage of nodes connected
func getMetricsData(enrichedClusters []model.EnrichedClusterInfo) {

	totalNodesConnected := 0
	totalNodes := 0

	for _, cluster := range enrichedClusters {
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

func enrichClusterData(clusters []model.ClusterInfo) ([]model.EnrichedClusterInfo, error) {
	enrichedClusters := make([]model.EnrichedClusterInfo, len(clusters))

	for i, cluster := range clusters {
		enriched := model.EnrichedClusterInfo{ClusterInfo: cluster}
		if cluster.AgentConnected {
			agentData, err := getAgentData(cluster.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to get agent data: %v", err)
			}

			if agentData.AgentStats != (model.AgentStats{}) {
				enriched.NodesConnected = fmt.Sprintf("%v", agentData.AgentStats.TotalCount)
			} else {
				enriched.NodesConnected = "0"
			}

			agentDetails := filterAgentDetails(agentData.Details, []model.AgentStatusType{
				model.AgentStatusAlmostOutOfDate,
				model.AgentStatusOutOfDate,
				model.AgentStatusUpToDate,
				model.AgentStatusDisconnected,
			})
			if len(agentDetails) > 0 {
				enriched.AgentStatus = agentDetails[0].AgentStatus
				enriched.AgentVersion = agentDetails[0].AgentVersion
			} else {
				enriched.AgentStatus = "N/A"
				enriched.AgentVersion = "N/A"
			}
		} else {
			enriched.NodesConnected = "0"
			enriched.AgentStatus = "N/A"
			enriched.AgentVersion = "N/A"
		}

		enrichedClusters[i] = enriched
	}

	return enrichedClusters, nil
}

// Function to write to CSV
func writeToCSV(fileName string, enrichedClusters []model.EnrichedClusterInfo) error {
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
	for _, enriched := range enrichedClusters {
		// Use the fields from `enriched`, e.g., enriched.Name, enriched.NodeCount, etc.
		writer.Write([]string{
			enriched.Name,
			fmt.Sprintf("%d", enriched.NodeCount),
			fmt.Sprintf("%v", enriched.AgentConnected),
			enriched.NodesConnected,
			enriched.AgentStatus,
			enriched.AgentVersion,
			enriched.Provider,
			getEnvironment(string(enriched.Name[3])),
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
