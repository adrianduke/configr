package configr

// Satisfies github.com/yourheropaul/inj:DatasourceReader interface, not intended
// for regular configr usage. Use `configr.Get(string)` instead.
func (c *Configr) Read(key string) (interface{}, error) {
	return c.Get(key)
}
