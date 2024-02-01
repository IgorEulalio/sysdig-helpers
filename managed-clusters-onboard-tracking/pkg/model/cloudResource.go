package model

import (
	"encoding/json"
	"time"
)

type CloudResourceWrapper struct {
	CloudResource CloudResource `json:"data"`
}

type CloudResource struct {
	ID                    string               `json:"id"`
	Name                  string               `json:"name"`
	Platform              string               `json:"platform"`
	Type                  string               `json:"type"`
	LastSeen              string               `json:"lastSeen"`
	Metadata              Metadata             `json:"metadata"`
	PosturePolicySummary  PosturePolicySummary `json:"posturePolicySummary"`
	Configuration         json.RawMessage      `json:"configuration"` // RawMessage used for arbitrary JSON
	ParsedConfiguration   Configuration
	KeyValueConfigs       []interface{}           `json:"keyValueConfigs"` // Empty array, so type is ambiguous
	Labels                []string                `json:"labels"`
	PostureControlSummary []PostureControlSummary `json:"postureControlSummary"`
	Zones                 []Zone                  `json:"zones"`
	Hash                  string                  `json:"hash"`
	ResourceOrigin        string                  `json:"resourceOrigin"`
	ClusterName           string                  `json:"clusterName"`
	Namespace             string                  `json:"namespace"`
	SourceHash            string                  `json:"sourceHash"`
	Category              string                  `json:"category"`
	SourceType            int                     `json:"sourceType"`
	SourcePath            string                  `json:"sourcePath"`
	RepositoryName        string                  `json:"repositoryName"`
	SourceURL             string                  `json:"sourceURL"`
	GitIntegrationId      string                  `json:"gitIntegrationId"`
	GitIntegrationName    string                  `json:"gitIntegrationName"`
	Kind                  string                  `json:"kind"`
}

type InventoryWrapper struct {
	Inventory  []Inventory `json:"data"`
	TotalCount int         `json:"totalCount"`
}

type Inventory struct {
	ID                        string               `json:"id"`
	Name                      string               `json:"name"`
	Platform                  string               `json:"platform"`
	Type                      string               `json:"type"`
	LastSeen                  string               `json:"lastSeen"`
	Metadata                  Metadata             `json:"metadata"`
	PosturePolicySummary      PosturePolicySummary `json:"posturePolicySummary"`
	Configuration             string               `json:"configuration"`
	KeyValueConfigs           []interface{}        `json:"keyValueConfigs"`
	Labels                    []string             `json:"labels"`
	PostureControlSummary     []PostureControl     `json:"postureControlSummary"`
	Zones                     []Zone               `json:"zones"`
	Hash                      string               `json:"hash"`
	ResourceOrigin            string               `json:"resourceOrigin"`
	ConfigApiEndpoint         string               `json:"configApiEndpoint"`
	VulnerabilitySummary      VulnerabilitySummary `json:"vulnerabilitySummary"`
	Tags                      []interface{}        `json:"tags"`
	Category                  string               `json:"category"`
	InUseVulnerabilitySummary VulnerabilitySummary `json:"inUseVulnerabilitySummary"`
	IsExposed                 bool                 `json:"isExposed"`
}

type PostureControl struct {
	// Define fields as per your JSON structure
}

type VulnerabilitySummary struct {
	CriticalSeverityCount   int  `json:"criticalSeverityCount"`
	HighSeverityCount       int  `json:"highSeverityCount"`
	MediumSeverityCount     int  `json:"mediumSeverityCount"`
	HasExploit              bool `json:"hasExploit"`
	LowSeverityCount        int  `json:"lowSeverityCount"`
	NegligibleSeverityCount int  `json:"negligibleSeverityCount"`
}

type Metadata struct {
	Account      string `json:"Account"`
	Organization string `json:"Organization"`
	Region       string `json:"Region"`
}

type PosturePolicySummary struct {
	PassPercentage int      `json:"passPercentage"`
	Policies       []Policy `json:"policies"`
}

type Policy struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Pass bool   `json:"pass"`
}

type PostureControlSummary struct {
	Name             string `json:"name"`
	PolicyId         string `json:"policyId"`
	FailedControls   int    `json:"failedControls"`
	AcceptedControls int    `json:"acceptedControls"`
}

type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Configuration struct {
	AmiLaunchIndex                          int                  `json:"AmiLaunchIndex"`
	Architecture                            string               `json:"Architecture"`
	BlockDeviceMappings                     []BlockDeviceMapping `json:"BlockDeviceMappings"`
	BootMode                                string               `json:"BootMode"`
	CapacityReservationId                   *string              `json:"CapacityReservationId"`
	ClientToken                             string               `json:"ClientToken"`
	CurrentInstanceBootMode                 string               `json:"CurrentInstanceBootMode"`
	EbsOptimized                            bool                 `json:"EbsOptimized"`
	ElasticGpuAssociations                  []interface{}        `json:"ElasticGpuAssociations"`
	ElasticInferenceAcceleratorAssociations []interface{}        `json:"ElasticInferenceAcceleratorAssociations"`
	EnaSupport                              bool                 `json:"EnaSupport"`
	Hypervisor                              string               `json:"Hypervisor"`
	ImageId                                 string               `json:"ImageId"`
	InstanceId                              string               `json:"InstanceId"`
	InstanceLifecycle                       string               `json:"InstanceLifecycle"`
	InstanceType                            string               `json:"InstanceType"`
	Ipv6Address                             *string              `json:"Ipv6Address"`
	KernelId                                *string              `json:"KernelId"`
	KeyName                                 *string              `json:"KeyName"`
	LaunchTime                              time.Time            `json:"LaunchTime"`
	Licenses                                []interface{}        `json:"Licenses"`
	NetworkInterfaces                       []NetworkInterface   `json:"NetworkInterfaces"`
	OutpostArn                              *string              `json:"OutpostArn"`
	Platform                                string               `json:"Platform"`
	PlatformDetails                         string               `json:"PlatformDetails"`
	PrivateDnsName                          string               `json:"PrivateDnsName"`
	PrivateIpAddress                        string               `json:"PrivateIpAddress"`
	ProductCodes                            []interface{}        `json:"ProductCodes"`
	PublicDnsName                           string               `json:"PublicDnsName"`
	PublicIpAddress                         *string              `json:"PublicIpAddress"`
	RamdiskId                               *string              `json:"RamdiskId"`
	RootDeviceName                          string               `json:"RootDeviceName"`
	RootDeviceType                          string               `json:"RootDeviceType"`
	SourceDestCheck                         bool                 `json:"SourceDestCheck"`
	SpotInstanceRequestId                   *string              `json:"SpotInstanceRequestId"`
	SriovNetSupport                         *string              `json:"SriovNetSupport"`
	StateTransitionReason                   string               `json:"StateTransitionReason"`
	SubnetId                                string               `json:"SubnetId"`
	TpmSupport                              *string              `json:"TpmSupport"`
	UsageOperation                          string               `json:"UsageOperation"`
	UsageOperationUpdateTime                time.Time            `json:"UsageOperationUpdateTime"`
	VirtualizationType                      string               `json:"VirtualizationType"`
	VpcId                                   string               `json:"VpcId"`
}

type BlockDeviceMapping struct {
	DeviceName string `json:"DeviceName"`
	Ebs        Ebs    `json:"Ebs"`
}

type Ebs struct {
	AttachTime          time.Time `json:"AttachTime"`
	DeleteOnTermination bool      `json:"DeleteOnTermination"`
	Status              string    `json:"Status"`
	VolumeId            string    `json:"VolumeId"`
}

type NetworkInterface struct {
	ConnectionTrackingConfiguration *interface{}  `json:"ConnectionTrackingConfiguration"`
	Description                     string        `json:"Description"`
	InterfaceType                   string        `json:"InterfaceType"`
	Ipv4Prefixes                    *interface{}  `json:"Ipv4Prefixes"`
	Ipv6Addresses                   []interface{} `json:"Ipv6Addresses"`
	Ipv6Prefixes                    *interface{}  `json:"Ipv6Prefixes"`
	MacAddress                      string        `json:"MacAddress"`
	NetworkInterfaceId              string        `json:"NetworkInterfaceId"`
	OwnerId                         string        `json:"OwnerId"`
	PrivateDnsName                  string        `json:"PrivateDnsName"`
	PrivateIpAddress                string        `json:"PrivateIpAddress"`
	SourceDestCheck                 bool          `json:"SourceDestCheck"`
	Status                          string        `json:"Status"`
	SubnetId                        string        `json:"SubnetId"`
	VpcId                           string        `json:"VpcId"`
}
