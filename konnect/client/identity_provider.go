package client

const (
	IdentityProviderPath = "identity-provider"
	EmailClaim           = "email"
	NameClaim            = "name"
	GroupsClaim          = "groups"
)

type IdentityProvider struct {
	Issuer        string            `json:"issuer,omitempty"`
	LoginPath     string            `json:"login_path,omitempty"`
	ClientId      string            `json:"client_id,omitempty"`
	ClientSecret  string            `json:"client_secret,omitempty"`
	Scopes        []string          `json:"scopes,omitempty"`
	ClaimMappings map[string]string `json:"claim_mappings,omitempty"`
}
