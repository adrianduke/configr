package configr

type SourceAdapter func() (map[string]interface{}, error)

func (f SourceAdapter) Unmarshal() (map[string]interface{}, error) {
	return f()
}
