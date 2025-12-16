package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

func newSearcher(parent string) *dbeaver {
	return &dbeaver{
		Parent: parent,
		Chunks: []string{
			"DBeaverData",
			"workspace6",
			"General",
			".dbeaver",
		},
		Files: map[string]string{
			"creds": "credentials-config.json",
			"datas": "data-sources.json",
		},
	}
}

func (d *dbeaver) get(file string) (string, error) {
	target := filepath.Join(append([]string{d.Parent}, d.Chunks...)...)
	target = filepath.Join(target, d.Files[file])

	if _, err := os.Stat(target); os.IsNotExist(err) {
		return "", err
	}

	return target, nil
}

func (d *dbeaver) getDatas() (*datas, error) {
	dat, err := d.get("datas")
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(dat)
	if err != nil {
		return nil, err
	}

	var res datas
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (d *dbeaver) decrypt() (map[string]creds, error) {
	crd, err := d.get("creds")
	if err != nil {
		return nil, err
	}

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
	raw, err := os.ReadFile(crd)
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
