package conf

import (
	"errors"

	"golang.org/x/sys/windows/registry"
)

// Conf provides a simple means for reading / writing registry values.
type Conf struct {
	key registry.Key
}

// New  creates a new Conf instance.
func New() (*Conf, error) {
	k, _, err := registry.CreateKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\nutsvc`,
		registry.ALL_ACCESS,
	)
	if err != nil {
		return nil, err
	}
	return &Conf{
		key: k,
	}, nil
}

// Get returns the data for the specified value, returning the provided default
// if it is not present.
func (c *Conf) Get(name, def string) (string, error) {
	v, _, err := c.key.GetStringValue(name)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return def, nil
		}
		return "", err
	}
	return v, nil
}

// Set changes the value of the specified key.
func (c *Conf) Set(name, data string) error {
	return c.key.SetStringValue(name, data)
}

// Close frees all of the resources used by the Conf.
func (c *Conf) Close() {
	c.key.Close()
}
