package main

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
)


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