package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

var apiUrl = "https://app.jike.ruguoapp.com/1.0/messages/history"

type Audio struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Type   string `json:"type"`
}

type Topic struct {
	Content string `json:"content"`
}

type User struct {
	ScreenName string `json:"screenName"`
}

type LinkInfo struct {
	LinkUrl string `json:"linkUrl"`
	Source  string `json:"source"`
	Title   string `json:"title"`
	Audio   Audio  `json:"audio"`
}

type Message struct {
	Content  string   `json:"content"`
	LinkInfo LinkInfo `json:"linkInfo"`
	Topic    Topic    `json:"topic"`
	User     User     `json:"user"`
}

type Response struct {
	Data        []Message `json:"data"`
	LoadMoreKey string    `json:"loadMoreKey"`
}

type Request struct {
	Limit       int    `json:"limit"`
	TopicId     string `json:"topic"`
	LoadMoreKey string `json:"loadMoreKey"`
}

func Fetch(topicId string) (Response, error) {
	return FetchMore(topicId, "")
}

func FetchMore(topicId string, skip string) (Response, error) {
	var res Response
	req := Request{
		TopicId: topicId,
		Limit:   20,
	}
	req.LoadMoreKey = skip
	jsonStr, err := json.Marshal(req)
	if err != nil {
		return res, err
	}
	request := NewRequest("POST", apiUrl, bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	resp, err := Client().Do(request)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func FetchFMOnly(topicId string) ([]Message, string, error) {
	return FetchMoreFM(topicId, "")
}

func FetchMoreFM(topicId string, skip string) ([]Message, string, error) {
	res, err := FetchMore(topicId, skip)
	if err != nil {
		return nil, "", err
	}
	var messages []Message
	for _, msg := range res.Data {
		// only support netease music right now
		if msg.LinkInfo.Source == "163.com" && msg.LinkInfo.Audio.Type == "AUDIO" {
			messages = append(messages, msg)
		}
	}
	return messages, res.LoadMoreKey, nil
}
