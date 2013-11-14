package materials

import (
	"testing"
)

func TestCreateNewUser(t *testing.T) {
	u, err := NewUserFrom("test_data")

	if err != nil {
		t.Fatalf("NewUserFrom returned an error\n")
	}

	if u.Username == "" {
		t.Fatalf("No username\n")
	}

	if u.Apikey == "" {
		t.Fatalf("No apikey\n")
	}
}

func TestSaveUser(t *testing.T) {
	u, _ := NewUserFrom("test_data")
	u.Apikey = "abc123"
	err := u.Save()
	if err != nil {
		t.Fatalf("Save returned error %s\n", err.Error())
	}

	u2, _ := NewUserFrom("test_data")
	if u2.Apikey != "abc123" {
		t.Fatalf("Expected apikey to be abc123, got %s\n", u2.Apikey)
	}
}
