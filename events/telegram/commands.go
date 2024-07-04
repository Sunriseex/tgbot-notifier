package telegram

import (
	"log"
	"net/url"
	"strings"

	"github.com/sunriseex/tgbot-notifier/lib/e"
	"github.com/sunriseex/tgbot-notifier/lib/storage"
)

const (
	RndCmd   = "/rnd"
	helpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	if isAddCmd(text) {
		// TODO: AddPage()

	}

	switch text {
	case RndCmd:
	case helpCmd:
	case StartCmd:
	default:

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

func isAddCmd(text string) bool {
	return isURL(text)

}
func isURL(pageURL string) bool {
	u, err := url.Parse(pageURL)
	return err == nil && u.Host != ""
}
