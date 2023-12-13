package model

// AgentStatusType defines a type for agent status enum.
type AgentStatusType string

// Define enum values for AgentStatusType.
const (
	AgentStatusAlmostOutOfDate AgentStatusType = "Almost Out of Date"
	AgentStatusOutOfDate       AgentStatusType = "Out of Date"
	AgentStatusUpToDate        AgentStatusType = "Up to Date"
	AgentStatusDisconnected    AgentStatusType = "Disconnected"
)
