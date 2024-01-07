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
	"time"
)

// Setting the environment variables
var sysdigURL = os.Getenv("API_URL")
var token = os.Getenv("SECURE_API_TOKEN")

// validate environment variables are set
func init() {
	validateEnvironment()
}

func validateEnvironment() {
	if sysdigURL == "" {
		log.Fatal("API_URL environment variable not set")
	}
	if token == "" {
		log.Fatal("SECURE_API_TOKEN environment variable not set")
	}
}

func main() {

	start := time.Now()

	args := parseArguments()

	clusters, err := getClusterData(args.Limit, args.Filter, args.Connected)
	if err != nil {
		log.Fatal("error getting cluster data: ", err)
		return
	}

	clustersWithAgentInfo, err := addAgentMetadata(clusters)
	if err != nil {
		log.Fatal("error enriching cluster data: ", err)
		return
	}

	getMetricsData(clustersWithAgentInfo)

	// Write to CSV
	err = writeToCSV(args.Output, clustersWithAgentInfo)
	if err != nil {
		fmt.Println("Failed to write to CSV:", err)
	}
	end_time := time.Now()
	fmt.Println("Execution time: ", end_time.Sub(start))
}

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

func updateClusterMetadataWithAgentData(clusterMetadata *model.ClusterWithAgentMetadata, agentData model.AgentData) {
	agentDetails := filterAgentDetails(agentData.Details, []model.AgentStatusType{
		model.AgentStatusAlmostOutOfDate,
		model.AgentStatusOutOfDate,
		model.AgentStatusUpToDate,
		model.AgentStatusDisconnected,
	})
	if agentData.AgentStats != (model.AgentStats{}) {
		clusterMetadata.NodesConnected = fmt.Sprintf("%v", agentData.AgentStats.TotalCount)
	}
	if len(agentDetails) > 0 {
		clusterMetadata.AgentStatus = agentDetails[0].AgentStatus
		clusterMetadata.AgentVersion = agentDetails[0].AgentVersion
	}
}

func addAgentMetadata(clusters []model.ClusterInfo) ([]model.ClusterWithAgentMetadata, error) {
	clusterWithAgentMetadata := make([]model.ClusterWithAgentMetadata, len(clusters))
	for i, cluster := range clusters {
		clusterMetadata := model.ClusterWithAgentMetadata{
			ClusterInfo:    cluster,
			NodesConnected: "0", // default value
			AgentStatus:    "N/A",
			AgentVersion:   "N/A",
		}

		if cluster.AgentConnected {
			agentData, err := getAgentData(cluster.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to get agent data for cluster %v. Error: %v", cluster.Name, err)
			}

			updateClusterMetadataWithAgentData(&clusterMetadata, agentData) // extracted function
		}

		clusterWithAgentMetadata[i] = clusterMetadata
	}
	return clusterWithAgentMetadata, nil
}

// func addAgentMetadata(clusters []model.ClusterInfo) ([]model.ClusterWithAgentMetadata, error) {
// 	clusterWithAgentMetadata := make([]model.ClusterWithAgentMetadata, len(clusters))
// 	var wg sync.WaitGroup
// 	errChan := make(chan error, len(clusters))

// 	for i, cluster := range clusters {
// 		wg.Add(1)
// 		go func(i int, cluster model.ClusterInfo) {
// 			defer wg.Done()

// 			clusterMetadata := model.ClusterWithAgentMetadata{
// 				ClusterInfo:    cluster,
// 				NodesConnected: "0",
// 				AgentStatus:    "N/A",
// 				AgentVersion:   "N/A",
// 			}

// 			if cluster.AgentConnected {
// 				agentData, err := getAgentData(cluster.Name)
// 				if err != nil {
// 					errChan <- fmt.Errorf("failed to get agent data for cluster %v. Error: %v", cluster.Name, err)
// 					return
// 				}

// 				updateClusterMetadataWithAgentData(&clusterMetadata, agentData)
// 			}

// 			clusterWithAgentMetadata[i] = clusterMetadata
// 		}(i, cluster)
// 	}

// 	wg.Wait()
// 	close(errChan)

// 	// Check if any errors occurred
// 	for err := range errChan {
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return clusterWithAgentMetadata, nil
// }

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
	for _, clustersWithAgentMetadata := range clusterWithAgentMetadata {
		// Use the fields from `clustersWithAgentMetadata`, e.g., clustersWithAgentMetadata.Name, clustersWithAgentMetadata.NodeCount, etc.
		writer.Write([]string{
			clustersWithAgentMetadata.Name,
			fmt.Sprintf("%d", clustersWithAgentMetadata.NodeCount),
			fmt.Sprintf("%v", clustersWithAgentMetadata.AgentConnected),
			clustersWithAgentMetadata.NodesConnected,
			clustersWithAgentMetadata.AgentStatus,
			clustersWithAgentMetadata.AgentVersion,
			clustersWithAgentMetadata.Provider,
			getEnvironment(string(clustersWithAgentMetadata.Name[3])),
		})
	}

	return nil
}

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
