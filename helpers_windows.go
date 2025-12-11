//go:build windows

package main

import (
	"os"
)

var (
	baseDir = os.Getenv("APPDATA")
)

func truncBytes(buf []byte, sz int) []byte {
	return buf[sz : len(buf)-2]
}
