package your_api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"test_task/pkg/e"
)

type Client struct {
	host   string
	client http.Client
}

func NewClient(host string) *Client {
	return &Client{
		host:   host,
		client: http.Client{},
	}
}

func (c *Client) GetSongInfo(ctx context.Context, group, song string) (*Response, error) {
	const fn = "your_api.GetSongInfo"

	q := url.Values{}
	q.Add("group", group)
	q.Add("song", song)

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "info",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(fn, err)
	}

	req.URL.RawQuery = q.Encode()

	req.Close = true

	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, e.Wrap(fn, err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, e.Wrap(fn, ErrBadRequest)
		} else {
			return nil, e.Wrap(fn, ErrInternalServerError)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(fn, err)
	}

	var res Response

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, e.Wrap(fn, err)
	}

	return &res, nil
}
