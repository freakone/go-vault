package models

import "time"

type Login struct {
	RequestID     string      `json:"request_id"`
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int         `json:"lease_duration"`
	Data          interface{} `json:"data"`
	WrapInfo      interface{} `json:"wrap_info"`
	Warnings      interface{} `json:"warnings"`
	Auth          struct {
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
	} `json:"auth"`
}

type Data struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		Data     string `json:"data"`
		Metadata struct {
			CreatedTime  time.Time `json:"created_time"`
			DeletionTime string    `json:"deletion_time"`
			Destroyed    bool      `json:"destroyed"`
			Version      int       `json:"version"`
		} `json:"metadata"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}
