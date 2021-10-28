package l2cache

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestCache(t *testing.T) {
	cache, err := New(1024, os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()

	// cache in memory
	_, err = cache.Write([]byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
	if cache.file != nil {
		t.Fatal("unexpected create file")
	}
	data, err := ioutil.ReadAll(cache)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, []byte("hello world")) {
		t.Fatal("invalid data")
	}

	// cache in file
	for i := 0; i < 1024; i++ {
		_, err = cache.Write([]byte{0})
		if err != nil {
			t.Fatal(err)
		}
	}
	if cache.file == nil {
		t.Fatal("unexpected create file")
	}
	data, err = ioutil.ReadAll(cache)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, bytes.Repeat([]byte{0}, 1024)) {
		t.Fatal("invalid data")
	}
}
