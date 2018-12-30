package crawler

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/ahmdrz/goinsta"
	"github.com/iveronanomi/goinstagrab"
)

// LatestMedia ...
func (s *service) LatestMedia() {
	names := goinstagrab.Config.ScanTargets
	for _, uName := range names {
		user, err := s.api.Profiles.ByName(uName)
		if err != nil {
			s.l.Print(err)
			return
		}
		j := 0
		feeds := user.Feed()
		for feeds.Next() {
			j++
			for _, img := range feeds.Items {
				_, _, err := s.download(&img, fmt.Sprintf("./data/%s/feed/", user.Username), "", user.Username)
				if err != nil {
					s.l.Printf("%d:%s stories saved: %d", user.ID, user.Username, j)
					s.l.Printf("getting feed error: %v", err)
					break
				}
			}
			s.l.Printf("%d:%s stories saved: %d", user.ID, user.Username, j)
			if goinstagrab.Config.DeepScan <= j {
				break
			}
		}

		j = 0
		stories := user.Stories()
		for stories.Next() {
			for _, story := range stories.Items {
				j++
				_, _, err := s.download(&story, fmt.Sprintf("./data/%s/story/", user.Username), "", user.Username)
				if err != nil {
					s.l.Printf("%d:%s stories saved: %d", user.ID, user.Username, j)
					s.l.Printf("getting story error: %v", err)
					break
				}
			}
			s.l.Printf("%d:%s stories saved: %d", user.ID, user.Username, j)
			if goinstagrab.Config.DeepScan <= j {
				break
			}
		}
	}
}

func (s *service) download(item *goinsta.Item, folder, name, username string) (imgs, vds string, err error) {
	log.SetPrefix("Download ")
	var u *url.URL
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
			u, err = url.Parse(vds)
			if err != nil {
				return
			}

			nname = path.Join(vidFolder, path.Base(u.Path))
		} else {
			nname = path.Join(vidFolder, nname)
		}
		imgName := name
		if goinstagrab.Dump.IsScanned(username, "videos", imgName) {
			return "", vds, nil
		}
		nname = getname(nname)

		vds, err = get(inst, vds, nname)
		if err != nil {
			goinstagrab.Dump.MarkMediaScanned(username, "videos", imgName)
		}
		return "", vds, err
	}

	imgs = goinsta.GetBest(item.Images.Versions)
	if imgs != "" {
		if name == "" {
			u, err = url.Parse(imgs)
			if err != nil {
				return
			}

			nname = path.Join(imgFolder, path.Base(u.Path))
		} else {
			nname = path.Join(imgFolder, nname)
		}
		imgName := name
		if goinstagrab.Dump.IsScanned(username, "images", imgName) {
			return imgs, vds, nil
		}
		nname = getname(nname)

		imgs, err = get(inst, imgs, nname)
		if err != nil {
			goinstagrab.Dump.MarkMediaScanned(username, "videos", imgName)
		}
		return imgs, "", err
	}

	return imgs, vds, fmt.Errorf("cannot find any image or video")
}

func get(inst *goinsta.Instagram, url, dst string) (string, error) {
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
