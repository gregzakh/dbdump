package main

import (
	"flag"
	"fmt"
)

func main() {
	altBase := flag.String("change-dir", "", "changes the default scan path")
	flag.StringVar(altBase, "c", "", "alias for --change-dir")
	flag.Parse()

	if *altBase != "" {
		baseDir = *altBase
	}

	searcher := newSearcher(baseDir)
	cred, err := searcher.decrypt()
	if err != nil {
		fmt.Printf("decrypt() error: %v\n", err)
		return
	}

	data, err := searcher.getDatas()
	if err != nil {
		fmt.Printf("getDatas() error: %v\n", err)
		return
	}

	for name, cons := range data.Connections {
		fmt.Printf("Name: %s\nHost: %s\nData: %s\nUser: %s\nPass: %s\n\n",
			cons.Name, cons.Config.Host, cons.Config.Data,
			cred[name].Connection.User, cred[name].Connection.Password,
		)
	}
}
