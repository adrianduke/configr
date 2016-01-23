package configr_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/adrianduke/configr"
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
