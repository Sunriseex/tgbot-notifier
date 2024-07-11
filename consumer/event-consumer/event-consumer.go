package event_consumer

import (
	"fmt"
	"time"

	"github.com/sunriseex/tgbot-notifier/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}
func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			fmt.Printf("[ERR] Consumer: %s", err.Error())
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := c.handleEvents(gotEvents); err != nil {
			fmt.Print(err)
			continue
		}
	}
}

// TODO: Параллельная обработка через waitgroup
func (c *Consumer) handleEvents(events []events.Event) error {
	for _, e := range events {
		fmt.Printf("got new event: %s", e.Text)
		if err := c.processor.Process(e); err != nil {
			fmt.Printf("cant handle event: %s", err.Error())
			continue
		}
	}
	return nil
}
