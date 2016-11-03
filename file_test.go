package configr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func resetGlobals() func() {
	RegisteredFileEncoders = make(map[string]Encoder)
	RegisteredFileDecoders = make(map[string]FileDecoder)
	ExtensionToDecoderName = make(map[string]string)
	ExtensionToEncoderName = make(map[string]string)

	return func() {
		resetGlobals()
	}
}

func Test_ItRegistersNameAsAFileDecoderExtension(t *testing.T) {
	defer resetGlobals()()

	name := "json"
	source := FileDecoderAdapter(nil)
	expectedExtensions := map[string]string{
		name: name,
	}

	RegisterFileDecoder(name, source)

	assert.Equal(t, expectedExtensions, ExtensionToDecoderName)
}

func Test_ItRegistersAllFileExtensionsForDecoder(t *testing.T) {
	defer resetGlobals()()

	name := "json"
	source := FileDecoderAdapter(nil)
	expectedExtensions := map[string]string{
		"js":   name,
		"JSON": name,
		name:   name,
	}

	RegisterFileDecoder(name, source, "js", "JSON", "json")

	assert.Equal(t, expectedExtensions, ExtensionToDecoderName)
}

func Test_ItInfersEncodingNameByFileExtension(t *testing.T) {
	defer resetGlobals()()

	ExtensionToDecoderName["js"] = "json"

	f := NewFile("/tmp/config.js")

	assert.Equal(t, "json", f.encodingName)
}

func Test_ItErrorsIfItCantFindFileEncoding(t *testing.T) {
	defer resetGlobals()()
	f := NewFile("/tmp/config.js")

	_, err := f.Unmarshal([]string{}, nil)

	assert.EqualError(t, err, ErrUnknownEncoding.Error())
}

func Test_ItRegistersNameAsAnEncoderExtension(t *testing.T) {
	defer resetGlobals()()

	name := "json"
	source := EncoderAdapter(nil)
	expectedExtensions := map[string]string{
		name: name,
	}

	RegisterFileEncoder(name, source)

	assert.Equal(t, expectedExtensions, ExtensionToEncoderName)
}

func Test_ItRegistersAllFileExtensionsForEncoder(t *testing.T) {
	defer resetGlobals()()

	name := "json"
	source := EncoderAdapter(nil)
	expectedExtensions := map[string]string{
		"js":   name,
		"JSON": name,
		name:   name,
	}

	RegisterFileEncoder(name, source, "js", "JSON", "json")

	assert.Equal(t, expectedExtensions, ExtensionToEncoderName)
}

func Test_MarshalingErrorsIfItCantFindFileEncoding(t *testing.T) {
	defer resetGlobals()()
	f := NewFile("/tmp/config.js")

	_, err := f.Marshal(nil)

	assert.EqualError(t, err, ErrUnknownEncoding.Error())
}

type MockFileSourcer struct {
	mock.Mock
}

func (m MockFileSourcer) Unmarshal(b []byte, v interface{}) error {
	args := m.Called(b, v)
	return args.Error(0)
}
