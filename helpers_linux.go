//go:build linux

package main

import (
	"os"
	"path/filepath"
)

var (
	baseDir = filepath.Join(os.Getenv("HOME"), ".local", "share")
)

func truncBytes(buf []byte, sz int) []byte {
	buf = buf[sz:len(buf)]

	i := len(buf)
	for ; i > 0 && buf[i-1] == 0x0F; i-- {
	}

	return buf[:i]
}
