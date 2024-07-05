package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	contentURLFmtStr = "https://xtp.dylibso.com/api/v1/c/%v"
)

// GetContent gets the content at the provided address.
func (c *Client) GetContent(address string) ([]byte, error) {
	if c.xtpToken == "" {
		return nil, ErrXTPTokenEnvVarNotSet
	}

	url := fmt.Sprintf(contentURLFmtStr, url.QueryEscape(address))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(AuthHeader, fmt.Sprintf("Bearer %v", c.xtpToken))
	req.Header.Add(ContentTypeHeader, ContentType)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("url: %v got status code %v: %s", url, res.StatusCode, body)
		return nil, nil
	}

	return body, nil
}
