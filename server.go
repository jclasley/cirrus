package main

import (
	"encoding/json"
	"fmt"
	"os"
	"net/http"
	"controller"
)

var PORT string
var MSGDIR string

func main() {
	Configure()
	http.HandleFunc("/api/saveMessage", SaveHandler)
	http.HandleFunc("/api/retrieveMessage", GetMsgHandler)
	http.ListenAndServe(PORT, nil)
}

type Config struct {
	PORT string
	MSGDIR string
}

func Configure() {
	conf := Config{}
	f, _ := os.ReadFile("config.json")
	json.Unmarshal(f, &conf)
	PORT = conf.PORT
	MSGDIR = conf.MSGDIR
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST route only", http.StatusMethodNotAllowed)
		return
	}
	user, err := controller.GetUserInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !controller.Authenticate(user) {
		http.Error(w, "Unauthorized user", http.StatusUnauthorized)
		return
	}
	msg, err := controller.ParseMessage(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	StoreFile(user, msg)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Message received")
}

// stores at messages subdir path
func StoreFile(u, m string) error {
	path := fmt.Sprintf("%s/%s.txt", MSGDIR, u)
  f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Error in storing file: %s", err)
	}
	f.WriteString(fmt.Sprintf("%s\n", m))
	return nil
}

func GetMsgHandler(w http.ResponseWriter, r *http.Request) {
	user, err := controller.GetUserInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	path := fmt.Sprintf("%s/%s.txt", MSGDIR, user)
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		fmt.Fprintf(w, "No message found. Please submit a message in JSON format to \"/api/saveMessage\"")
	} else {
		w.Write(b)
	}
}