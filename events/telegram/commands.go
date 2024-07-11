package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"

	e "github.com/sunriseex/tgbot-notifier/lib/errors/err"
	"github.com/sunriseex/tgbot-notifier/lib/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, username, text)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)

	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendStart(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, username string, pageURL string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}
	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}
	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}
func (p *Processor) sendStart(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)

}
func isURL(pageURL string) bool {
	u, err := url.Parse(pageURL)
	return err == nil && u.Host != ""
}
