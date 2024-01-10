package model

type Page struct {
	Returned int         `json:"returned"`
	Matched  int         `json:"matched"`
	Next     interface{} `json:"next"`
}

type Labels struct {
	AssetType                  string `json:"asset.type"`
	KubernetesClusterName      string `json:"kubernetes.cluster.name"`
	KubernetesNamespaceName    string `json:"kubernetes.namespace.name"`
	KubernetesPodContainerName string `json:"kubernetes.pod.container.name"`
	KubernetesWorkloadName     string `json:"kubernetes.workload.name"`
	KubernetesWorkloadType     string `json:"kubernetes.workload.type"`
}

type RecordDetails struct {
	MainAssetName string `json:"mainAssetName"`
	Labels        Labels `json:"labels"`
}

type RuntimeResult struct {
	HashId                  string        `json:"hashId"`
	ResultId                string        `json:"resultId"`
	RecordDetails           RecordDetails `json:"recordDetails"`
	VulnsBySev              []int         `json:"vulnsBySev"`
	RunningVulnsBySev       []int         `json:"runningVulnsBySev"`
	ExploitCount            int           `json:"exploitCount"`
	IsEVEEnabled            bool          `json:"isEVEEnabled"`
	PolicyEvaluationsResult string        `json:"policyEvaluationsResult"`
	HasAcceptedRisk         bool          `json:"hasAcceptedRisk"`
}

type RuntimeData struct {
	Page Page            `json:"page"`
	Data []RuntimeResult `json:"data"`
}

type RuntimeCluster struct {
	ClusterName string
	IsEnabled   bool
}
