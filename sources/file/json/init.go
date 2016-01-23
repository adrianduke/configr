package json

import (
	"encoding/json"

	"github.com/adrianduke/configr"
)

const Name = "json"

func init() {
	Register()
}

func Register() {
	configr.RegisterFileDecoder(Name, configr.FileDecoderAdapter(json.Unmarshal), "json", "JSON")

	jsonEncoder := func(v interface{}) ([]byte, error) {
		return json.MarshalIndent(v, "", "	")
	}
	configr.RegisterFileEncoder(Name, configr.EncoderAdapter(jsonEncoder), "json", "JSON")
}
