package adapter

import (
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/model"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func WriteToCSV(fileName string, clusterWithAgentMetadata []model.ClusterWithAgentMetadata) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"name", "node_count", "agentConnected", "nodes_connected", "agent_status", "agent_version", "provider", "environment", "runtime_enabled"})

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
			strconv.FormatBool(clustersWithAgentMetadata.RuntimeEnabled),
		})
	}

	return nil
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
