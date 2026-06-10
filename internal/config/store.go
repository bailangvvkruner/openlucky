package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	Listen  string          `json:"listen"`
	Theme   string          `json:"theme"`
	Modules map[string]bool `json:"modules"`
}

type Store struct {
	path string
	mu   sync.Mutex
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func Default() Config {
	return Config{
		Listen: "127.0.0.1:16601",
		Theme:  "light",
		Modules: map[string]bool{
			"status":      true,
			"logs":        true,
			"settings":    true,
			"ddns":        true,
			"web":         true,
			"portforward": true,
			"ssl":         true,
			"cron":        true,
		},
	}
}

func (s *Store) Load() (Config, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return Default(), nil
	}
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	if cfg.Listen == "" {
		cfg.Listen = Default().Listen
	}
	if cfg.Theme == "" {
		cfg.Theme = Default().Theme
	}
	if cfg.Modules == nil {
		cfg.Modules = Default().Modules
	}
	return cfg, nil
}

func (s *Store) Save(cfg Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cfg.Listen == "" {
		cfg.Listen = Default().Listen
	}
	if cfg.Theme == "" {
		cfg.Theme = Default().Theme
	}
	if cfg.Modules == nil {
		cfg.Modules = Default().Modules
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
