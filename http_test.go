package wc

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func testServer() http.Handler {
	s := Storage{
		Static: "static",
		Users:  &MemUsers{},
	}

	return NewMux(s)
}

func newTestUser(t *testing.T, h http.Handler, username string) (*User, error) {
	args := url.Values{"username": []string{username}}

	req, err := http.NewRequest("GET", "http://wc.com/api/user/new?"+args.Encode(), nil)
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	var resp = struct {
		Success bool
		Error   *string
		Result  *User
	}{}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(*resp.Error)
	}

	return resp.Result, nil
}

func TestNewUser(t *testing.T) {
	s := testServer()

	user, err := newTestUser(t, s, "Alice")
	if err != nil {
		t.Fatal(err)
	}

	if user.Username != "Alice" {
		t.Fatalf("Expected Username == Alice (was %s)", user.Username)
	}
}

func TestTwoUsers(t *testing.T) {
	s := testServer()

	user1, err := newTestUser(t, s, "Alice")
	if err != nil {
		t.Fatal(err)
	}

	user2, err := newTestUser(t, s, "Bob")
	if err != nil {
		t.Fatal(err)
	}

	if user1.ID == user2.ID {
		t.Fatalf("Expected user1.ID != user2.ID (was %d, %d)", user1.ID, user2.ID)
	}
}
