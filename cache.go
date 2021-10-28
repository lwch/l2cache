package l2cache

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

// Cache cache object
type Cache struct {
	sync.Mutex
	limit  int
	dir    string
	buffer []byte
	file   *os.File
	offset int64
	closed bool
}

// New new cache
func New(limit int, dir string) (*Cache, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	return &Cache{
		limit: limit,
		dir:   dir,
	}, nil
}

// Close close cache
func (c *Cache) Close() {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return
	}
	c.buffer = nil
	if c.file != nil {
		c.file.Close()
		os.Remove(c.file.Name())
		c.file = nil
	}
	c.closed = true
}

// Write write data
func (c *Cache) Write(data []byte) (int, error) {
	c.Lock()
	defer c.Unlock()
	if c.file != nil {
		return c.file.Write(data)
	}
	if len(c.buffer) > c.limit {
		f, err := ioutil.TempFile(c.dir, "l2")
		if err != nil {
			return 0, err
		}
		c.file = f
		_, err = io.Copy(f, bytes.NewReader(c.buffer))
		if err != nil {
			return 0, err
		}
		c.buffer = nil
	}
	if c.file == nil {
		c.buffer = append(c.buffer, data...)
		return len(data), nil
	}
	_, err := c.file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	return c.file.Write(data)
}

// Read read data
func (c *Cache) Read(data []byte) (int, error) {
	c.Lock()
	defer c.Unlock()
	if c.file != nil {
		_, err := c.file.Seek(c.offset, io.SeekStart)
		if err != nil {
			return 0, err
		}
		n, err := c.file.Read(data)
		if err != nil {
			return n, err
		}
		c.offset += int64(n)
		return n, nil
	}
	n := copy(data, c.buffer[c.offset:])
	if n == 0 {
		return 0, io.EOF
	}
	c.offset += int64(n)
	return n, nil
}

// Limit get limit size
func (c *Cache) Limit() int {
	return c.limit
}
