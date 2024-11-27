package validate

import (
	"context"
	"net/http"
	"net/url"
)

func Link(ctx context.Context, URL string) bool {

	parsedURL, err := url.Parse(URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return false
	}
	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return false
	}

	if resp.StatusCode == http.StatusNotFound {
		return false
	}

	return true
}
