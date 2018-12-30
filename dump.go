package goinstagrab

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/iveronanomi/goinstagrab/crawler"
)

// Dump ...
var Dump *dump

type dump struct {
	Users         map[string]string                     `json:"users"`
	Media         map[string]map[string]map[string]bool `json:"scanned"`
	Conversations map[string][]crawler.Message          `json:"conversations"`
}

// IsScanned ...
func (d *dump) IsScanned(user, media, filename string) bool {
	if _, ok := d.Media[user]; !ok {
		return false
	}
	if _, ok := d.Media[user][media]; !ok {
		return false
	}
	_, ok := d.Media[user][media][filename]
	return ok
}

// MarkMediaScanned ...
func (d *dump) MarkMediaScanned(user, media, filename string) {
	log.SetPrefix("MarkMediaScanned ")
	if d.Media == nil {
		d.Media = make(map[string]map[string]map[string]bool)
	}
	if d.Media[user] == nil {
		d.Media[user] = make(map[string]map[string]bool)
	}
	if d.Media[user][media] == nil {
		d.Media[user][media] = make(map[string]bool)
	}
	d.Media[user][media][filename] = true
}

// MarkConversationAsRead ...
func (d *dump) DumpMessage(username string, message crawler.Message) {
	if d.Conversations == nil {
		d.Conversations = make(map[string][]crawler.Message)
	}
	d.Conversations[username] = append(d.Conversations[username], message)
}

// ReadDump from local file
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

// SaveDump to local file
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
