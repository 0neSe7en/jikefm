package jike

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

const host = "app.jike.ruguoapp.com"
const protocol = "https"
const version = "1.0"

var topicSelected = fmt.Sprintf("%s://%s/%s/messages/history", protocol, host, version)
var createSession = fmt.Sprintf("%s://%s/sessions.create", protocol, host)
var waitLogin = fmt.Sprintf("%s://%s/sessions.wait_for_login", protocol, host)
var confirmLogin = fmt.Sprintf("%s://%s/sessions.wait_for_confirmation", protocol, host)

var tokenPath string

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	tokenPath = filepath.Join(usr.HomeDir, ".local", "jike", "jikefm.json")
	if _, err := os.Stat(filepath.Dir(tokenPath)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(tokenPath), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
