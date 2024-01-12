package client

import "IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/model"

type API interface {
	GetRuntimeData(clusterName string) (model.RuntimeCluster, error)
	GetClusterData(limit int, filter, connected string) ([]model.ClusterInfo, error)
	GetAgentData(clusterName string) (model.AgentData, error)
}
