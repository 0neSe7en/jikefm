package jike

import (
	"bytes"
	"encoding/json"
)

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

func (m Message) GetTitle() string {
	return m.LinkInfo.Title
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

func FetchMoreSelected(session *Session, topicId string, skip string) (res Response, err error) {
	req := Request{
		TopicId: topicId,
		Limit:   20,
	}
	req.LoadMoreKey = skip
	jsonStr, err := json.Marshal(req)

	if err != nil {
		return
	}

	body, err := session.Post(topicSelected, bytes.NewBuffer(jsonStr))

	if err != nil {
		return
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return res, err
	}
	return res, err
}

func FetchMoreSelectedFM(session *Session, topicId string, skip string) ([]Message, string, error) {
	res, err := FetchMoreSelected(session, topicId, skip)
	if err != nil {
		return nil, "", err
	}
	var messages []Message
	for _, msg := range res.Data {
		if msg.LinkInfo.Source == "163.com" && msg.LinkInfo.Audio.Type == "AUDIO" {
			messages = append(messages, msg)
		}
	}
	return messages, res.LoadMoreKey, nil
}
