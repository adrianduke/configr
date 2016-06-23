package configr

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ItConvertsKeyPathsIntoEnvVarCompatibleOnes(t *testing.T) {
	envVars := NewEnvVars("")

	keys := []string{
		"t1",
		"t2.t21",
		"t3.t31.t311",
		"t4.t41.t411.t4111",
	}

	expectedKeysToUnmarshal := []string{
		"T1",
		"T2_T21",
		"T3_T31_T311",
		"T4_T41_T411_T4111",
	}

	envVars.KeysToUnmarshal(keys, newKeySplitter("."))

	assert.Equal(t, expectedKeysToUnmarshal, envVars.envVarsToUnmarshal)
}

func Test_ItAppendsPrefixToKeysToUnmarsal(t *testing.T) {
	envVars := NewEnvVars("prefix")

	keys := []string{
		"t1",
		"t2.t21",
		"t3.t31.t311",
		"t4.t41.t411.t4111",
	}

	expectedKeysToUnmarshal := []string{
		"PREFIX_T1",
		"PREFIX_T2_T21",
		"PREFIX_T3_T31_T311",
		"PREFIX_T4_T41_T411_T4111",
	}

	envVars.KeysToUnmarshal(keys, newKeySplitter("."))

	assert.Equal(t, expectedKeysToUnmarshal, envVars.envVarsToUnmarshal)
}

func Test_ItUnmarhsalsKeysFromEnvironmentalVariables(t *testing.T) {
	prefix := "configr"
	keySplitter := newKeySplitter(".")
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

	envVars.KeysToUnmarshal(configrKeys, keySplitter)
	actual, err := envVars.Unmarshal()

	assert.Nil(t, err)
	assert.Equal(t, expectedKeyValues, actual)
}
