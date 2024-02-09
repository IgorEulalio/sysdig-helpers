package main

import (
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/adapter"
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/biz"
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/client"
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/config"
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/logging"
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/model"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

var mutex sync.Mutex

func init() {
	errs := config.LoadConfig()
	logging.InitLogger(config.Config)

	if len(errs) > 0 {
		for _, err := range errs {
			logging.Log.Error(err)
		}
		os.Exit(1) // or another code to signify abnormal termination
	}
}

func main() {
	args := parseArguments()

	var sysdigClient client.API = &client.Client{
		SecureApiToken: config.Config.SecureApiToken,
		BaseURL:        config.Config.ApiURL,
		MaxRetries:     config.Config.ApiMaxRetries,
	}
	logging.Log.Debugf("Created HTTP client with following configs: %+v", sysdigClient)

	start := time.Now()

	logging.Log.Info("Getting cluster CSPM data from datasources...")
	clusters, err := sysdigClient.GetClusterData(args.Limit, args.Filter, args.Connected)
	if err != nil {
		logging.Log.Fatal("error getting cluster data: ", err)
		return
	}

	logging.Log.Info("Getting runtime information for CSPM clusters...")
	clustersWithAgentInfo, runtimeClusters, err := getExtraFeaturesInformationFromClusters(clusters, sysdigClient)
	if err != nil {
		logging.Log.Fatal("error enriching cluster data: ", err)
		return
	}

	inventoryHostsData, err := biz.GetHostsData(sysdigClient)
	if err != nil {
		logging.Log.Errorf("Error getting complete inventory: %v", err)
	}

	resourcesByHash, err := biz.GetResourcesByHash(sysdigClient, inventoryHostsData)
	if err != nil {
		logging.Log.Errorf("Error retrieving resources by hash: %v", err)
	}

	connectedAgents, err := biz.GetConnectedHosts(sysdigClient)

	connectedHostsMacAddress := biz.GetConnectedHostsMacAddress(connectedAgents)
	hosts, err := biz.MapCloudResourceToHosts(resourcesByHash, connectedHostsMacAddress)

	mergeClusterInfoWithRuntime(clustersWithAgentInfo, runtimeClusters)
	getMetricsData(clustersWithAgentInfo)

	// Write cluster data to CSV
	err = adapter.WriteClusterData(args.Output, clustersWithAgentInfo)
	if err != nil {
		logging.Log.Info("Failed to write to CSV:", err)
	}

	// Write host data to CSV
	err = adapter.WriteHostDataToCSV(args.Output, hosts)
	if err != nil {
		logging.Log.Info("Failed to write to CSV:", err)
	}

	endTime := time.Now()
	logging.Log.Info("Execution time: ", endTime.Sub(start))
}

func mergeClusterInfoWithRuntime(clusters []model.ClusterWithAgentMetadata, runtimeClusters map[string]model.RuntimeCluster) {
	for i := range clusters {
		cluster := &clusters[i] // Get a pointer to the actual element in the slice
		runtimeCluster, exist := runtimeClusters[cluster.Name]
		if exist {
			cluster.RuntimeEnabled = runtimeCluster.IsEnabled
		} else {
			cluster.RuntimeEnabled = false
		}
	}
}

func getMetricsData(clusterWithAgentMetadata []model.ClusterWithAgentMetadata) {

	totalNodesConnected := 0
	totalNodes := 0

	for _, cluster := range clusterWithAgentMetadata {
		nodesConnected, err := strconv.Atoi(cluster.NodesConnected)
		if err != nil {
			logging.Log.Fatal("error converting NodesConnected to int: ", err)
			return
		}
		totalNodesConnected += nodesConnected
		totalNodes += cluster.NodeCount
	}

	logging.Log.Info("Total Nodes Connected: ", totalNodesConnected)
	logging.Log.Info("Total Nodes: ", totalNodes)
	logging.Log.Info("Percentage of Nodes Connected: ", float64(totalNodesConnected)/float64(totalNodes)*100)
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

func getExtraFeaturesInformationFromClusters(clusters []model.ClusterInfo, sysdigClient client.API) ([]model.ClusterWithAgentMetadata, map[string]model.RuntimeCluster, error) {
	clustersWithAgentMetadata := make([]model.ClusterWithAgentMetadata, len(clusters))

	// map of cluster name to runtime data
	runtimeClusters := make(map[string]model.RuntimeCluster)
	var wg sync.WaitGroup
	errChan := make(chan error, len(clusters))

	for i, cluster := range clusters {
		wg.Add(2)
		go func(i int, cluster model.ClusterInfo) {
			defer wg.Done()

			clusterMetadata := model.ClusterWithAgentMetadata{
				ClusterInfo:    cluster,
				NodesConnected: "0",
				AgentStatus:    "N/A",
				AgentVersion:   "N/A",
				RuntimeEnabled: false,
			}

			if cluster.AgentConnected {
				agentData, err := sysdigClient.GetAgentData(cluster.Name)
				if err != nil {
					errChan <- fmt.Errorf("failed to get agent data for cluster %v. Error: %v", cluster.Name, err)
					return
				}

				updateClusterMetadataWithAgentData(&clusterMetadata, agentData)
			}

			clustersWithAgentMetadata[i] = clusterMetadata
		}(i, cluster)
		go func(i int, cluster model.ClusterInfo) {
			defer wg.Done()
			runtimeCluster, err := sysdigClient.GetRuntimeData(cluster.Name)
			if err != nil {
				errChan <- fmt.Errorf("failed to get runtime data for cluster %v: %v", cluster.Name, err)
				return
			}
			mutex.Lock()
			runtimeClusters[cluster.Name] = runtimeCluster
			mutex.Unlock()
		}(i, cluster)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, nil, err
		}
	}
	return clustersWithAgentMetadata, runtimeClusters, nil
}

func bodyReader(body io.ReadCloser) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		logging.Log.Fatal(err)
	}

	logging.Log.Info(string(bodyBytes))
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
