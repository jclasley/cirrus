package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	_, err := GetUserInfo(r)
	if err == nil {
		t.Error("Expected empty auth to throw error")
	}

	u := "test@test.com"
	r.SetBasicAuth(u, "")
	
	_, err = GetUserInfo(r)
	if err == nil {
		t.Error("Expected auth with no password to throw error")
	}

	r.SetBasicAuth(u, "test")
	user, err := GetUserInfo(r) 
	if err != nil {
		t.Error(err.Error())
	} else if user != u {
		t.Errorf("Expected %s but got %s", u, user)
	}
}

func TestAuthenticate(t *testing.T) {
	tests := []struct{
		user string
		want bool
	}{
		{"george.bluth@reqres.in", true},
		{"test@test.com", false},
	}
	
	for _, tt := range tests {
		got := Authenticate(tt.user)
		if got != tt.want {
			t.Errorf("%s: Expected %v for %s, got %v", t.Name(), tt.want, tt.user, got)
		}
	}
}

func TestParseMessage(t *testing.T) {
	tests := map[int]struct{
		method string
		route string
		message *AuthMessage
	}{
	  http.StatusCreated:	{"GET", "/", &AuthMessage{Message: "test"}},
	  http.StatusBadRequest:	{"GET", "/", &AuthMessage{}},
	}
	for _, tt := range tests {
		body, _ := json.Marshal(tt.message)
		rBody := bytes.NewBuffer(body)
		r := httptest.NewRequest(tt.method, tt.route, rBody)
		got, err := ParseMessage(r)
		if err != nil {
			t.Error(err.Error())
		} else if got != tt.message.Message {
			t.Errorf("Expected %s, got %s", tt.message.Message, got)
		}
	}
}