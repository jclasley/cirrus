package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"net/http"
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


type AuthMessage struct {
	Message string `json:"message"`
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST route only", http.StatusMethodNotAllowed)
		return
	}
	user, err := GetUserInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !Authenticate(user) {
		http.Error(w, "Unauthorized user", http.StatusUnauthorized)
		return
	}
	var m AuthMessage
	b, err := io.ReadAll(r.Body)
	err = json.Unmarshal(b, &m)
	if err != nil {
		http.Error(w, "Improperly formatted JSON", http.StatusBadRequest)
		return
	}
  if err = StoreFile(user, m.Message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

// retrieves username and throws errors based on pre-defined specs
func GetUserInfo(r *http.Request) (string, error) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		return "", errors.New("Basic auth not provided")
	}
	if pass == "" {
		return "", errors.New("Invalid password")
	}
	return user, nil
}


// retrieve allowed usernames from reqres route, returns bool of whether or not the given username is in the allowed set

func Authenticate(u string) bool {
	res, err := http.Get("https://reqre.in/api/users")
	if err != nil {
		fmt.Printf("Error in Authenticate: %s", err)
		return false
	}
	var r map[string][]map[string]interface{}
	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	json.Unmarshal(b, &r)
	for _,v := range r["data"] {
		if v["email"] == u {
			return true
		}
	}
	return false
}

func GetMsgHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserInfo(r)
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