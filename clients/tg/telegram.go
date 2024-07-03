package tg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
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

func New(host, token string) Client {
	return Client{
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

func (c *Client) SendMessages(chatID int, text string) {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	_, err := c.doReq(SendMessagesMethod, q)
	if err != nil {

		log.Printf("Error sending messages: %v", err)
		// return fmt.Errorf("error sending messages: %v", err)
	}
}

func (c *Client) doReq(method string, query url.Values) ([]byte, error) {

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request %v", err)
	}
	req.URL.RawQuery = query.Encode()
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't send request %v", err)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	return body, nil
}
