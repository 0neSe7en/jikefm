package jike

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Session struct {
	Token string `json:"token"`
}

func NewSession() *Session {
	session := &Session{}
	jsonFile, err := os.Open(tokenPath)
	if err != nil {
		return NewLoginSession()
	}
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)
	if err = json.Unmarshal(bytes, session); err != nil {
		return NewLoginSession()
	}
	if session.Token == "" {
		return NewLoginSession()
	}
	return session
}

func NewLoginSession() *Session {
	session := &Session{}
	_ = session.Login()
	return session
}

func (s Session) Save() error {
	str, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		log.Fatal(err, "save token fail")
		return err
	}
	return ioutil.WriteFile(tokenPath, str, os.ModePerm)
}

func (s *Session) Login() error {
	body, err := s.Get(createSession)
	if err != nil {
		log.Fatal(err, "create session failed")
	}
	var data struct {
		Uuid string `json:"uuid"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err, "create session failed")
	}

	log.Printf("uuid: %s", data.Uuid)
	GenerateQRCode(data.Uuid)

	logging := false
	attemptCounter := 1
	for !logging {
		logging = WaitLogin(data.Uuid)
		attemptCounter += 1
		if attemptCounter > 3 {
			return errors.New("login tasks too long")
		}
	}
	token := ""
	attemptCounter = 1
	for token == "" {
		token, _ = ConfirmLogin(data.Uuid)
		attemptCounter += 1
		if attemptCounter > 3 {
			return errors.New("login tasks too long")
		}
	}
	s.Token = token
	return s.Save()
}

func (s *Session) Get(url string) ([]byte, error) {
	return s.Request("GET", url, nil)
}

func (s *Session) Post(url string, body io.Reader) ([]byte, error) {
	return s.Request("POST", url, body)
}

func (s *Session) Request(method string, url string, body io.Reader) ([]byte, error) {
	r := newRequest(method, url, body)
	if s.Token != "" {
		r.Header.Set("x-jike-app-auth-jwt", s.Token)
	}
	resp, err := client().Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func WaitLogin(uuid string) bool {
	req := newRequest("GET", waitLogin, nil)
	q := req.URL.Query()
	q.Add("uuid", uuid)
	req.URL.RawQuery = q.Encode()
	resp, err := client().Do(req)
	if err != nil {
		log.Println(err, "wait login fail")
		return false
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		LoggedIn bool `json:"logged_in"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println(err, "wait login fail")
		return false
	}
	return data.LoggedIn
}

func ConfirmLogin(uuid string) (string, error) {
	req := newRequest("GET", confirmLogin, nil)
	q := req.URL.Query()
	q.Add("uuid", uuid)
	req.URL.RawQuery = q.Encode()
	resp, err := client().Do(req)
	if err != nil {
		return "", errors.New("confirm login fail")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", errors.New("confirm login fail")
	}
	return data.Token, nil
}
