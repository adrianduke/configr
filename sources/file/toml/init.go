package toml

import (
	"bufio"
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/adrianduke/configr"
)

const Name = "toml"

func init() {
	Register()
}

func Register() {
	configr.RegisterFileDecoder(Name, configr.FileDecoderAdapter(toml.Unmarshal), "toml", "TOML")

	tomlEncoder := func(v interface{}) ([]byte, error) {
		var tomlBytes bytes.Buffer
		tomlEncoder := toml.NewEncoder(bufio.NewWriter(&tomlBytes))
		err := tomlEncoder.Encode(v)
		if err != nil {
			return tomlBytes.Bytes(), err
		}

		return tomlBytes.Bytes(), nil
	}
	configr.RegisterFileEncoder(Name, configr.EncoderAdapter(tomlEncoder), "toml", "TOML")
}
