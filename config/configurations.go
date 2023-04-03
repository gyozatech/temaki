package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var vars = map[string]string{}

func init() {
	if err := LoadEnvVariablesFromFile(); err != nil {
		fmt.Println(".env file not provided: ", err)
	}
	LoadEnvVariablesFromOs()
}

// Get gets the config from the .env file or from the environment variables
func Get(name string) string {
	return vars[name]
}

// Set inserts a new config
func Set(name string, value string) {
	vars[name] = value
}

// Override can be used to mock the whole configs in tests
func Override(newConfig map[string]string) {
	vars = newConfig
}

func LoadEnvVariablesFromFile() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.ReplaceAll(line, " ", ""), "#") {
			continue
		}
		keyVal := strings.SplitN(line, "=", 2)
		if len(keyVal) == 2 {
			key := strings.ReplaceAll(keyVal[0], " ", "")
			if strings.Contains(key, "#") {
				continue
			}
			value := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(keyVal[1], "\"", ""), "'", ""), "`", "")
			vars[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func LoadEnvVariablesFromOs() {
	for _, env := range os.Environ() {
		pair := strings.Split(env, "=")
		vars[pair[0]] = pair[1]
	}
}
