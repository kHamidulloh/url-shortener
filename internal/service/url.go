package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
)

type URLRepository interface {
	Save(original, short string) error
	GetOriginal(short string) (string, error)
}

type URLService struct {
	repo URLRepository
}

func NewURLService(repo URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) Shorten(original string) (string, error) {
	if !isValidURL(original) {
		return "", errors.New("invalid URL")
	}
	short := generateShortURL()
	err := s.repo.Save(original, short)
	return short, err
}

func (s *URLService) Resolve(short string) (string, error) {
	return s.repo.GetOriginal(short)
}

func generateShortURL() string {
	x := make([]byte, 4)
	rand.Read(x)
	return strings.TrimRight(base64.URLEncoding.EncodeToString(x), "=")
}

func isValidURL(s string) bool {
	parsed, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}
