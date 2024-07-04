package telegram

import (
	"log"
	"net/url"
	"strings"
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

func (p *Processor) savePage(chatID int, username string, text string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

}

func isAddCmd(text string) bool {
	return isURL(text)

}
func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
