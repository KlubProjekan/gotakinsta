package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	goinsta "github.com/klubprojekan/gotakinsta"
)

func main() {
	insta := goinsta.New(
		"bryanjmb11",
		"jombang1122",
	)
	if err := insta.Login(); err != nil {
		log.Println(err)
		return
	}
	defer insta.Logout()

	log.Println("logged in")

	feeds, _ := insta.Feed.Tags("anjinglucu")

	for _, feed := range feeds.Images {
		if feed.MediaType == 2 {
			log.Println("start to download feed video")
			video, err := urlToBuffer(feed.Videos[0].URL)
			if err != nil {
				log.Println("failed read video ", feed.Videos[0].URL)
				continue
			}

			photo, err := urlToBuffer(feed.Images.Versions[0].URL)
			if err != nil {
				log.Println("failed read video ", feed.Images.Versions[0].URL)
				continue
			}

			request := goinsta.PostMediaRequest{
				Video: video,
				Photo: photo,
				Item:  feed,
			}

			log.Println("start to repost video")
			_, err = insta.UploadVideo(request)
			if err != nil {
				log.Println("failed to upload video ", err)
				continue
			}
			log.Println("success repost video")
		} else if feed.MediaType == 1 {
			log.Println("start to download feed photo")

			photo, err := urlToBuffer(feed.Images.Versions[0].URL)
			if err != nil {
				log.Println("failed read photo ", feed.Images.Versions[0].URL)
				continue
			}

			request := goinsta.PostMediaRequest{
				Photo: photo,
				Item:  feed,
			}

			log.Println("start to repost photo")
			_, err = insta.UploadPhoto(request)
			if err != nil {
				log.Println("failed to upload photo ", err)
				continue
			}
			log.Println("success repost photo")
		}

		time.Sleep(time.Duration(10) * time.Second)
	}

}

var client = &http.Client{
	Timeout: time.Duration(150) * time.Second,
}

func urlToBuffer(url string) (io.Reader, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("failed create new request ", err)
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("failed to do http request ", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("failed to read response body ", err)
		return nil, err
	}

	return bytes.NewBuffer(body), nil
}
