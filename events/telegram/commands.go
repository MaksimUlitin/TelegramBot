package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"

	//	e "github.com/MaksimUlitin/cliens/l"
	e "github.com/MaksimUlitin/error"
	"github.com/MaksimUlitin/storage"
)

const (
	RmdCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(Text string, ChatID int, Username string) error {

	Text = strings.TrimSpace(Text)

	log.Printf("got new command `%s` from `%s", Text, Username)

	if isAddCmd(Text) {
		return p.savePage(ChatID, Text, Username)
	}

	switch Text {
	case RmdCmd:
		return p.sendRandom(ChatID, Username)

	case HelpCmd:
		return p.sendHelp(ChatID)

	case StartCmd:
		return p.sendHEllo(ChatID)

	default:
		p.tg.SendMessage(ChatID, msgUnknownCommand)
	}

	return nil

}

func (p *Processor) savePage(ChatID int, pageURL string, Username string) (err error) {
	defer func() { err = e.Wrapiferr("can't do command: save page ", err) }()

	page := &storage.Page{
		URL:      pageURL,
		Username: Username,
	}

	isExists, err := p.storage.IsExists(page)

	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(ChatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(ChatID, msgSaved); err != nil {
		return err

	}

	return nil

}

func (p *Processor) sendRandom(ChatID int, Username string) (err error) {
	defer func() { err = e.Wrapiferr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(Username)

	if err != nil && !errors.Is(err, storage.ErrNoSevedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSevedPages) {
		return p.tg.SendMessage(ChatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(ChatID, page.URL); err != nil {

		return nil

	}

	return p.storage.Remove(page)

}

func (p *Processor) sendHelp(ChatID int) error {

	return p.tg.SendMessage(ChatID, msgHelp)

}

func (p *Processor) sendHEllo(ChatID int) error {

	return p.tg.SendMessage(ChatID, msgHello)

}

func isAddCmd(Text string) bool {

	return isUrl(Text)
}

func isUrl(Text string) bool {
	u, err := url.Parse(Text)

	return err == nil && u.Host != ""
}
