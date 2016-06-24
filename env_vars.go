package configr

import (
	"os"
	"strings"
)

const (
	EnvVarSeparator = "_"
)

var lookupEnv = shimLookupEnv

type EnvVars struct {
	prefix             string
	envVarsToUnmarshal []string
	keysToUnmarshal    []string
}

func NewEnvVars(prefix string) *EnvVars {
	return &EnvVars{
		prefix:             prefix,
		envVarsToUnmarshal: []string{},
	}
}

func (e *EnvVars) Unmarshal() (map[string]interface{}, error) {
	returnMap := map[string]interface{}{}

	for i, envVarKey := range e.envVarsToUnmarshal {
		if envVarValue, exists := lookupEnv(envVarKey); exists {
			returnMap[e.keysToUnmarshal[i]] = envVarValue
		}
	}

	return returnMap, nil
}

func (e *EnvVars) KeysToUnmarshal(keys []string, keySplitter KeySplitter) {
	e.keysToUnmarshal = keys

	for _, key := range keys {
		e.envVarsToUnmarshal = append(
			e.envVarsToUnmarshal,
			toEnvVarKey(e.prefix, key, keySplitter),
		)
	}
}

func toEnvVarKey(prefix, key string, keySplitter KeySplitter) string {
	keyParts := keySplitter(strings.ToUpper(key))

	if prefix != "" {
		keyParts = append([]string{strings.ToUpper(prefix)}, keyParts...)
	}

	return strings.Join(keyParts, EnvVarSeparator)
}

// os.lookupEnv was only introduced in go1.5, this is a shim for < go1.5
func shimLookupEnv(key string) (string, bool) {
	var value string
	var found bool

	if value = os.Getenv(key); value == "" {
		// Check key is present in environmental variables
		for _, envVarKeyVal := range os.Environ() {
			envVarKey := strings.Split(envVarKeyVal, "=")[0]
			if envVarKey == key {
				found = true
			}
		}
	} else {
		found = true
	}

	return value, found
}
