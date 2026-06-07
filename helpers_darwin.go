//go:build darwin

package main

import (
   "crypto/aes"
   "os"
   "path/filepath"
)

var (
   baseDir = filepath.Join(os.Getenv("HOME"), "Library")
)

func truncBytes(buf []byte, sz int) []byte {
   buf = buf[sz:]

   if len(buf) == 0 {
      return buf
   }

   padLen := int(buf[len(buf)-1])
   if padLen > 0 && padLen <= aes.BlockSize && padLen <= len(buf) {
      buf = buf[:len(buf)-padLen]
   }

   return buf
}
