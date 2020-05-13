package models

// Destination - Event dispatch destination
type Destination struct {
	Config Configuration
}

func (d Destination) GetCallbackURL() string {
	_, value := d.Config.GetByEnv("callback_url")
	if len(value) == 0 {
		value = d.Config.CallbackURL
	}
	return value
}
