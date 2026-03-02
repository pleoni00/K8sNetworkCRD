package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NetworkGraphList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []NetworkGraph `json:"items"`
}

type NetworkGraph struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkGraphSpec   `json:"spec,omitempty"`
	Status NetworkGraphStatus `json:"status,omitempty"`
}

type NetworkGraphSpec struct {
	Nodes    []Node                `json:"nodes,omitempty"`
	Edges    []Edge                `json:"edges,omitempty"`
	Metadata *NetworkGraphMetadata `json:"metadata,omitempty"`
}

type NetworkGraphStatus struct {
	Phase       string       `json:"phase,omitempty"`
	LastUpdated *metav1.Time `json:"last_updated,omitempty"`
	NodeCount   int          `json:"node_count,omitempty"`
	EdgeCount   int          `json:"edge_count,omitempty"`
	Health      *HealthInfo  `json:"health,omitempty"`
}

type Node struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Type               string `json:"type"`
	IP                 string `json:"ip,omitempty"`
	State              string `json:"state,omitempty"`
	CPUUsageMillicores int64  `json:"cpu_usage_millicores,omitempty"`
	MemoryBytes        int64  `json:"memory_bytes,omitempty"`
}

type Edge struct {
	ID        string       `json:"id"`
	Source    string       `json:"source"`
	Target    string       `json:"target"`
	Traffic   TrafficInfo  `json:"traffic"`
	Timestamp *metav1.Time `json:"timestamp,omitempty"`
	Quality   *QualityInfo `json:"quality,omitempty"`
}

type TrafficInfo struct {
	Protocol    string `json:"protocol,omitempty"`
	Port        int    `json:"port,omitempty"`
	Bytes       int64  `json:"bytes,omitempty"`
	Packets     int64  `json:"packets,omitempty"`
	Connections int    `json:"connections,omitempty"`
	Latency     string `json:"latency,omitempty"`
	Bandwidth   string `json:"bandwidth,omitempty"`
	Direction   string `json:"direction,omitempty"`
}

type QualityInfo struct {
	SuccessRate float64 `json:"success_rate,omitempty"`
	ErrorRate   float64 `json:"error_rate,omitempty"`
	Jitter      string  `json:"jitter,omitempty"`
}

type NetworkGraphMetadata struct {
	CollectionTime     *metav1.Time `json:"collection_time,omitempty"`
	CollectionInterval string       `json:"collection_interval,omitempty"`
	DataSource         string       `json:"data_source,omitempty"`
	Version            string       `json:"version,omitempty"`
}

type HealthInfo struct {
	HealthyNodes       int    `json:"healthy_nodes,omitempty"`
	UnhealthyNodes     int    `json:"unhealthy_nodes,omitempty"`
	TotalTrafficVolume string `json:"total_traffic_volume,omitempty"`
}
