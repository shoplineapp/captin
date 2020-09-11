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

func (d Destination) GetDocumentStore() string {
	_, value := d.Config.GetByEnv("document_store")
	if len(value) == 0 {
		if len(d.Config.DocumentStore) == 0 {
			return "default"
		}
		return d.Config.DocumentStore
	}
	return value
}
