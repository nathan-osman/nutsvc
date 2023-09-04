package conf

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"sync"

	"github.com/hectane/go-acl"
)

var ErrKeyNotFound = errors.New("key was not found")

// Conf provides a means for reading / writing persistent values to disk in a
// way that keeps the data secure. All methods are considered thread-safe.
type Conf struct {
	mutex    sync.RWMutex
	filename string
	values   map[string]string
}

func (c *Conf) load() error {
	f, err := os.Open(c.filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&c.values)
}

func (c *Conf) save() error {
	f, err := os.Create(c.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := acl.Chmod(c.filename, 0600); err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(&c.values)
}

// New creates a new Conf instance.
func New() (*Conf, error) {
	e, err := os.Executable()
	if err != nil {
		return nil, err
	}
	c := &Conf{
		filename: path.Join(path.Dir(e), "conf.json"),
		values:   make(map[string]string),
	}
	if err := c.load(); err != nil {
		return nil, err
	}
	return c, nil
}

// Get attempts to return the value for the specified key. If the key is not
// found, ErrKeyNotFound is returned.
func (c *Conf) Get(key string) (string, error) {
	defer c.mutex.RUnlock()
	c.mutex.RLock()
	if v, ok := c.values[key]; ok {
		return v, nil
	} else {
		return "", ErrKeyNotFound
	}
}

// GetAll returns all of the currently stored values.
func (c *Conf) GetAll() map[string]string {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	m := make(map[string]string)
	for k, v := range c.values {
		m[k] = v
	}
	return m
}

// SetMultiple merges the specified map of keys/values.
func (c *Conf) SetMultiple(changes map[string]string) error {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	for k, v := range changes {
		c.values[k] = v
	}
	return c.save()
}

// Set changes the value of the specified key.
func (c *Conf) Set(key, value string) error {
	return c.SetMultiple(map[string]string{key: value})
}
