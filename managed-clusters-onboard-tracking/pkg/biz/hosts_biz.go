package biz

import (
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/client"
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/logging"
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/model"
	"encoding/json"
	"strings"
)

func GetHostsData(sysdigClient client.API) ([]model.Inventory, error) {
	var inventoryList []model.Inventory
	pageSize := 500
	for {
		inventory, err := sysdigClient.GetInventoryData(pageSize)
		if err != nil {
			logging.Log.Errorf("Error getting inventory data: %v", err)
		}
		inventoryList = append(inventoryList, inventory.Inventory...)
		if inventory.TotalCount < pageSize {
			break
		}
	}

	logging.Log.Infof("Total EC2 hosts found: %d", len(inventoryList))

	return inventoryList, nil
}

func GetResourcesByHash(sysdigClient client.API, inventory []model.Inventory) ([]model.CloudResource, error) {
	var cloudResources []model.CloudResource

	for _, ec2FromInventory := range inventory {
		cloudResource, err := sysdigClient.GetCloudResourceFromHash(ec2FromInventory.Hash)
		if err != nil {
			logging.Log.Errorf("Error getting cloud resource: %v. Resource Hash: %s", err, ec2FromInventory.Hash)
			continue // Proceed with the next item in the loop
		}

		// Step 1: Unmarshal json.RawMessage into a string
		var jsonString string
		err = json.Unmarshal(cloudResource.Configuration, &jsonString)
		if err != nil {
			logging.Log.Errorf("Fail to unmarshal configuration JSON string: %v", err)
			continue // Proceed with the next item in the loop
		}

		// Step 2: Unmarshal the string into the model.Configuration struct
		var config model.Configuration
		err = json.Unmarshal([]byte(jsonString), &config)
		if err != nil {
			logging.Log.Errorf("Fail to unmarshal config object into Configuration struct: %v", err)
			continue // Proceed with the next item in the loop
		}

		logging.Log.Infof("Configuration: %+v", config)
		cloudResource.ParsedConfiguration = config // Assuming you have this field in your CloudResource struct
		cloudResources = append(cloudResources, cloudResource)
	}

	return cloudResources, nil
}

func MapCloudResourceToHosts(cloudResources []model.CloudResource, macAddresses []string) ([]model.Host, error) {
	var hosts []model.Host
	for _, cloudResource := range cloudResources {
		host := model.Host{
			Name:             cloudResource.Name,
			Account:          cloudResource.Metadata.Account,
			Organization:     cloudResource.Metadata.Organization,
			Region:           cloudResource.Metadata.Region,
			ClusterName:      getClusterNameFromLabels(cloudResource.Labels),
			IsKubernetesHost: isKubernetesHost(cloudResource.Labels),
			NodeGroup:        getNodeGroupName(cloudResource.Labels),
			Connected:        verifyIfENIMacAddressIsPresent(cloudResource.ParsedConfiguration.NetworkInterfaces, macAddresses),
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

func verifyIfENIMacAddressIsPresent(interfaces []model.NetworkInterface, addresses []string) bool {
	for _, networkInterface := range interfaces {
		for _, address := range addresses {
			if networkInterface.MacAddress == address {
				return true
			}
		}
	}
	return false
}

func getNodeGroupName(labels []string) string {
	for _, label := range labels {
		if strings.Contains(label, "eks:nodegroup-name") {
			stringParts := strings.SplitAfter(label, ":")
			nodeGroupName := stringParts[len(stringParts)-1]
			return strings.TrimSpace(nodeGroupName)
		}
	}
	return "N/A"
}

func isKubernetesHost(labels []string) bool {
	for _, label := range labels {
		if strings.Contains(label, "eks:cluster-name") {
			return true
		}
	}
	return false
}

func getClusterNameFromLabels(labels []string) string {
	for _, label := range labels {
		if strings.Contains(label, "eks:cluster-name") {
			stringParts := strings.SplitAfter(label, ":")
			clusterName := stringParts[len(stringParts)-1]
			return strings.TrimSpace(clusterName)
		}
	}
	return "N/A"
}

func GetConnectedHosts(sysdigClient client.API) ([]model.AgentsConnected, error) {
	connectedAgents, err := sysdigClient.GetConnectedAgents()
	if err != nil {
		logging.Log.Errorf("Error obtaining connected agents. Error: %v", err)
		return []model.AgentsConnected{}, err
	}

	return connectedAgents.Agents, nil
}

func GetConnectedHostsMacAddress(hosts []model.AgentsConnected) []string {
	var macAddresses []string
	for _, host := range hosts {
		macAddresses = append(macAddresses, host.MachineID)
	}
	return macAddresses
}
