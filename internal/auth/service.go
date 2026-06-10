package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidToken = errors.New("invalid token")

type Challenge struct {
	Mode        string   `json:"mode"`
	Username    string   `json:"username"`
	TokenHeader string   `json:"tokenHeader"`
	Compat      []string `json:"compat"`
}

type Session struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      User      `json:"user"`
}

type User struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type Service struct {
	username     string
	passwordHash [32]byte
	ttl          time.Duration
	now          func() time.Time
	mu           sync.Mutex
	sessions     map[string]Session
}

func New(username string, password string, ttl time.Duration) *Service {
	if username == "" {
		username = "openlucky"
	}
	if password == "" {
		password = "openlucky-dev"
	}
	if ttl <= 0 {
		ttl = 12 * time.Hour
	}
	return &Service{
		username:     username,
		passwordHash: sha256.Sum256([]byte(password)),
		ttl:          ttl,
		now:          time.Now,
		sessions:     make(map[string]Session),
	}
}

func (s *Service) Challenge() Challenge {
	return Challenge{
		Mode:        "password",
		Username:    s.username,
		TokenHeader: "OpenLucky-Admin-Token",
		Compat:      []string{"Lucky-Admin-Token", "Authorization: Bearer"},
	}
}

func (s *Service) Login(username string, password string) (Session, error) {
	if username != s.username || subtle.ConstantTimeCompare(hash(password), s.passwordHash[:]) != 1 {
		return Session{}, ErrInvalidCredentials
	}

	token, err := randomToken()
	if err != nil {
		return Session{}, err
	}

	session := Session{
		Token:     token,
		ExpiresAt: s.now().Add(s.ttl),
		User:      User{Name: s.username, Role: "admin"},
	}

	s.mu.Lock()
	s.sessions[token] = session
	s.mu.Unlock()
	return session, nil
}

func (s *Service) Validate(token string) (Session, error) {
	if token == "" {
		return Session{}, ErrInvalidToken
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[token]
	if !ok || !session.ExpiresAt.After(s.now()) {
		delete(s.sessions, token)
		return Session{}, ErrInvalidToken
	}
	return session, nil
}

func (s *Service) Logout(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

func hash(value string) []byte {
	sum := sha256.Sum256([]byte(value))
	return sum[:]
}

func randomToken() (string, error) {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}
