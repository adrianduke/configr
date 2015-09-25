package configr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Encoding int

const (
	Unknown Encoding = iota - 1
	JSON
	TOML
)

type FileSource struct {
	filePath string
	encoding Encoding
}

var SupportedFileExtensions = []string{"json", "toml"}
var ErrUnknownEncoding = errors.New("configr: Unable to determine file encoding, please set manually")

var globalFileSource *FileSource = NewFileSource("")

func NewFileSource(filePath string) *FileSource {
	f := &FileSource{encoding: Unknown}
	f.SetFilePath(filePath)

	return f
}

// SetFilePath sets the file path of the configuration file and try to determine
// the encoding of the file using its extension. See SupportedFileExtensions for
// a list of supported extensions
func SetFilePath(path string) {
	globalFileSource.SetFilePath(path)
}
func (f *FileSource) SetFilePath(path string) {
	f.filePath = path

	fileExt := getFileExtension(path)
	switch fileExt {
	case SupportedFileExtensions[TOML]:
		f.SetEncoding(TOML)
	case SupportedFileExtensions[JSON]:
		f.SetEncoding(JSON)
	}
}

func FilePath() string {
	return globalFileSource.FilePath()
}
func (f *FileSource) FilePath() string {
	return f.filePath
}

// SetEncoding allows the caller to override the infered file encoding format
func SetEncoding(encoding Encoding) {
	globalFileSource.SetEncoding(encoding)
}
func (f *FileSource) SetEncoding(encoding Encoding) {
	f.encoding = encoding
}

func Unmarshal() (map[string]interface{}, error) {
	return globalFileSource.Unmarshal()
}
func (f *FileSource) Unmarshal() (map[string]interface{}, error) {
	var unmarshaller func([]byte, interface{}) error
	values := make(map[string]interface{})

	switch f.encoding {
	case JSON:
		unmarshaller = json.Unmarshal
	case TOML:
		unmarshaller = toml.Unmarshal
	default:
		return values, ErrUnknownEncoding
	}

	fileBytes, err := ioutil.ReadFile(f.filePath)
	if err != nil {
		return values, err
	}

	err = unmarshaller(fileBytes, &values)
	if err != nil {
		return values, err
	}

	return values, nil
}

func Marshal(v interface{}) ([]byte, error) {
	return globalFileSource.Marshal(v)
}
func (f *FileSource) Marshal(v interface{}) ([]byte, error) {
	switch f.encoding {
	case JSON:
		return json.MarshalIndent(v, "", "	")
	case TOML:
		var tomlBytes bytes.Buffer
		tomlEncoder := toml.NewEncoder(bufio.NewWriter(&tomlBytes))
		err := tomlEncoder.Encode(v)
		if err != nil {
			return tomlBytes.Bytes(), err
		}

		return tomlBytes.Bytes(), nil
	default:
		return []byte{}, ErrUnknownEncoding
	}
}

func getFileExtension(filePath string) string {
	return strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))
}
