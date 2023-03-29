package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	e "github.com/MaksimUlitin/error"
	"github.com/MaksimUlitin/storage"
)

const (
	RmdCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, ChatID int, username string) error {

	text = strings.TrimSpace(text)

	log.Printf("got new command `%s` from `%s", text, username)

	if isAddCmd(text) {
		return p.savePage(ChatID, text, username)
	}

	switch text {
	case RmdCmd:
		return p.sendRandom(ChatID, username)

	case HelpCmd:
		return p.sendHelp(ChatID)

	case StartCmd:
		return p.sendHello(ChatID)

	default:
		return p.tg.SendMessage(ChatID, msgUnknownCommand)
	}

}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.Wrapiferr("can't do command: save page ", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(context.Background(), page)

	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err

	}

	return nil

}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.Wrapiferr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(context.Background(), username)

	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {

		return nil

	}

	return p.storage.Remove(context.Background(), page)

}

func (p *Processor) sendHelp(chatID int) error {

	return p.tg.SendMessage(chatID, msgHelp)

}

func (p *Processor) sendHello(chatID int) error {

	return p.tg.SendMessage(chatID, msgHello)

}

func isAddCmd(text string) bool {

	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
