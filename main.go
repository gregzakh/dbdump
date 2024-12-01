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
	"regexp"
	"runtime"
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

	rex := regexp.MustCompile("^(aut|clo|con|pro|typ|url)")
	for marker, cons := range data["connections"].(map[string]interface{}) {
		usr := cred[marker].(map[string]interface{})["#connection"]
		cfg := cons.(map[string]interface{})["configuration"].(map[string]interface{})

		for k, v := range cfg {
			if rex.MatchString(k) {
				continue
			}
			fmt.Printf("%-13s: %s\n", k, v)
		}

		for k, v := range usr.(map[string]interface{}) {
			fmt.Printf("%-13s: %s\n", k, v)
		}
		fmt.Println()
	}
}

func decryptCredentials(file string) (map[string]interface{}, error) {
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

	var res map[string]interface{}
	if err := json.Unmarshal(dec[sz:len(dec)-sz+1], &res); err != nil {
		return nil, err
	}

	return res, nil
}

func getDBases(file string) (map[string]interface{}, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, err
	}

	return res, nil
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
