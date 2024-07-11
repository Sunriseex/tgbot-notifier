package tg

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/sunriseex/tgbot-notifier/lib/e"
)

type Client struct {
	host     string
	basePath string
	Client   http.Client
}

const (
	getUpdatesMethod   = "getUpdates"
	SendMessagesMethod = "sendMessages"
)

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		Client:   http.Client{},
	}

}

func newBasePath(token string) string {
	return "bot" + token

}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doReq(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}
	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	_, err := c.doReq(SendMessagesMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) doReq(method string, query url.Values) (data []byte, err error) {

	defer func() { err = e.WrapIfErr("can't save", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = query.Encode()
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
