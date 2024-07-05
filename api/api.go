// Package api provides methods to access the XTP API.
package api

import (
	"errors"
	"os"

	jsoniter "github.com/json-iterator/go"
)

const (
	AuthHeader        = "Authorization"
	ContentTypeHeader = "Content-Type"
	ContentType       = "application/json; charset=utf-8"
	XTPTokenEnvVar    = "XTP_TOKEN"
)

var (
	ErrXTPTokenEnvVarNotSet = errors.New("env var XTP_TOKEN not set")

	jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary
)

// Client represents an XTP API client.
type Client struct {
	xtpToken string
}

// New returns a new API client.
func New() *Client {
	return &Client{xtpToken: os.Getenv(XTPTokenEnvVar)}
}
