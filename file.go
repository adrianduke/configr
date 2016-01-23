package configr

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type FileDecoder interface {
	Unmarshal([]byte, interface{}) error
}

var (
	RegisteredFileEncoders = make(map[string]Encoder)
	RegisteredFileDecoders = make(map[string]FileDecoder)
	ExtensionToDecoderName = make(map[string]string)
	ExtensionToEncoderName = make(map[string]string)

	ErrUnknownEncoding = errors.New("configr: Unable to determine file encoding, please set manually")
)

func RegisterFileEncoder(name string, encoder Encoder, fileExtensions ...string) {
	RegisteredFileEncoders[name] = encoder

	registerExtensionsToMap(name, fileExtensions, ExtensionToEncoderName)
}

func RegisterFileDecoder(name string, source FileDecoder, fileExtensions ...string) {
	RegisteredFileDecoders[name] = source

	registerExtensionsToMap(name, fileExtensions, ExtensionToDecoderName)
}

func registerExtensionsToMap(name string, fileExtensions []string, extMap map[string]string) {
	if len(fileExtensions) == 0 {
		fileExtensions = append(fileExtensions, name)
	}

	for _, fileExtension := range fileExtensions {
		extMap[fileExtension] = name
	}
}

type File struct {
	filePath     string
	encodingName string
}

func NewFile(path string) *File {
	f := &File{}
	f.SetPath(path)

	return f
}

func (f *File) SetPath(path string) {
	f.filePath = path

	fileExt := getFileExtension(path)
	if encodingName, found := ExtensionToDecoderName[fileExt]; found {
		f.encodingName = encodingName
	}
}

func (f *File) SetEncodingName(name string) {
	f.encodingName = name
}

func (f *File) Path() string {
	return f.filePath
}

func (f *File) Unmarshal() (map[string]interface{}, error) {
	if decoder, found := RegisteredFileDecoders[f.encodingName]; found {
		values := make(map[string]interface{})

		fileBytes, err := ioutil.ReadFile(f.filePath)
		if err != nil {
			return values, err
		}

		if err := decoder.Unmarshal(fileBytes, &values); err != nil {
			return values, err
		}

		return values, nil
	}

	return map[string]interface{}{}, ErrUnknownEncoding
}

func (f *File) Marshal(v interface{}) ([]byte, error) {
	if encoder, found := RegisteredFileEncoders[f.encodingName]; found {
		return encoder.Marshal(v)
	}

	return []byte{}, ErrUnknownEncoding
}

func getFileExtension(filePath string) string {
	return strings.TrimPrefix(filepath.Ext(filePath), ".")
}
