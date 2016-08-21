package utils

type (
	ConfigurationReader interface {
		GetValue(key string) string
		GetMapValue(key string) map[string]interface{}
		GetArrayValue(key string) []string
	}

	configurationReaderFactory func(environment string, basePath string) ConfigurationReader
)

var NewConfigurationReader configurationReaderFactory
