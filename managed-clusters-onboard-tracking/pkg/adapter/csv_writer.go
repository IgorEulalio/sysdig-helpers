package adapter

import (
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/logging"
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/model"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

const dateFormat = "02-01-2006"

func WriteClusterData(fileName string, clusterWithAgentMetadata []model.ClusterWithAgentMetadata) error {

	file, err := os.Create(fmt.Sprintf("clusters-%s-%s", fileName, time.Now().Format(dateFormat)))
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
	dir, err := os.Getwd()
	if err != nil {
		logging.Log.Errorf("Failed to obtain current dir. Error %s", err)
	}
	logging.Log.Debugf("Created csv file successfully. File name: %s, file path: %s", fileName, dir)

	return nil
}

func WriteHostDataToCSV(fileName string, hosts []model.Host) error {
	time.Now()
	file, err := os.Create(fmt.Sprintf("host-%s-%s", fileName, time.Now().Format(dateFormat)))

	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"name", "account", "organization", "region", "is_kubernetes_host", "cluster_name", "nodegroup_name", "connected"})

	// Write data
	for _, host := range hosts {
		// Use the fields from `clustersWithAgentMetadata`, e.g., clustersWithAgentMetadata.Name, clustersWithAgentMetadata.NodeCount, etc.
		writer.Write([]string{
			host.Name,
			host.Account,
			host.Organization,
			host.Region,
			strconv.FormatBool(host.IsKubernetesHost),
			host.ClusterName,
			host.NodeGroup,
			strconv.FormatBool(host.Connected),
		})
	}
	dir, err := os.Getwd()
	if err != nil {
		logging.Log.Errorf("Failed to obtain current dir. Error %s", err)
	}
	logging.Log.Debugf("Created csv file successfully. File name: %s, file path: %s", fileName, dir)

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
