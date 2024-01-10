package model

type ClusterWithAgentMetadata struct {
	ClusterInfo
	NodesConnected string
	AgentStatus    string
	AgentVersion   string
	RuntimeEnabled bool
}
