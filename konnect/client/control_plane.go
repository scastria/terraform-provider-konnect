package client

const (
	ControlPlanePath    = "v2/control-planes"
	ControlPlanePathGet = ControlPlanePath + "/%s"
)

type ControlPlane struct {
	Id          string              `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description"`
	Config      *ControlPlaneConfig `json:"config,omitempty"`
}
type ControlPlaneConfig struct {
	ClusterType          string `json:"cluster_type,omitempty"`
	ControlPlaneEndpoint string `json:"control_plane_endpoint,omitempty"`
	TelemetryEndpoint    string `json:"telemetry_endpoint,omitempty"`
}
type ControlPlaneCollection struct {
	ControlPlanes []ControlPlane `json:"data"`
}
