package main

const (
	keyer = "babb4a9f774ab853c96c2d653dfe544a"
)

type creds struct {
	Connection struct {
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"#connection"`
}

type datas struct {
	Connections map[string]struct {
		Name   string `json:"name"`
		Config struct {
			Host string `json:"host"`
			Port string `json:"port"`
			Data string `json:"database"`
		} `json:"configuration"`
	} `json:"connections"`
}
