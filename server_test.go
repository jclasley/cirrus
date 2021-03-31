package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestStoreFile(t *testing.T) {
	u := "test@test.com"
	m := "Test message"
	Configure()
	StoreFile(u, m)

	path := fmt.Sprintf("%s/%s.txt", MSGDIR, u)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, instead found nothing", path)
	} else {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Error(err.Error())
		}
		if string(b) != m + "\n" {
			t.Errorf("Expected %s but got %s", m, b)
		}
	}

	os.Remove(path)
}


func TestSaveHandler(t *testing.T) {
	msg := []byte("{\"message\": \"test\"}")
	tests := map[int]struct{
		method string
		route string
		username string
		password string
		message []byte
	}{
		http.StatusUnauthorized: {"POST", "/api/saveMessage", "george.bluth@reqres.in", "", msg},
		http.StatusBadRequest: {"POST", "/api/saveMessage", "george.bluth@reqres.in", "1", nil},
		http.StatusCreated: {"POST", "/api/saveMessage", "george.bluth@reqres.in", "1", msg},
	}

	
	for code, tt := range tests {
		var req *http.Request
		if tt.message != nil {
			rBody := bytes.NewBuffer(tt.message)
			req = httptest.NewRequest(tt.method, tt.route, rBody)
		} else {
			req = httptest.NewRequest(tt.method, tt.route, nil)
		}
		req.SetBasicAuth(tt.username, tt.password)

		rec := httptest.NewRecorder()
		
		SaveHandler(rec, req)

		if rec.Code != code {
			t.Errorf("Expected code %v but got %v", code, rec.Code)
		}
	}
}

func TestGetHandler(t *testing.T) {
	path := "./messages/george.bluth@reqres.in.txt"
	f, _ := os.Create(path)
	f.WriteString("test")

	tests := []struct{
		method string
		route string
		username string
		password string
		want string
	}{
		{"GET", "/api/retrieveMessage", "george.bluth@reqres.i", "1", "No message found. Please submit a message in JSON format to \"/api/saveMessage\""},
		{"GET", "/api/retrieveMessage", "george.bluth@reqres.in", "1", "test"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.route, nil)
		req.SetBasicAuth(tt.username, tt.password)
		rec := httptest.NewRecorder()

		GetMsgHandler(rec, req)

		b, _ := io.ReadAll(rec.Result().Body)
		defer rec.Result().Body.Close()
		if string(b) != tt.want {
			t.Errorf("Expected %s but got %s", tt.want, string(b))
		}
	}
	os.Remove(path)
}