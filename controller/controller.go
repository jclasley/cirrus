package controller

import (
	"net/http"
	"encoding/json"
	"io"
	"fmt"
	"errors"
)

type AuthMessage struct {
	Message string `json:"message"`
}

func ParseMessage(r *http.Request) (string, error) {
	var m AuthMessage
	b, err := io.ReadAll(r.Body)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return "", err
	}
	return m.Message, nil
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
	res, err := http.Get("https://reqres.in/api/users")
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