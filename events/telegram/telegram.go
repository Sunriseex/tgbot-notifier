package telegram

import (
	"errors"

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

var ErrUnknownEventType = errors.New("unknown event type")
var ErrUnknownMetaType = errors.New("unknown meta type")

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
	//allocate memory for updates
	res := make([]events.Event, 0, len(updates))
	for i, u := range updates {
		res = append(res, event(u))
	}
	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch Event.Type {
	case events.Message:
		p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}
func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
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
