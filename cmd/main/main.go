package main

import (
	"fmt"
	"log"

	grab "github.com/iveronanomi/goinstagrab"
	"os"

	api "github.com/ahmdrz/goinsta"
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

	l.Printf("username: `%s` password: `%s`", grab.Config.UserName, grab.Config.UserPassword)
	walker := api.New(grab.Config.UserName, grab.Config.UserPassword)
	if err := walker.Login(); err != nil {
		l.Print(err)
		return
	}
	defer func() {
		if err := grab.SaveDump(); err != nil {
			l.Print(err)
		}
		if err := walker.Logout(); err != nil {
			panic(err)
		}
	}()

	if err := walker.Inbox.Sync(); err != nil {
		l.Printf("inbox sync error `%v`", err)
		return
	}

	for _, c := range walker.Inbox.Conversations {
		for c.Next() {
			for _, i := range c.Items {
				if i.Type == "text" {
					fmt.Printf("%d: %s\n", i.UserID, i.Text)
				}
			}
		}
		fmt.Println("=======")
		if err := walker.Inbox.Sync(); err != nil {
			l.Printf("inbox sync error `%v`", err)
			return
		}
	}

	//grab.GrabMedia(grab.Config.ScanTargets, walker, l)
}
