package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Configuration map from externalized file
type Configuration map[string]string

var config Configuration

func init() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}
}

// GetIntConfig - get int property by key
func GetIntConfig(key string) int {
	value := GetConfig(key)
	i, _ := strconv.Atoi(value)
	return i
}

// GetConfig - get property by key
func GetConfig(key string) string {
	return config[key]
}
