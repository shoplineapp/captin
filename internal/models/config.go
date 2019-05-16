package models

// Configuration - Webhook Configuration Model
type Configuration struct {
	ConfigID        string   `json:"id"`
	CallbackURL     string   `json:"callback_url"`
	Validate        string   `json:"validate"`
	Actions         []string `json:"actions"`
	Source          string   `json:"source"`
	Throttle        string   `json:"throttle"`
	IncludeDocument bool     `json:"include_document"`
	Name            string   `json:"name"`
}
