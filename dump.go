package goinstagrab

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var Dump *dump

type dump struct {
	Users   map[string]string                     `json:"users"`
	Scanned map[string]map[string]map[string]bool `json:"scanned"`
}

func (d *dump) IsScanned(user, media, filename string) bool {
	if _, ok := d.Scanned[user]; !ok {
		return false
	}
	if _, ok := d.Scanned[user][media]; !ok {
		return false
	}
	_, ok := d.Scanned[user][media][filename]
	return ok
}

func (d *dump) MarkScanned(user, media, filename string) {
	log.SetPrefix("MarkScanned ")
	if d.Scanned == nil {
		d.Scanned = make(map[string]map[string]map[string]bool)
	}
	if d.Scanned[user] == nil {
		d.Scanned[user] = make(map[string]map[string]bool)
	}
	if d.Scanned[user][media] == nil {
		d.Scanned[user][media] = make(map[string]bool)
	}
	d.Scanned[user][media][filename] = true
}

func ReadDump() error {
	log.SetPrefix("ReadDump ")
	var (
		err error
		f   *os.File
	)
	if f, err = os.Open("./dump.json"); err != nil {
		log.Print("couldn't read dump.json")
	}
	defer func() { err = f.Close() }()

	Dump = &dump{}
	return json.NewDecoder(f).Decode(Dump)
}

func SaveDump() error {
	log.SetPrefix("SaveDump ")
	var err error
	b, err := json.Marshal(Dump)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile("./dump.json", b, os.ModePerm); err != nil {
		return err
	}
	return err
}
