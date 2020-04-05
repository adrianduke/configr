package configr_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/adrianduke/configr"
	"github.com/adrianduke/configr/sources"
	"github.com/adrianduke/configr/sources/file/json"
	"github.com/adrianduke/configr/sources/file/toml"
	"github.com/stretchr/testify/assert"
)

func writeTempFile(t *testing.T, filePath string, content string) {
	err := ioutil.WriteFile(filePath, []byte(content), os.ModePerm)
	assert.NoError(t, err)
}

func Test_ItParsesAllValuesFromJSONConfig(t *testing.T) {
	// Not required outside of this package
	json.Register()

	filePath := "/tmp/test.json"
	writeTempFile(t, filePath, `{
	"t1": "1",
	"t2": {
		"t21": 2,
		"t22": {
			"t221": true
		}
	}
}`)
	defer os.Remove(filePath)
	f := configr.NewFile(filePath)

	config := configr.New()
	config.AddSource(f)
	config.RequireKey("t1", "")
	config.RequireKey("t2.t21", "")
	config.RequireKey("t2.t22.t221", "")
	config.RegisterKey("t3", "", 3)

	assert.NoError(t, config.Parse())

	t1, err := config.String("t1")
	assert.NoError(t, err)
	t2t21, err := config.Float64("t2.t21")
	assert.NoError(t, err)
	t2t22t221, err := config.Bool("t2.t22.t221")
	assert.NoError(t, err)
	t3, err := config.Int("t3")
	assert.NoError(t, err)

	assert.Equal(t, "1", t1)
	assert.Equal(t, float64(2), t2t21)
	assert.Equal(t, true, t2t22t221)
	assert.Equal(t, 3, t3)
}

func Test_ItParsesAllValuesFromTOMLConfig(t *testing.T) {
	// Not required outside of this package
	toml.Register()

	filePath := "/tmp/test.toml"
	writeTempFile(t, filePath, `
[t1]
t11 = "1"
t12 = 2

[t2]
	[t2.t21]
	t211 = false
		[t2.t21.t212]
		t2121 = "4"
`)
	defer os.Remove(filePath)
	f := configr.NewFile(filePath)

	config := configr.New()
	config.AddSource(f)
	config.RequireKey("t1.t11", "")
	config.RequireKey("t1.t12", "")
	config.RequireKey("t2.t21.t211", "")
	config.RequireKey("t2.t21.t212.t2121", "")
	config.RegisterKey("t3", "", "sup")

	assert.NoError(t, config.Parse())

	t1t11, err := config.String("t1.t11")
	assert.NoError(t, err)
	t1t12, err := config.Int("t1.t12")
	assert.NoError(t, err)
	t2t21t211, err := config.Bool("t2.t21.t211")
	assert.NoError(t, err)
	t2t21t212t2121, err := config.Get("t2.t21.t212.t2121")
	assert.NoError(t, err)
	t3, err := config.String("t3")
	assert.NoError(t, err)

	assert.Equal(t, "1", t1t11)
	assert.Equal(t, 2, t1t12)
	assert.Equal(t, false, t2t21t211)
	assert.Equal(t, "4", t2t21t212t2121.(string))
	assert.Equal(t, "sup", t3)
}

func Test_ItGeneratesBlankJSONConfig(t *testing.T) {
	// Not required outside of this package
	json.Register()

	config := configr.New()
	expectedOutput := `{
	"t1": {
		"t11": "*** You need this ***",
		"t12": "*** Me too ***"
	},
	"t2": {
		"t21": {
			"t211": "*** And me ***",
			"t212": {
				"t2121": "*** Also me! ***"
			}
		}
	},
	"t3": 0
}`
	config.RequireKey("t1.t11", "You need this")
	config.RequireKey("t1.t12", "Me too")
	config.RequireKey("t2.t21.t211", "And me")
	config.RequireKey("t2.t21.t212.t2121", "Also me!")
	config.RegisterKey("t3", "", 0)

	f := configr.NewFile("")
	f.SetEncodingName("json")

	configBytes, err := config.GenerateBlank(f)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, string(configBytes))
}

func Test_ItGeneratesBlankTOMLConfig(t *testing.T) {
	// Not required outside of this package
	toml.Register()

	config := configr.New()
	expectedOutput := `t3 = 0

[t1]
  t11 = "*** You need this ***"
  t12 = "*** Me too ***"

[t2]
  [t2.t21]
    t211 = "*** And me ***"
    [t2.t21.t212]
      t2121 = "*** Also me! ***"
`
	config.RequireKey("t1.t11", "You need this")
	config.RequireKey("t1.t12", "Me too")
	config.RequireKey("t2.t21.t211", "And me")
	config.RequireKey("t2.t21.t212.t2121", "Also me!")
	config.RegisterKey("t3", "", 0)

	f := configr.NewFile("")
	f.SetEncodingName("toml")

	configBytes, err := config.GenerateBlank(f)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, string(configBytes))
}

func Test_ItParsesValuesFromEnvironmentalVariables(t *testing.T) {
	os.Setenv("CONFIGR_T1", "1")
	os.Setenv("CONFIGR_T2_T21", "2")
	os.Setenv("CONFIGR_T2_T22_T221", "true")

	config := configr.New()
	config.AddSource(sources.NewEnvVars("configr"))

	config.RequireKey("t1", "")
	config.RequireKey("t2.t21", "")
	config.RequireKey("t2.t22.t221", "")
	config.RegisterKey("t3", "", 3)

	assert.NoError(t, config.Parse())

	t1, err := config.String("t1")
	assert.NoError(t, err)
	t2t21, err := config.Float64("t2.t21")
	assert.NoError(t, err)
	t2t22t221, err := config.Bool("t2.t22.t221")
	assert.NoError(t, err)
	t3, err := config.Int("t3")
	assert.NoError(t, err)

	assert.Equal(t, "1", t1)
	assert.Equal(t, float64(2), t2t21)
	assert.Equal(t, true, t2t22t221)
	assert.Equal(t, 3, t3)
}

func Test_ItUnmarshalsIntoAStruct(t *testing.T) {
	type T22 struct {
		T221 bool
	}
	type T2 struct {
		T21 int
		T22 T22
	}
	type TestConfig struct {
		T1 string
		T2 T2
		T3 string `configr:"t4"`
		T5 bool
	}

	expectedTestConfig := TestConfig{
		T1: "1",
		T2: T2{
			T21: 2,
			T22: T22{
				T221: true,
			},
		},
		T3: "4",
		T5: true,
	}

	source := MemorySource{
		values: map[string]interface{}{
			"t1": "1",
			"t2": map[string]interface{}{
				"t21": 2,
				"t22": map[string]interface{}{
					"t221": true,
				},
			},
			"t4": "4",
		},
	}

	config := configr.New()
	config.AddSource(source)

	config.RequireKey("t1", "")
	config.RequireKey("t2.t21", "")
	config.RequireKey("t2.t22.t221", "")
	config.RegisterKey("t3", "", 3)
	config.RegisterKey("t4", "", "")
	config.RegisterKey("t5", "", true)

	assert.NoError(t, config.Parse())

	testConfig := TestConfig{}
	assert.NoError(t, config.Unmarshal(&testConfig))

	assert.Equal(t, expectedTestConfig, testConfig)
}

func Test_ItRegistersKeysAndDefaultsFromStruct(t *testing.T) {
	type Email struct {
		FromAddress string `configr:",required"`
		Subject     string
		MaxRetries  int `configr:"maximumRetries"`
		RetryOnFail bool
	}
	defaultEmail := Email{
		RetryOnFail: true,
	}

	expectedEmail := Email{
		FromAddress: "test@testing.com",
		Subject:     "",
		MaxRetries:  3,
		RetryOnFail: true,
	}

	source := MemorySource{
		values: map[string]interface{}{
			"fromAddress":    "test@testing.com",
			"maximumRetries": 3,
		},
	}

	config := configr.New()
	config.AddSource(source)
	config.RegisterFromStruct(&defaultEmail, configr.ToLowerCamelCase)

	assert.NoError(t, config.Parse())

	email := Email{}
	assert.NoError(t, config.Unmarshal(&email))

	assert.Equal(t, expectedEmail, email)
}

type MemorySource struct {
	values map[string]interface{}
}

func (m MemorySource) Unmarshal(_ []string, _ configr.KeySplitter) (map[string]interface{}, error) {
	return m.values, nil
}
