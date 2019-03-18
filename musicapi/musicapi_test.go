package musicapi

import "testing"

var urls = []string{
	"https://music.163.com/#/song?id=4875306",
	"https://music.163.com/song/4875306/",
	"https://music.163.com/song?id=4875306",
	"https://music.163.com/song/4875306",
	"https://music.163.com/song/4875306?userid=12345",
	"https://music.163.com/song/4875306/?userid=12345",
}

func TestNeteaseUrlToMp3(t *testing.T) {
	expect := "http://music.163.com/song/media/outer/url?id=4875306.mp3"
	for _, urlstr := range urls {
		res := NeteaseUrlToMp3(urlstr)
		if res != expect {
			t.Fatalf("Test %s Expect %s, Got %s", urlstr, expect, res)
		}
	}
}
