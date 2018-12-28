package goinstagrab

import (
	"encoding/json"
	"log"
	"os"
)

var Config *config

type config struct {
	UserPassword string   `json:"user_password"`
	UserName     string   `json:"user_name"`
	ScanTargets  []string `json:"scan_targets"`
	DeepScan     int      `json:"deep_scan"`
}

func ReadConfig() error {
	log.SetPrefix("ReadConfig ")
	var (
		err error
		f   *os.File
	)
	if f, err = os.Open("./config.json"); err != nil {
		return err
	}
	Config = &config{}
	return json.NewDecoder(f).Decode(Config)
}
