package model

type EnrichedClusterInfo struct {
	ClusterInfo
	NodesConnected string
	AgentStatus    string
	AgentVersion   string
}
