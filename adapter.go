package configr

// SourceAdapter allows you to convert a func:
//    func() (map[string]interface{}, error)
// into a type that satisfies the Source interface
type SourceAdapter func() (map[string]interface{}, error)

func (f SourceAdapter) Unmarshal() (map[string]interface{}, error) {
	return f()
}

// EncoderAdapter allows you to convert a func:
//    func(interface{}) ([]byte, error)
// into a type that satisfies the Encoder interface
type EncoderAdapter func(interface{}) ([]byte, error)

func (f EncoderAdapter) Marshal(v interface{}) ([]byte, error) {
	return f(v)
}
