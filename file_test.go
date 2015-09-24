package configr

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ItErrorsIfFileExtensionIsNotSupported(t *testing.T) {
	f := NewFileSource("/tmp/test.zoml")

	_, err := f.Unmarshal()
	assert.Equal(t, ErrUnknownEncoding, err)
}

func Test_ItReturnsAnyParserErrors(t *testing.T) {
	filePath := "/tmp/test.json"
	err := ioutil.WriteFile(filePath, []byte(`}{`), os.ModePerm)
	defer os.Remove(filePath)
	assert.NoError(t, err)

	f := NewFileSource(filePath)

	_, err = f.Unmarshal()
	assert.Error(t, err)
}
