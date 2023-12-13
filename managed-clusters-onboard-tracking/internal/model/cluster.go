package model

type ClusterInfo struct {
	CustomerID           int    `json:"customerID"`
	AccountID            string `json:"accountID"`
	Provider             string `json:"provider"`
	Name                 string `json:"name"`
	Region               string `json:"region"`
	AgentConnected       bool   `json:"agentConnected"`
	CreatedAt            string `json:"createdAt"`
	NodeCount            int    `json:"nodeCount"`
	ClusterResourceGroup string `json:"clusterResourceGroup"`
	Version              string `json:"version"`
	AgentConnectString   string `json:"agentConnectString"`
}
