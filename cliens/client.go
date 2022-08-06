package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	e "github.com/MaksimUlitin/error"
	//	e "github.com/MaksimUlitin/cliens/l"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	GetUpdatesMethod  = "GetUpdates"
	sendmessageMethod = "sendMessage"
)

func New(host string, token /*basaPath*/ string) *Client {
	return &Client{
		host:     host,
		basePath: NewBasePath(token),
		client:   http.Client{},
	}
}

func NewBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) (updates []Updates, err error) {
	defer func() { err = e.Wrapiferr("can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	datas, err := c.doRequest(GetUpdatesMethod, q)

	if err != nil {
		return nil, err
	}

	var res UpdateRespouns

	if err := json.Unmarshal(datas, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {

	q := url.Values{}
	q.Add("chatID", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendmessageMethod, q)

	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil

}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {

	defer func() { err = e.Wrapiferr("can't do request", err) }() //defer func() { err = e.WrapIfErr("can't get updates", err) }()
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, err //e.Wrap(errMsg,err)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err //e.Wrap(errMsg, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
