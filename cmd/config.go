package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Port  int    `json:"port"`
	Token string `json:"token"`
	Repo  string `json:"repo"`
	Cert  string `json:"cert"`
	Key   string `json:"key"`
	Log   string `json:"log"`
}

func loadConfig(path string, defaultPort int, defaultLog string) (*config, error) {
	cfg := &config{
		Port: defaultPort,
		Log:  defaultLog,
	}
	if !isFile(path) {
		return cfg, nil
	}
	if file, err := os.Open(path); err != nil {
		return nil, err
	} else {
		defer file.Close()
		if err = json.NewDecoder(file).Decode(cfg); err != nil {
			return nil, err
		} else {
			return cfg, nil
		}
	}
}

func isFile(path string) bool {
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		return false
	} else {
		return true
	}
}
