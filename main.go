package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type creds struct {
	Connection struct {
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"#connection"`
}

type connection struct {
	Name   string `json:"name"`
	Config struct {
		Host string `json:"host"`
		Port string `json:"port"`
		Data string `json:"database"`
	} `json:"configuration"`
}

type datas struct {
	Connections map[string]connection `json:"connections"`
}

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

func main() {
	pair, err := getPair()
	if err != nil {
		log.Fatal(err)
	}

	cred, err := decryptCredentials(pair[0])
	if err != nil {
		log.Fatal(err)
	}

	data, err := getDBases(pair[1])
	if err != nil {
		log.Fatal(err)
	}

	for name, cons := range data.Connections {
		fmt.Printf("Name: %s\nHost: %s\nData: %s\nUser: %s\nPass: %s\n\n",
			cons.Name, cons.Config.Host, cons.Config.Data,
			cred[name].Connection.User, cred[name].Connection.Password,
		)
	}
}

func decryptCredentials(file string) (map[string]creds, error) {
	key, err := hex.DecodeString("babb4a9f774ab853c96c2d653dfe544a")
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

	var res map[string]creds
	if err := json.Unmarshal(dec[sz:len(dec)-sz/2], &res); err != nil {
		return nil, err
	}

	return res, nil
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

func iif[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}

	return vfalse
}

func getPair() ([2]string, error) {
	base := iif(runtime.GOOS == "windows", os.Getenv("APPDATA"),
		filepath.Join(os.Getenv("HOME"), ".local", "share"),
	)

	points := getPathPoints()
	follow := append([]string{base}, points["chunks"]...)
	target := filepath.Join(follow...)

	pair := [2]string{
		filepath.Join(target, points["files"][0]),
		filepath.Join(target, points["files"][1]),
	}

	for _, file := range pair {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return pair, err
		}
	}

	return pair, nil
}
