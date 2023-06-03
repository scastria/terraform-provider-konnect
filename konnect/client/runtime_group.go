package client

const (
	RuntimeGroupPath    = "runtime-groups"
	RuntimeGroupPathGet = RuntimeGroupPath + "/%s"
)

type RuntimeGroup struct {
	Id          string              `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description"`
	Config      *RuntimeGroupConfig `json:"config,omitempty"`
}
type RuntimeGroupConfig struct {
	ClusterType          string `json:"cluster_type,omitempty"`
	ControlPlaneEndpoint string `json:"control_plane_endpoint,omitempty"`
	TelemetryEndpoint    string `json:"telemetry_endpoint,omitempty"`
}
type RuntimeGroupCollection struct {
	RuntimeGroups []RuntimeGroup `json:"data"`
}
