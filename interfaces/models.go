package interfaces

import (
	"time"
)

// IncomingHandler - Interface for creating handler to trigger captin execute
type IncomingHandler interface {
	SetConfigMapper(configMapper *ConfigMapperInterface)
	Setup(c CaptinInterface)
}

// ConfigMapperInterface - Interface for config mapper
type ConfigMapperInterface interface {
	ConfigsForKey(eventKey string) []ConfigurationInterface
}

type IncomingEventInterface interface {
	GetTraceInfo() map[string]interface{}
	GetControl() map[string]interface{}
	GetOutstandingDelaySeconds() time.Duration
	IsValid() bool
}

type DestinationInterface interface {
	GetConfig() ConfigurationInterface
	GetCallbackURL() string
	GetSqsSenderConfig(key string) string
	GetDocumentStore() string
}

type ConfigurationInterface interface {
	GetByEnv(key string) (string, string)
	GetThrottleValue() time.Duration
	GetDelayValue() time.Duration
	GetTimeValueMillis(timeValue string) time.Duration
	GetActions() []string
	GetConfigID() string
	GetCallbackURL() string
	GetValidate() string
	GetSource() string
	GetThrottle() string
	GetDelay() string
	GetThrottleTrailingDisabled() bool
	GetKeepThrottledPayloads() bool
	GetKeepThrottledDocuments() bool
	GetIncludeDocument() bool
	GetName() string
	GetAllowLoopback() bool
	GetSender() string
	GetDocumentStore() string
	GetRetryBackoff() []string
	GetIncludeDocumentAttrs() []string
	GetExcludeDocumentAttrs() []string
	GetIncludePayloadAttrs() []string
	GetExcludePayloadAttrs() []string
	GetExtras() map[string]string
}
