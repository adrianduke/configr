// Configr provides an abstraction above configuration sources, allowing you to
// use a single interface to get all your configuration values
package configr

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/cast"
)

type Manager interface {
	RegisterKey(string, string, interface{}, ...Validator)
	RequireKey(string, string, ...Validator)

	AddSource(Source)

	GenerateBlank(Encoder) ([]byte, error)
}

// Validator is a validation function which would be coupled with a configuration
// key, anytime the config key is found in a Source it's value is validated.
type Validator func(interface{}) error

type Config interface {
	Parse() error
	Parsed() bool
	MustParse()

	Get(string) (interface{}, error)

	String(string) (string, error)
	Bool(string) (bool, error)
	Int(string) (int, error)
	Float64(string) (float64, error)
}

// Source is a source of configuration keys and values, calling unmarshal should
// return a map[string]interface{} of all key/value pairs (nesting is supported)
type Source interface {
	Unmarshal() (map[string]interface{}, error)
}

// Encoder would be used to encode registered and required values (along with
// their defaults or descriptions) into bytes.
type Encoder interface {
	Marshal(interface{}) ([]byte, error)
}

type Configr struct {
	valueValidators    map[string][]Validator
	registeredValues   map[string]string
	requiredValues     map[string]struct{}
	defaultValues      map[string]interface{}
	values             map[string]interface{}
	sources            []Source
	parsed             bool
	configFilePath     string
	keyDelimeter       string
	descriptionWrapper string
}

func New() *Configr {
	return &Configr{
		valueValidators:    make(map[string][]Validator),
		registeredValues:   make(map[string]string),
		requiredValues:     make(map[string]struct{}),
		defaultValues:      make(map[string]interface{}),
		values:             make(map[string]interface{}),
		keyDelimeter:       ".",
		descriptionWrapper: "***",
	}
}

func GetConfigr() *Configr {
	return globalConfigr
}

func SetConfigr(c *Configr) {
	globalConfigr = c
}

var globalConfigr *Configr = New()
var ErrKeyNotFound = errors.New("configr: Key not found")
var ErrParseHasntBeenCalled = errors.New("configr: Trying to get values before calling Parse()")
var ErrNoRegisteredValues = errors.New("configr: No registered values to generate")

type ErrRequiredKeysMissing []string

func (e ErrRequiredKeysMissing) Error() string {
	sort.Strings(e)
	return fmt.Sprintf("configr: Missing required configuration values: %v", []string(e))
}

// RegisterValue registers a configuration key (name) along with a description
// of what the configuration key does, a default value and optional validators
//
// name supports nested notation in the form of '.' delimitered keys, e.g.
//     "user.age.month"
func RegisterKey(name, description string, defaultVal interface{}, validators ...Validator) {
	globalConfigr.RegisterKey(name, description, defaultVal, validators...)
}
func (c *Configr) RegisterKey(name, description string, defaultVal interface{}, validators ...Validator) {
	name = strings.ToLower(name)
	c.registeredValues[name] = description

	if defaultVal != nil {
		c.defaultValues[name] = defaultVal
	}

	if len(validators) > 0 {
		c.valueValidators[name] = validators
	}
}

// RequireValue wraps the RegisterValue() call but upon parsing sources, if the
// configuration key (name) is not found, Parse() will return a
// ErrRequiredValuesMissing error
func RequireKey(name, description string, validators ...Validator) {
	globalConfigr.RequireKey(name, description, validators...)
}
func (c *Configr) RequireKey(name, description string, validators ...Validator) {
	name = strings.ToLower(name)
	c.requiredValues[name] = struct{}{}
	c.RegisterKey(name, description, nil, validators...)
}

// AddSource registers Sources with the Configr instance to Unmarshal()
// when Parse() is called. Sources are parsed in a FILO order, meaning
// the first source added is considered the highest priority, and any
// keys from lower priority sources that are present in a higher will be
// overwritten
func AddSource(p Source) {
	globalConfigr.AddSource(p)
}
func (c *Configr) AddSource(p Source) {
	c.sources = append(c.sources, p)
}

// Parse calls Unmarshal on all registered sources, and caches the subsequent
// key/value's. Additional calls to Parse can be made to add additional config
// from sources.
//
// Sources are called in a FILO order, meaning the first source added is
// considered the highest priority, any keys set from lower priority sources
// found in higher priority will be overwritten.
func Parse() error {
	return globalConfigr.Parse()
}
func (c *Configr) Parse() error {
	if err := c.populateValues(); err != nil {
		return err
	}

	if err := c.checkRequiredKeys(); err != nil {
		return err
	}

	c.parsed = true
	return nil
}

func (c *Configr) checkRequiredKeys() error {
	missingKeys := []string{}

	for requiredKey := range c.requiredValues {
		if _, err := c.get(requiredKey); err != nil {
			missingKeys = append(missingKeys, requiredKey)
		}
	}

	if len(missingKeys) > 0 {
		return ErrRequiredKeysMissing(missingKeys)
	}

	return nil
}

func (c *Configr) populateValues() error {
	for i := len(c.sources) - 1; i >= 0; i-- {
		source := c.sources[i]

		sourceValues, err := source.Unmarshal()
		if err != nil {
			return err
		}
		for key, value := range sourceValues {
			if err := c.set(key, value); err != nil {
				return err
			}
		}
	}

	return nil
}

// MustParse wraps Parse() and will panic if there are any resulting errors
func MustParse() {
	globalConfigr.MustParse()
}
func (c *Configr) MustParse() {
	if err := c.Parse(); err != nil {
		panic(err)
	}
}

func (c *Configr) set(key string, value interface{}) error {
	key = strings.ToLower(key)
	if err := c.runValidators(key, value); err != nil {
		return err
	}

	c.values = c.mergeMap(key, value, c.values)

	return nil
}

func (c *Configr) mergeMap(key string, value interface{}, targetMap map[string]interface{}) map[string]interface{} {
	if reflect.TypeOf(value).Kind() == reflect.Map {
		targetMap = c.traverseSubMap(key, cast.ToStringMap(value), targetMap)
	} else {
		path := strings.SplitN(key, c.keyDelimeter, 2)
		if len(path) == 2 {
			targetMap = c.traverseKeyPath(path[0], path[1], value, targetMap)
		} else {
			targetMap[key] = value
		}
	}

	return targetMap
}

func (c *Configr) traverseKeyPath(currentKey, keyRemainder string, value interface{}, targetMap map[string]interface{}) map[string]interface{} {
	if _, found := targetMap[currentKey]; !found {
		targetMap[currentKey] = make(map[string]interface{})
	}

	targetMap[currentKey] = c.mergeMap(keyRemainder, value, targetMap[currentKey].(map[string]interface{}))

	return targetMap
}

func (c *Configr) traverseSubMap(key string, value map[string]interface{}, targetMap map[string]interface{}) map[string]interface{} {
	for subKey, subValue := range value {
		if _, found := targetMap[key]; !found {
			targetMap[key] = make(map[string]interface{})
		}
		targetMap[key] = c.mergeMap(subKey, subValue, targetMap[key].(map[string]interface{}))
	}

	return targetMap
}

func (c *Configr) runValidators(key string, value interface{}) error {
	keysAndValues, err := c.findKeysAndValuesToValidate(key, value)
	if err != nil {
		return err
	}

	for validatorKey, valueToValidate := range keysAndValues {
		if validators, found := c.valueValidators[validatorKey]; found {
			for _, validate := range validators {
				if err := validate(valueToValidate); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *Configr) findKeysAndValuesToValidate(key string, value interface{}) (map[string]interface{}, error) {
	keysAndValues := make(map[string]interface{})
	if reflect.TypeOf(value).Kind() == reflect.Map {
		for validatorKey := range c.valueValidators {
			if !strings.HasPrefix(validatorKey, key) {
				continue
			}
			valueToValidate := searchMap(cast.ToStringMap(value), strings.Split(validatorKey, c.keyDelimeter)[1:])
			keysAndValues[validatorKey] = valueToValidate
		}
	} else {
		keysAndValues[key] = value
	}

	return keysAndValues, nil
}

// Get can only be called after a Parse() has been done. Keys support the
// nested notation format:
//    "user.age.month"
//
// If a key is not found but has been registered with a default, the default
// will be returned
func Get(key string) (interface{}, error) {
	return globalConfigr.Get(key)
}
func (c *Configr) Get(key string) (interface{}, error) {
	if !c.Parsed() {
		return nil, ErrParseHasntBeenCalled
	}

	return c.get(key)
}

func (c *Configr) get(key string) (interface{}, error) {
	key = strings.ToLower(key)
	if value, found := c.values[key]; found {
		return value, nil
	}

	path := strings.Split(key, c.keyDelimeter)
	parent, found := c.values[path[0]]
	if found {
		if reflect.TypeOf(parent).Kind() == reflect.Map {
			if val := searchMap(cast.ToStringMap(parent), path[1:]); val != nil {
				return val, nil
			}
		}
	}

	if defaultValue, found := c.defaultValues[key]; found {
		return defaultValue, nil
	}

	return nil, ErrKeyNotFound
}

// From github.com/spf13/viper
func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	if next, ok := source[path[0]]; ok {
		switch next.(type) {
		case map[string]interface{}:
			// Type assertion is safe here since it is only reached
			// if the type of `next` is the same as the type being asserted
			return searchMap(next.(map[string]interface{}), path[1:])
		default:
			return next
		}
	} else {
		return nil
	}
}

// String wraps Get() and will attempt to cast the resulting value to a string
// or error
func String(key string) (string, error) {
	return globalConfigr.String(key)
}
func (c *Configr) String(key string) (string, error) {
	val, err := c.Get(key)
	if err != nil {
		return "", err
	}
	return cast.ToStringE(val)
}

// Bool wraps Get() and will attempt to cast the resulting value to a bool
// or error
func Bool(key string) (bool, error) {
	return globalConfigr.Bool(key)
}
func (c *Configr) Bool(key string) (bool, error) {
	val, err := c.Get(key)
	if err != nil {
		return false, err
	}
	return cast.ToBoolE(val)
}

// Int wraps Get() and will attempt to cast the resulting value to a int
// or error
func Int(key string) (int, error) {
	return globalConfigr.Int(key)
}
func (c *Configr) Int(key string) (int, error) {
	val, err := c.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToIntE(val)
}

// Float64 wraps Get() and will attempt to cast the resulting value to a float64
// or error
func Float64(key string) (float64, error) {
	return globalConfigr.Float64(key)
}
func (c *Configr) Float64(key string) (float64, error) {
	val, err := c.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToFloat64E(val)
}

// Parsed lets the caller know if a Parse() call has been made or not
func Parsed() bool {
	return globalConfigr.Parsed()
}
func (c *Configr) Parsed() bool {
	return c.parsed
}

// GenerateBlank generates a 'blank' configuration using the passed Encoder,
// it will honour nested keys, use default values where possible and when not
// fall back to placing the description as the value.
func GenerateBlank(e Encoder) ([]byte, error) {
	return globalConfigr.GenerateBlank(e)
}
func (c *Configr) GenerateBlank(e Encoder) ([]byte, error) {
	if len(c.registeredValues) == 0 {
		return []byte{}, ErrNoRegisteredValues
	}

	blankMap := make(map[string]interface{})
	for key, description := range c.registeredValues {
		if defaultValue, found := c.defaultValues[key]; found {
			blankMap = c.mergeMap(key, defaultValue, blankMap)
		} else {
			blankMap = c.mergeMap(key, c.wrapDescription(description), blankMap)
		}
	}

	return e.Marshal(blankMap)
}

func (c *Configr) wrapDescription(description string) string {
	return strings.Join([]string{c.descriptionWrapper, description, c.descriptionWrapper}, " ")
}
