package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	e "github.com/MaksimUlitin/error"
	//e "github.com/MaksimUlitin/cliens/l"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(Username string) (*Page, error)
	Remove(P *Page) error
	IsExists(P *Page) (bool, error)
}

var (
	ErrNoSevedPages = errors.New("no seved pages")
)

type Page struct {
	URL      string
	Username string
}

func (p Page) Hash() (string, error) {
	h := sha1.New() //h:= sha1.New() можно и через sha1

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.Username); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil

}
