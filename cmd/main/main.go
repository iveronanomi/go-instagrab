package main

import (
	"fmt"
	"io"
	"log"
	neturl "net/url"
	"os"
	"path"
	"strings"

	grab "github.com/iveronanomi/goinstagrab"

	"github.com/ahmdrz/goinsta"
)

func main() {
	//fmt.Printf("%#v", []byte{0,1,2,3}[:1])
	//return
	l := log.New(os.Stdout, "", log.Lshortfile|log.Ltime)
	if err := grab.ReadConfig(); err != nil {
		panic(err)
	}
	if err := grab.ReadDump(); err != nil {
		l.Printf("could not read dump file %v", err)
	}
	l.Printf("username: `%s` password: `%s`", grab.Config.UserName, grab.Config.UserPassword)
	walker := goinsta.New(grab.Config.UserName, grab.Config.UserPassword)
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

	//if inbox := walker.Inbox; inbox != nil {
	//	for _, dialog := range inbox.Conversations {
	//		fmt.Printf("%v", dialog)
	//	}
	//}

	for _, uName := range grab.Config.ScanTargets {
		user, err := walker.Profiles.ByName(uName)
		if err != nil {
			l.Print(err)
			return
		}
		feeds := user.Feed()
		j := 0
		for feeds.Next() {
			j++
			for _, img := range feeds.Items {
				_, _, err := Download(&img, fmt.Sprintf("./data/%s/feed/", user.Username), "", user.Username)
				if err != nil {
					l.Printf("getting feed error: %v", err)
					break
				}
			}
			if grab.Config.DeepScan > 0 && grab.Config.DeepScan <= j {
				break
			}
			l.Printf("%d:%s feeds saved: %d", user.ID, user.Username, j)
		}

		stories := user.Stories()
		for stories.Next() {
			j := 0
			for _, story := range stories.Items {
				j++
				_, _, err := Download(&story, fmt.Sprintf("./data/%s/story/", user.Username), "", user.Username)
				if err != nil {
					l.Printf("getting story error: %v", err)
					break
				}
			}
			if grab.Config.DeepScan > 0 && grab.Config.DeepScan <= j {
				break
			}
			l.Printf("%d:%s stories saved: %d", user.ID, user.Username, j)
		}
	}

}

func Download(item *goinsta.Item, folder, name, username string) (imgs, vds string, err error) {
	log.SetPrefix("Download ")
	var u *neturl.URL
	var nname string
	imgFolder := path.Join(folder, "images")
	vidFolder := path.Join(folder, "videos")
	inst := item.Media.Instagram()

	os.MkdirAll(folder, 0777)
	os.MkdirAll(imgFolder, 0777)
	os.MkdirAll(vidFolder, 0777)

	vds = goinsta.GetBest(item.Videos)
	if vds != "" {
		if name == "" {
			u, err = neturl.Parse(vds)
			if err != nil {
				return
			}

			nname = path.Join(vidFolder, path.Base(u.Path))
		} else {
			nname = path.Join(vidFolder, nname)
		}
		imgName := name
		if grab.Dump.IsScanned(username, "videos", imgName) {
			return "", vds, nil
		}
		nname = getname(nname)

		vds, err = download(inst, vds, nname)
		if err != nil {
			grab.Dump.MarkScanned(username, "videos", imgName)
		}
		return "", vds, err
	}

	imgs = goinsta.GetBest(item.Images.Versions)
	if imgs != "" {
		if name == "" {
			u, err = neturl.Parse(imgs)
			if err != nil {
				return
			}

			nname = path.Join(imgFolder, path.Base(u.Path))
		} else {
			nname = path.Join(imgFolder, nname)
		}
		imgName := name
		if grab.Dump.IsScanned(username, "images", imgName) {
			return imgs, vds, nil
		}
		nname = getname(nname)

		imgs, err = download(inst, imgs, nname)
		if err != nil {
			grab.Dump.MarkScanned(username, "videos", imgName)
		}
		return imgs, "", err
	}

	return imgs, vds, fmt.Errorf("cannot find any image or video")
}

func download(inst *goinsta.Instagram, url, dst string) (string, error) {
	file, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer file.Close()
	log.Print(url)
	resp, err := inst.C.Get(url)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, resp.Body)
	return dst, err
}

func getname(name string) string {
	nname := name
	i := 1
	for {
		ext := path.Ext(name)

		_, err := os.Stat(name)
		if err != nil {
			break
		}
		if ext != "" {
			nname = strings.Replace(nname, ext, "", -1)
		}
		name = fmt.Sprintf("%s.%d%s", nname, i, ext)
		i++
	}
	return name
}
