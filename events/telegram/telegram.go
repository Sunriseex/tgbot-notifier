package telegram

import (
	"github.com/sunriseex/tgbot-notifier/clients/tg"
	"github.com/sunriseex/tgbot-notifier/events"
	"github.com/sunriseex/tgbot-notifier/lib/storage"
)

type Processor struct {
	tg      *tg.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

func New(client *tg.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      nil,
		storage: nil,
	}

}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(update) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for i, u := range updates {
		res = append(res, event(u))
	}
	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func event(upd tg.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

func fetchType(upd tg.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text

}

func fetchText(upd tg.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
