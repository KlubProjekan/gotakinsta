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

	feeds, _ := insta.Feed.Tags("kucinglucuk")

	for _, feed := range feeds.Images {
		if feed.MediaType == 2 {
			log.Println("start to download feed photo")
			video, err := urlToBuffer(feed.Videos[0].URL)
			if err != nil {
				log.Println("failed read photo ", feed.Videos[0].URL)
				continue
			}

			photo, err := urlToBuffer(feed.Images.Versions[0].URL)
			if err != nil {
				log.Println("failed read photo ", feed.Images.Versions[0].URL)
				continue
			}

			request := goinsta.PostMediaRequest{
				Video:      video,
				Photo:      photo,
				Caption:    feed.Caption.Text,
				Quality:    feed.NumberOfQualities,
				FilterType: feed.FilterType,
				Item:       feed,
			}

			log.Println("start to repost video")
			_, err = insta.UploadVideo(request)
			if err != nil {
				log.Println("failed to upload photo ", err)
				continue
			}
			log.Println("success repost photo")
		}

		// break
	}

}

var client = &http.Client{
	Timeout: time.Duration(5) * time.Second,
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
