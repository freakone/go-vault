package models

type Auth struct {
	ClientToken   string   `json:"client_token"`
	Accessor      string   `json:"accessor"`
	Policies      []string `json:"policies"`
	TokenPolicies []string `json:"token_policies"`
	Metadata      struct {
		RoleName string `json:"role_name"`
	} `json:"metadata"`
	LeaseDuration int      `json:"lease_duration"`
	Renewable     bool     `json:"renewable"`
	EntityID      string   `json:"entity_id"`
	TokenType     string   `json:"token_type"`
	Orphan        bool     `json:"orphan"`
	Errors        []string `json:"errors"`
}
