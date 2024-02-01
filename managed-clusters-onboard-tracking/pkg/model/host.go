package model

type Host struct {
	Name             string
	Account          string
	Organization     string
	Region           string
	IsKubernetesHost bool
	ClusterName      string
	NodeGroup        string
	Connected        bool
}
