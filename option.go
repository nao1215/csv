package csv

// Option is a function that sets a configuration option for CSV struct.
type Option func(c *CSV) error

// WithTabDelimiter is an Option that sets the delimiter to a tab character.
func WithTabDelimiter() Option {
	return func(c *CSV) error {
		c.reader.Comma = '\t'
		return nil
	}
}

// WithHeaderless is an Option that sets the headerless flag to true.
func WithHeaderless() Option {
	return func(c *CSV) error {
		c.headerless = true
		return nil
	}
}
