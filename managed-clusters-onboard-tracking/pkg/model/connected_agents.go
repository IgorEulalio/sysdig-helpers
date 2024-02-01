package model

type AgentsConnectedWrapper struct {
	Total  int               `json:"total"`
	Agents []AgentsConnected `json:"agents"`
}

type AgentsConnected struct {
	Customer   int        `json:"customer"`
	ID         int64      `json:"id"`
	MachineID  string     `json:"machineId"`
	HostName   string     `json:"hostName"`
	Attributes Attributes `json:"attributes"`
	Connected  bool       `json:"connected"`
	Status     string     `json:"status"`
	TeamName   string     `json:"teamName"`
}

type Attributes struct {
	CustomName                     string            `json:"customName"`
	Tags                           map[string]string `json:"tags"`
	CustomMap                      *interface{}      `json:"customMap"` // Assuming the type is uncertain
	Hidden                         bool              `json:"hidden"`
	HiddenProcesses                string            `json:"hiddenProcesses"`
	HostName                       string            `json:"hostName"`
	ClusterId                      string            `json:"clusterId"`
	ClusterName                    string            `json:"clusterName"`
	Delegated                      bool              `json:"delegated"`
	Version                        string            `json:"version"`
	AgentType                      string            `json:"agentType"`
	IPAddresses                    []string          `json:"ipAddresses"`
	DateCreated                    int64             `json:"dateCreated"`
	Ts                             int64             `json:"ts"`
	SwarmNodeId                    *interface{}      `json:"swarmNodeId"` // Assuming the type is uncertain
	ProtocolVersion                int               `json:"protocolVersion"`
	AggregationInterval            int               `json:"aggregationInterval"`
	MaxSupportedCustomMetricsLimit int               `json:"maxSupportedCustomMetricsLimit"`
	FastProtoSupported             bool              `json:"fastProtoSupported"`
	Serverless                     bool              `json:"serverless"`
	PromConfigSupported            bool              `json:"promConfigSupported"`
	PromscrapeVersion              int               `json:"promscrapeVersion"`
	LiveLogsSupported              bool              `json:"liveLogsSupported"`
	K8sCommandsSupported           bool              `json:"k8sCommandsSupported"`
	OrchestratorType               string            `json:"orchestratorType"`
}
