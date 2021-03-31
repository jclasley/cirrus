package main

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
)


func TestAuthenticate(t *testing.T) {
	u1 := "george.bluth@reqres.in"
	u2 := "test@test.com"
	truthy := Authenticate(u1) // should be true
	falsy := Authenticate(u2) // should be false

	if !truthy {
		t.Errorf("Expected %s to be in authenticated users", u1)
	} else if falsy {
		t.Errorf("Expected %s to NOT be in authenticated users", u2)
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
	path := fmt.Sprintf("./messages/%s.txt", u)
	StoreFile(u, m)
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
