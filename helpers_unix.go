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
	return buf[sz:len(buf)]
}
