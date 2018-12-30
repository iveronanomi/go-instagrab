package main

import (
	"log"
	"os"

	"github.com/ahmdrz/goinsta"
	grab "github.com/iveronanomi/goinstagrab"
	"github.com/iveronanomi/goinstagrab/crawler"
)

func main() {
	//fmt.Printf("%#v", []byte{0,1,2,3}[:1])
	grab.Init()
	l := log.New(os.Stdout, "", log.Lshortfile|log.Ltime)
	if err := grab.ReadConfig(); err != nil {
		panic(err)
	}
	if err := grab.ReadDump(); err != nil {
		l.Printf("could not read dump file %v", err)
	}

	//l.Printf("username: `%s` password: `%s`", grab.Config.UserName, grab.Config.UserPassword)
	api := goinsta.New(grab.Config.UserName, grab.Config.UserPassword)
	if err := api.Login(); err != nil {
		l.Print(err)
		return
	}
	cr := crawler.New(api, l)

	defer func() {
		if err := grab.SaveDump(); err != nil {
			l.Print(err)
		}
		if err := api.Logout(); err != nil {
			panic(err)
		}
	}()

	cr.GrabConversations()
	//grab.GrabMedia(grab.Config.ScanTargets, walker, l)
}
