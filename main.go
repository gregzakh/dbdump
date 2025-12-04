package main

import (
	"fmt"
)

func main() {
	paths, err := getPaths(baseDir)
	if err != nil {
		fmt.Printf("getPaths() error: %v\n", err)
		return
	}

	cred, err := decryptCredentials(paths[0])
	if err != nil {
		fmt.Printf("decryptCredentials() error: %v\n", err)
		return
	}

	data, err := getDBases(paths[1])
	if err != nil {
		fmt.Printf("getDBases() error: %v\n", err)
		return
	}

	for name, cons := range data.Connections {
		fmt.Printf("Name: %s\nHost: %s\nData: %s\nUser: %s\nPass: %s\n\n",
			cons.Name, cons.Config.Host, cons.Config.Data,
			cred[name].Connection.User, cred[name].Connection.Password,
		)
	}
}
