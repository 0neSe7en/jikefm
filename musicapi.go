package main

import (
	"fmt"
	"io"
	"net/url"
	"strings"
)

func toMp3UrlWithId(musicId string) string {
	return fmt.Sprintf("http://music.163.com/song/media/outer/url?id=%s.mp3", musicId)
}

func NeteaseUrlToMp3(urlstr string) string {
	u, err := url.Parse(urlstr)
	if err != nil {
		return ""
	}
	q := u.Query()
	musicId := q.Get("id")
	if musicId != "" {
		return toMp3UrlWithId(musicId)
	}
	// ["", "song", "id"]
	values := strings.Split(u.Path, "/")
	if len(values) >= 3 && values[1] == "song" {
		return toMp3UrlWithId(values[2])
	}
	return ""
}

func NeteaseDownload(mp3Url string) (io.ReadCloser, error) {
	resp, err := Client().Do(NewRequest("GET", mp3Url, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
