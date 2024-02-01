package model

type AgentData struct {
	AgentStats AgentStats    `json:"agentStats"`
	Details    []AgentDetail `json:"details"`
}

type AgentStats struct {
	AlmostOutOfDateCount    int `json:"almostOutOfDateCount"`
	OutOfDateCount          int `json:"outOfDateCount"`
	DisconnectedCount       int `json:"disconnectedCount"`
	HealthyCount            int `json:"healthyCount"`
	NeverConnected          int `json:"neverConnected"`
	Unknown                 int `json:"unknown"`
	TotalContainerisedCount int `json:"totalContainerisedCount"`
	TotalCount              int `json:"totalCount"`
}

type AgentDetail struct {
	AgentStatus    string            `json:"agentStatus"`
	AgentVersion   string            `json:"agentVersion"`
	AgentLastSeen  string            `json:"agentLastSeen"`
	ClusterName    string            `json:"clusterName"`
	DeploymentType string            `json:"deploymentType"`
	Containerised  bool              `json:"containerised"`
	Labels         map[string]string `json:"labels"`
}
