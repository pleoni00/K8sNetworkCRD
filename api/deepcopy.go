package api

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *NetworkGraph) DeepCopyInto(out *NetworkGraph) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta

	// Safe copy del Spec
	out.Spec = NetworkGraphSpec{
		Nodes: make([]Node, len(in.Spec.Nodes)),
		Edges: make([]Edge, len(in.Spec.Edges)),
	}

	// ✅ Safe copy di Metadata - controlla se è nil
	if in.Spec.Metadata != nil {
		out.Spec.Metadata = in.Spec.Metadata.DeepCopy()
	} else {
		out.Spec.Metadata = nil
	}

	// Copy nodes
	for i := range in.Spec.Nodes {
		in.Spec.Nodes[i].DeepCopyInto(&out.Spec.Nodes[i])
	}

	// Copy edges
	for i := range in.Spec.Edges {
		in.Spec.Edges[i].DeepCopyInto(&out.Spec.Edges[i])
	}

	// ✅ Safe copy del Status
	out.Status = NetworkGraphStatus{
		Phase:     in.Status.Phase,
		NodeCount: in.Status.NodeCount,
		EdgeCount: in.Status.EdgeCount,
	}

	// ✅ Safe copy di LastUpdated - controlla se è nil
	if in.Status.LastUpdated != nil {
		out.Status.LastUpdated = in.Status.LastUpdated.DeepCopy()
	}

	// ✅ Safe copy di Health - controlla se è nil
	if in.Status.Health != nil {
		out.Status.Health = &HealthInfo{
			HealthyNodes:       in.Status.Health.HealthyNodes,
			UnhealthyNodes:     in.Status.Health.UnhealthyNodes,
			TotalTrafficVolume: in.Status.Health.TotalTrafficVolume,
		}
	}
}

func (in *Node) DeepCopyInto(out *Node) {
	out.ID = in.ID
	out.Name = in.Name
	out.Type = in.Type
	out.IP = in.IP
	out.State = in.State
	out.CPUUsageMillicores = in.CPUUsageMillicores
	out.MemoryBytes = in.MemoryBytes
}

func (in *Edge) DeepCopyInto(out *Edge) {
	out.ID = in.ID
	out.Source = in.Source
	out.Target = in.Target
	out.Traffic = TrafficInfo{
		Protocol:    in.Traffic.Protocol,
		Port:        in.Traffic.Port,
		Bytes:       in.Traffic.Bytes,
		Packets:     in.Traffic.Packets,
		Connections: in.Traffic.Connections,
		Latency:     in.Traffic.Latency,
		Bandwidth:   in.Traffic.Bandwidth,
		Direction:   in.Traffic.Direction,
	}

	// ✅ Safe copy di Timestamp
	if in.Timestamp != nil {
		out.Timestamp = in.Timestamp.DeepCopy()
	}

	// ✅ Safe copy di Quality
	if in.Quality != nil {
		out.Quality = &QualityInfo{
			SuccessRate: in.Quality.SuccessRate,
			ErrorRate:   in.Quality.ErrorRate,
			Jitter:      in.Quality.Jitter,
		}
	}
}

// ✅ Safe DeepCopy per NetworkGraphMetadata
func (in *NetworkGraphMetadata) DeepCopy() *NetworkGraphMetadata {
	if in == nil {
		return nil
	}

	out := &NetworkGraphMetadata{
		CollectionInterval: in.CollectionInterval,
		DataSource:         in.DataSource,
		Version:            in.Version,
	}

	// ✅ Safe copy di CollectionTime
	if in.CollectionTime != nil {
		out.CollectionTime = in.CollectionTime.DeepCopy()
	}

	return out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *NetworkGraph) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := NetworkGraph{}
	in.DeepCopyInto(&out)
	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *NetworkGraphList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}

	out := NetworkGraphList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]NetworkGraph, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
