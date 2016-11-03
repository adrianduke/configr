package sources

import (
	"os"
	"strings"

	"github.com/adrianduke/configr"
)

const (
	EnvVarSeparator = "_"
)

var lookupEnv = shimLookupEnv

type EnvVars struct {
	prefix string
}

func NewEnvVars(prefix string) *EnvVars {
	return &EnvVars{
		prefix: prefix,
	}
}

func (e *EnvVars) Unmarshal(keys []string, keySplitter configr.KeySplitter) (map[string]interface{}, error) {
	returnMap := map[string]interface{}{}

	for _, key := range keys {
		if envVarValue, exists := lookupEnv(toEnvVarKey(e.prefix, key, keySplitter)); exists {
			returnMap[key] = envVarValue
		}
	}

	return returnMap, nil
}

func toEnvVarKey(prefix, key string, keySplitter configr.KeySplitter) string {
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
