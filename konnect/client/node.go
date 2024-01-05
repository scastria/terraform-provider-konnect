package client

const (
	NodePath    = ControlPlanePathGet + "/nodes"
	NodePathGet = NodePath + "/%s"
)

type Node struct {
	ControlPlaneId  string `json:"-"`
	Id              string `json:"id"`
	Version         string `json:"version"`
	Hostname        string `json:"hostname"`
	LastPing        int64  `json:"last_ping"`
	Type            string `json:"type"`
	ConfigHash      string `json:"config_hash"`
	DataPlaneCertId string `json:"data_plane_cert_id"`
}
type NodeCollection struct {
	Nodes []Node `json:"items"`
}
