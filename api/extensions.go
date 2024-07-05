package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	epURLFmtStr = "https://xtp.dylibso.com/api/v1/apps/%v/extension-points"
)

// AppsExtensionPointsResponse represents the response from the /apps endpoint.
type AppsExtensionPointsResponse struct {
	ExtensionPoints []*ExtensionPoint `json:"objects"`
	Next            *string           `json:"next,omitempty"`
	Prev            *string           `json:"prev,omitempty"`
	Total           int               `json:"total"`
	PerPage         int               `json:"perPage"`
}

func (r *AppsExtensionPointsResponse) String() string {
	extensionPoints := make([]string, 0, len(r.ExtensionPoints))
	for _, extensionPoint := range r.ExtensionPoints {
		extensionPoints = append(extensionPoints, extensionPoint.String())
	}
	next, prev := "nil", "nil"
	if r.Next != nil {
		next = *r.Next
	}
	if r.Prev != nil {
		prev = *r.Prev
	}
	return fmt.Sprintf("{ExtensionPoints:[%v],Next:%v,Prev:%v,Total:%v,PerPage:%v}",
		strings.Join(extensionPoints, ","), next, prev, r.Total, r.PerPage)
}

// ExtensionPoint represents an extension point associated with an app.
type ExtensionPoint struct {
	ID    string `json:"id"` // "pattern": "^usr_[a-z0-9]{26}$"
	Name  string `json:"name,omitempty"`
	AppID string `json:"appId,omitempty"` // "pattern": "^app_[a-z0-9]{26}$"
	// Schema *Schema `json:"schema,omitempty"` // TODO
	SchemaYaml string `json:"schemaYaml,omitempty"`
	CreatedAt  string `json:"createdAt,omitempty"` // e.g. "2024-06-26T19:58:22.674Z"
	UpdatedAt  string `json:"updatedAt,omitempty"` // e.g. "2024-06-26T19:58:22.674Z"
}

func (e *ExtensionPoint) String() string {
	return fmt.Sprintf("{ID:%q,Name:%q,AppID:%q,SchemaYaml:%q,CreatedAt:%q,UpdatedAt:%q}",
		e.ID, e.Name, e.AppID, e.SchemaYaml, e.CreatedAt, e.UpdatedAt)
}

// GetAppsExtensionPoints returns extension points for the provided App ID.
func (c *Client) GetAppsExtensionPoints(appID string) (*AppsExtensionPointsResponse, error) {
	if c.xtpToken == "" {
		return nil, ErrXTPTokenEnvVarNotSet
	}

	url := fmt.Sprintf(epURLFmtStr, url.QueryEscape(appID))
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

	resp := &AppsExtensionPointsResponse{}
	if err := jsoncomp.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
