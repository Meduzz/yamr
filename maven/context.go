package maven

import "io"

type (
	Context struct {
		data map[string]interface{}
		chain *Chain
		index int32
	}
)

func NewContext(chain *Chain) *Context {
	return &Context{
		chain:chain,
		data:make(map[string]interface{}),
		index:0,
	}
}

func (c *Context) Set(param string, value interface{}) {
	c.data[param] = value
}

func (c *Context) Get(param string) interface{} {
	iface, exists := c.data[param]

	if exists {
		return iface
	}
	return nil
}

/* Not as elegant as a .Next(), but model's a bit richer. */
func (c *Context) Write(bytes io.ReadCloser) error {
	c.index++
	chain := *c.chain
	return chain[c.index].Write(c, bytes)
}

func (c *Context) Exists() (bool, error) {
	c.index++
	chain := *c.chain
	return chain[c.index].Exists(c)
}

func (c *Context) Read() ([]byte, error) {
	c.index++
	chain := *c.chain
	return chain[c.index].Read(c)
}