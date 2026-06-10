package auth

import "testing"

func TestChallengeLoginAndValidate(t *testing.T) {
	service := New("admin", "secret", 0)
	challenge := service.Challenge()
	if challenge.TokenHeader != "OpenLucky-Admin-Token" {
		t.Fatalf("TokenHeader = %q", challenge.TokenHeader)
	}

	session, err := service.Login("admin", "secret")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if session.Token == "" {
		t.Fatal("Token is empty")
	}
	validated, err := service.Validate(session.Token)
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if validated.User.Name != "admin" {
		t.Fatalf("User.Name = %q", validated.User.Name)
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	service := New("admin", "secret", 0)
	if _, err := service.Login("admin", "wrong"); err != ErrInvalidCredentials {
		t.Fatalf("Login error = %v, want ErrInvalidCredentials", err)
	}
}
