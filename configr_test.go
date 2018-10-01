package configr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupErroringSources(sources ...*MockSource) {
	sources[0].On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(map[string]interface{}{}, nil)

	sources[1].On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(map[string]interface{}{}, errors.New("!"))

	sources[2].On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(map[string]interface{}{}, nil)
}

func setupNonErroringSources(sources ...*MockSource) {
	for _, source := range sources {
		source.On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(map[string]interface{}{}, nil)
	}
}

func addUnMarshalers(m Manager, sources ...Source) {
	for _, source := range sources {
		m.AddSource(source)
	}
}

func Test_Parse_ItReturnsErrorOnFirstUnMarshalerError(t *testing.T) {
	config := New()
	s1, s2, s3 := &MockSource{}, &MockSource{}, &MockSource{}
	setupErroringSources(s1, s2, s3)
	addUnMarshalers(config, s1, s2, s3)

	err := config.Parse()

	assert.Equal(t, err, errors.New("!"))
}

func Test_Parse_ItDoesntSetUnmarshaldToTrueOnUnMarshalerError(t *testing.T) {
	config := New()
	s1, s2, s3 := &MockSource{}, &MockSource{}, &MockSource{}
	setupErroringSources(s1, s2, s3)
	addUnMarshalers(config, s1, s2, s3)

	config.Parse()

	assert.Equal(t, false, config.parsed)
}

func Test_Parse_ItSetsUnmarshaldToTrueOnSuccessfulParsing(t *testing.T) {
	config := New()
	s1, s2, s3 := &MockSource{}, &MockSource{}, &MockSource{}
	setupNonErroringSources(s1, s2, s3)
	addUnMarshalers(config, s1, s2, s3)

	config.Parse()

	assert.Equal(t, true, config.parsed)
}

func Test_MustParse_ItPanicsOnFirstUnMarshalerError(t *testing.T) {
	config := New()
	s1, s2, s3 := &MockSource{}, &MockSource{}, &MockSource{}
	setupErroringSources(s1, s2, s3)
	addUnMarshalers(config, s1, s2, s3)

	assert.Panics(t, func() { config.MustParse() })
}

func Test_set_ItSetsValue(t *testing.T) {
	config := New()

	assert.NoError(t, config.set("test", 1))
	assert.Equal(t, 1, config.cache["test"].(int))
}

func Test_RegisterKey_ItReturnsErrorOnFirstFailingValidator(t *testing.T) {
	config := New()
	v1 := func(v interface{}) error {
		return errors.New("!")
	}

	config.RegisterKey("test", "", nil, v1)

	assert.Error(t, config.set("test", 1))
}

func Test_set_ItWrapsValidationErrors(t *testing.T) {
	config := New()
	expectedError := NewValidationError("test", errors.New("!!!"))
	v1 := func(v interface{}) error {
		return errors.New("!!!")
	}

	config.RegisterKey("test", "", nil, v1)

	err := config.set("test", 1)

	assert.Equal(t, expectedError, err)
	assert.EqualError(t, err, "Validation error on key 'test': !!!")
}

func Test_set_ItRunsValidatorsWhenSettingValue(t *testing.T) {
	config := New()
	v1HasRun := false
	v1 := func(v interface{}) error {
		v1HasRun = true
		return nil
	}
	v2HasRun := false
	v2 := func(v interface{}) error {
		v2HasRun = true
		return nil
	}

	config.RegisterKey("test", "", nil, v1, v2)
	config.set("test", 1)

	assert.True(t, v1HasRun)
	assert.True(t, v2HasRun)
}

func Test_set_ItHonoursNestedKeysRunningAllValidators(t *testing.T) {
	config := New()
	v1HasRun := false
	v1 := func(v interface{}) error {
		v1HasRun = v.(bool)
		return nil
	}
	v2HasRun := false
	v2 := func(v interface{}) error {
		v2HasRun = v.(bool)
		return nil
	}

	config.RegisterKey("t1.t11.t111", "", nil, v1)
	config.RegisterKey("t1.t12.t121", "", nil, v2)
	config.set("t1", map[string]interface{}{
		"t11": map[string]interface{}{
			"t111": true,
		},
		"t12": map[string]interface{}{
			"t121": true,
		},
	})

	assert.True(t, v1HasRun)
	assert.True(t, v2HasRun)
}

func Test_Parse_ItPopulatesConfigrValuesFromSource(t *testing.T) {
	config := New()
	s1 := &MockSource{}
	expectedValues := map[string]interface{}{
		"t1": 1,
		"t2": 2,
		"t3": 3,
	}

	s1.On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(expectedValues, nil)

	config.AddSource(s1)
	config.Parse()

	assert.Equal(t, expectedValues, config.cache)
}

func Test_AddSource_ItReturnsErrorFromSet(t *testing.T) {
	config := New()
	s1 := &MockSource{}
	expectedValues := map[string]interface{}{
		"t1": 1,
		"t2": 2,
		"t3": 3,
	}
	config.RegisterKey("t2", "", nil, func(v interface{}) error {
		return errors.New("!")
	})

	s1.On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(expectedValues, nil)

	config.AddSource(s1)

	assert.Error(t, config.Parse())
}

func Test_Parse_ItOverwritesValuesFromHigherPrioritySources(t *testing.T) {
	config := New()
	s1, s2 := &MockSource{}, &MockSource{}
	s1Values := map[string]interface{}{
		"t1": 1,
		"t3": 4,
	}
	s2Values := map[string]interface{}{
		"t2": 2,
		"t3": 3,
	}
	expectedValues := map[string]interface{}{
		"t1": 1,
		"t2": 2,
		"t3": 4,
	}

	s1.On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(s1Values, nil)
	s2.On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(s2Values, nil)

	config.AddSource(s1)
	config.AddSource(s2)
	config.Parse()

	assert.Equal(t, expectedValues, config.cache)
}

func Test_MustParse_ItDoesntPanicOnValueNotRegisteredErrors(t *testing.T) {
	config := New()

	s1 := &MockSource{}
	s1Values := map[string]interface{}{
		"t1": 1,
		"t2": 2,
		"t3": 3,
	}

	s1.On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(s1Values, nil)

	config.AddSource(s1)

	assert.NotPanics(t, func() { config.MustParse() })
}

func Test_Get_ItRetrivesNestedValues(t *testing.T) {
	config := New()
	config.cache = map[string]interface{}{
		"t1": map[string]interface{}{
			"t11": 1,
		},
		"t2": map[string]interface{}{
			"t21": map[string]interface{}{
				"t211": 2,
			},
		},
		"t3": map[interface{}]interface{}{
			31: 3,
		},
	}
	config.parsed = true
	t1t11Expected := 1
	t2t21t211Expected := 2
	t331Expected := 3

	t1t11, err := config.Get("t1.t11")
	assert.NoError(t, err)
	assert.Equal(t, t1t11Expected, t1t11)

	t2t21t211, err := config.Get("t2.t21.t211")
	assert.NoError(t, err)
	assert.Equal(t, t2t21t211Expected, t2t21t211)

	t331, err := config.Get("t3.31")
	assert.NoError(t, err)
	assert.Equal(t, t331Expected, t331)
}

func Test_Parse_ItErrorsIfNotAllRequiredValuesAreFound(t *testing.T) {
	config := New()
	s1 := &MockSource{}
	s1Values := map[string]interface{}{
		"t1": 1,
		"t3": 3,
		"t4": map[string]interface{}{
			"t41": map[string]interface{}{
				"t411": 4,
			},
		},
	}

	s1.On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(s1Values, nil)

	config.AddSource(s1)
	config.RequireKey("t1", "")
	config.RequireKey("t2", "")
	config.RequireKey("t3.t31", "")
	config.RequireKey("t4.t41.t411", "")

	assert.Equal(t, ErrRequiredKeysMissing{"t2", "t3.t31"}.Error(), config.Parse().Error())
}

func Test_Parse_ItRespectsNestedValuesFromMultipleSources(t *testing.T) {
	config := New()
	s1, s2 := &MockSource{}, &MockSource{}
	s1Values := map[string]interface{}{
		"t1": false,
		"t2": map[string]interface{}{
			"t22": 5,
		},
		"t3": map[string]interface{}{
			"t31": map[string]interface{}{
				"t312": 6,
			},
		},
	}
	s2Values := map[string]interface{}{
		"t1": true,
		"t2": map[string]interface{}{
			"t21": 1,
			"t22": 2,
		},
		"t3": map[string]interface{}{
			"t31": map[string]interface{}{
				"t311": 3,
				"t312": 4,
			},
		},
	}
	expectedValues := map[string]interface{}{
		"t1": false,
		"t2": map[string]interface{}{
			"t21": 1,
			"t22": 5,
		},
		"t3": map[string]interface{}{
			"t31": map[string]interface{}{
				"t311": 3,
				"t312": 6,
			},
		},
	}

	s1.On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(s1Values, nil)
	s2.On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(s2Values, nil)

	config.AddSource(s1)
	config.AddSource(s2)
	config.Parse()

	assert.Equal(t, expectedValues, config.cache)
}

func Test_set_ItHandlesPathStyleKeysToSetValues(t *testing.T) {
	config := New()
	t1t11 := "1"
	t1t12t121 := int(2)
	t2t21 := float64(3.0)
	expectedValues := map[string]interface{}{
		"t1": map[string]interface{}{
			"t11": t1t11,
			"t12": map[string]interface{}{
				"t121": t1t12t121,
			},
		},
		"t2": map[string]interface{}{
			"t21": t2t21,
		},
	}

	assert.NoError(t, config.set("t1.t11", "1"))
	assert.NoError(t, config.set("t1.t12.t121", 2))
	assert.NoError(t, config.set("t2.t21", 3.0))

	assert.Equal(t, expectedValues, config.cache)
}

func Test_Get_ItErrorsIfYouTryGetBeforeParsing(t *testing.T) {
	config := New()

	_, err := config.Get("test")

	assert.Equal(t, ErrParseHasntBeenCalled, err)
}

func Test_Get_ItReturnsDefaultValueIfNoValueFoundFromSources(t *testing.T) {
	config := New()
	config.parsed = true

	config.RegisterKey("test", "its a test!", 1)

	value, err := config.Get("test")
	assert.NoError(t, err)

	assert.Equal(t, 1, value.(int))
}

func Test_Get_ItReturnsDefaultValuesInSubtree(t *testing.T) {
	config := New()
	config.parsed = true

	config.RequireKey("t1", "")
	config.RegisterKey("t1.t1", "", true)

	value, err := config.Get("t1")
	assert.NoError(t, err)

	assert.True(t, value.(map[string]interface{})["t1"].(bool))
}

func Test_GenerateBlank_ItReturnsErrorIfNoRegisteredValuesToGenerate(t *testing.T) {
	config := New()
	g := &MockGenerator{}

	_, err := config.GenerateBlank(g)
	assert.Equal(t, ErrNoRegisteredValues, err)
}

func Test_GenerateBlank_ItPassesANestedMapOfConfigNamesAndDefaults(t *testing.T) {
	config := New()
	g := &MockGenerator{}
	t1 := int(1)
	t2t21 := "2"
	t2t22 := float64(3.0)
	t3t31t311 := int(4)
	expectedValues := map[string]interface{}{
		"t1": t1,
		"t2": map[string]interface{}{
			"t21": t2t21,
			"t22": t2t22,
		},
		"t3": map[string]interface{}{
			"t31": map[string]interface{}{
				"t311": t3t31t311,
			},
		},
	}

	config.RegisterKey("t1", "test 1", t1)
	config.RegisterKey("t2.t21", "test 2", t2t21)
	config.RegisterKey("t2.t22", "test 3", t2t22)
	config.RegisterKey("t3.t31.t311", "test 4", t3t31t311)

	g.On("Marshal", expectedValues).Return([]byte{}, nil)

	config.GenerateBlank(g)

	g.AssertExpectations(t)
}

func Test_Parse_ItIsCaseSensitiveByDefault(t *testing.T) {
	config := New()
	s1 := &MockSource{}
	expectedValues := map[string]interface{}{
		"T1": 1,
		"T2": 2,
		"T3": 3,
	}

	s1.On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(expectedValues, nil)

	config.AddSource(s1)
	config.Parse()

	assert.Equal(t, expectedValues, config.cache)
}

func Test_Parse_ItIsCaseInensitive(t *testing.T) {
	config := New()
	s1 := &MockSource{}
	values := map[string]interface{}{
		"T1": 1,
		"T2": 2,
		"T3": 3,
	}
	expectedValues := map[string]interface{}{
		"t1": 1,
		"t2": 2,
		"t3": 3,
	}

	s1.On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).Return(values, nil)

	config.AddSource(s1)
	config.SetIsCaseSensitive(false)
	config.Parse()

	assert.Equal(t, expectedValues, config.cache)
}

func Test_Parse_ItPassesAllKeysToUnmarshalToSource(t *testing.T) {
	config := New()
	s1 := &MockSource{}

	t1 := int(1)
	t2t21 := "2"
	t2t22 := float64(3.0)

	expectedKeys := []string{
		"t1",
		"t2.t21",
		"t2.t22",
		"t3.t31.t311",
	}

	config.RegisterKey("t1", "test 1", t1)
	config.RegisterKey("t2.t21", "test 2", t2t21)
	config.RegisterKey("t2.t22", "test 3", t2t22)
	config.RequireKey("t3.t31.t311", "test 4")

	s1.On("Unmarshal", expectedKeys, mock.AnythingOfType("KeySplitter")).Return(map[string]interface{}{}, nil)

	config.AddSource(s1)
	config.Parse()

	s1.AssertExpectations(t)
}

func Test_Parse_ItPassesAValidKeySplittingFuncToSource(t *testing.T) {
	config := New()
	s1 := &MockSource{}

	var result []string
	t1t11t111 := "2"
	expectedValue := []string{"t1", "t11", "t111"}

	config.RegisterKey("t1.t11.t111", "test 1", t1t11t111)

	s1.
		On("Unmarshal", mock.AnythingOfType("[]string"), mock.AnythingOfType("KeySplitter")).
		Return(map[string]interface{}{}, nil).
		Run(func(args mock.Arguments) {
			expectedKeys := args.Get(0).([]string)
			splitterFn := args.Get(1).(KeySplitter)
			result = splitterFn(expectedKeys[0])
		})

	config.AddSource(s1)
	config.Parse()

	assert.Equal(t, expectedValue, result)
}

func Test_Unmarshal_ItReturnsErrorIfConfigrHasntBeenParsed(t *testing.T) {
	config := New()

	assert.Equal(t, ErrParseHasntBeenCalled, config.Unmarshal(&struct{}{}))
}

func Test_UnmarshalKey_ItReturnsAnyGetErrors(t *testing.T) {
	config := New()
	config.parsed = true

	assert.Equal(t, ErrKeyNotFound, config.UnmarshalKey("t1", &struct{}{}))
}

type MockGenerator struct {
	mock.Mock
}

func (m *MockGenerator) Marshal(i interface{}) ([]byte, error) {
	args := m.Called(i)
	return args.Get(0).([]byte), args.Error(1)
}

type MockSource struct {
	mock.Mock
}

func (p *MockSource) Unmarshal(keys []string, splitter KeySplitter) (map[string]interface{}, error) {
	args := p.Called(keys, splitter)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
