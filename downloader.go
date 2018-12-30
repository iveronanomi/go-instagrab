package goinstagrab

import (
	"fmt"
	api "github.com/ahmdrz/goinsta"
	"io"
	"log"
	neturl "net/url"
	"os"
	"path"
	"strings"
)

func GrabMedia(names []string, walker *api.Instagram, l *log.Logger) {
	for _, uName := range names {
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
			if Config.DeepScan > 0 && Config.DeepScan <= j {
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
			if Config.DeepScan > 0 && Config.DeepScan <= j {
				break
			}
			l.Printf("%d:%s stories saved: %d", user.ID, user.Username, j)
		}
	}
}

func Download(item *api.Item, folder, name, username string) (imgs, vds string, err error) {
	log.SetPrefix("Download ")
	var u *neturl.URL
	var nname string
	imgFolder := path.Join(folder, "images")
	vidFolder := path.Join(folder, "videos")
	inst := item.Media.Instagram()

	os.MkdirAll(folder, 0777)
	os.MkdirAll(imgFolder, 0777)
	os.MkdirAll(vidFolder, 0777)

	vds = api.GetBest(item.Videos)
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
		if Dump.IsScanned(username, "videos", imgName) {
			return "", vds, nil
		}
		nname = getname(nname)

		vds, err = download(inst, vds, nname)
		if err != nil {
			Dump.MarkScanned(username, "videos", imgName)
		}
		return "", vds, err
	}

	imgs = api.GetBest(item.Images.Versions)
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
		if Dump.IsScanned(username, "images", imgName) {
			return imgs, vds, nil
		}
		nname = getname(nname)

		imgs, err = download(inst, imgs, nname)
		if err != nil {
			Dump.MarkScanned(username, "videos", imgName)
		}
		return imgs, "", err
	}

	return imgs, vds, fmt.Errorf("cannot find any image or video")
}

func download(inst *api.Instagram, url, dst string) (string, error) {
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
