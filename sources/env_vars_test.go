package sources

import (
	"os"
	"testing"

	"github.com/adrianduke/configr"
	"github.com/stretchr/testify/assert"
)

func Test_ItUnmarhsalsKeysFromEnvironmentalVariables(t *testing.T) {
	prefix := ""
	keySplitter := configr.NewKeySplitter(".")
	envVars := NewEnvVars(prefix)

	configrKeys := []string{
		"t2.t21",
		"t1",
		"t4.t41.t411.t4111",
		"t3.t31.t311",
		"t5",
		"t6",
	}

	expectedKeyValues := map[string]interface{}{
		"t1":                "1",
		"t2.t21":            "2",
		"t3.t31.t311":       "3.0",
		"t4.t41.t411.t4111": "true",
		"t5":                "",
	}

	os.Clearenv()
	for key, value := range expectedKeyValues {
		os.Setenv(toEnvVarKey(prefix, key, keySplitter), value.(string))
	}

	actual, err := envVars.Unmarshal(configrKeys, keySplitter)

	assert.Nil(t, err)
	assert.Equal(t, expectedKeyValues, actual)
}

func Test_ItUnmarhsalsKeysFromEnvironmentalVariablesWithPrefix(t *testing.T) {
	prefix := "configr"
	keySplitter := configr.NewKeySplitter(".")
	envVars := NewEnvVars(prefix)

	configrKeys := []string{
		"t2.t21",
		"t1",
		"t4.t41.t411.t4111",
		"t3.t31.t311",
		"t5",
		"t6",
	}

	expectedKeyValues := map[string]interface{}{
		"t1":                "1",
		"t2.t21":            "2",
		"t3.t31.t311":       "3.0",
		"t4.t41.t411.t4111": "true",
		"t5":                "",
	}

	os.Clearenv()
	for key, value := range expectedKeyValues {
		os.Setenv(toEnvVarKey(prefix, key, keySplitter), value.(string))
	}

	actual, err := envVars.Unmarshal(configrKeys, keySplitter)

	assert.Nil(t, err)
	assert.Equal(t, expectedKeyValues, actual)
}
