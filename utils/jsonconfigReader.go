package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type JsonConfigReader struct {
	configFile string
	basePath   string
	jsonData   map[string]interface{}
}

func (jsonConfigReader *JsonConfigReader) GetValue(key string) string {
	if val, ok := jsonConfigReader.jsonData[key]; ok {
		return val.(string)
	}

	return ""
}

func (jsonConfigReader *JsonConfigReader) GetArrayValue(key string) []string {
	arrayValS := []string{}
	if arrayValI, ok := jsonConfigReader.jsonData[key]; ok {
		arrayVal := arrayValI.([]interface{})
		for _, strngI := range arrayVal {
			arrayValS = append(arrayValS, strngI.(string))
		}
	}

	return arrayValS
}

func (jsonConfigReader *JsonConfigReader) GetMapValue(key string) map[string]interface{} {
	if val, ok := jsonConfigReader.jsonData[key]; ok {
		return val.(map[string]interface{})
	}

	return map[string]interface{}{}
}

func NewJsonConfigurationReader(environment string, basePath string) ConfigurationReader {
	configFile := resolveConfigFile(environment)

	filename := os.Getenv("GOPATH") + basePath + configFile
	file, _ := os.Open(filename)

	defer file.Close()

	jsonParser := json.NewDecoder(file)
	jsonData := make(map[string]interface{})

	if err := jsonParser.Decode(&jsonData); err != nil {
		fmt.Printf("Could not parse the configuration file %s : %s \n", filename, err)
	}

	return &JsonConfigReader{
		configFile: configFile,
		basePath:   basePath,
		jsonData:   jsonData,
	}
}

func resolveConfigFile(environment string) string {
	switch environment {
	case "dev":
		return "config.dev.json"
	case "prod":
		return "config.prod.json"
	case "preprod":
		return "config.preprod.json"
	case "qa":
		return "config.qa.json"
	default:
		panic("Environment name incorrect : " + environment)
	}
}

func init() {
	NewConfigurationReader = NewJsonConfigurationReader
}
