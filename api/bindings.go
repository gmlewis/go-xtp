package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	bindingsURLFmtStr = "https://xtp.dylibso.com/api/v1/extension-points/%v/bindings"
)

// BindingsMap represents a collection of bindings keyed by name.
type BindingsMap map[string]*Binding

// Binding represents an XTP Extension Point Binding.
type Binding struct {
	ID             string `json:"id"`
	ContentAddress string `json:"contentAddress"`
	UpdatedAt      string `json:"updatedAt"`
}

// GetExtensionPointBindings returns extension points for the provided App ID.
func (c *Client) GetExtensionPointBindings(ep *ExtensionPoint) (BindingsMap, error) {
	if c.xtpToken == "" {
		return nil, ErrXTPTokenEnvVarNotSet
	}

	url := fmt.Sprintf(bindingsURLFmtStr, url.QueryEscape(ep.ID))
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

	resp := BindingsMap{}
	if err := jsoncomp.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
