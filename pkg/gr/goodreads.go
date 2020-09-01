// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-07-18

// Package gr is a partial implementation of the Goodreads API.
package gr

import (
	"net/http"
	"time"

	"github.com/mrjones/oauth"
	"github.com/pkg/errors"
)

// Workflow version. set via LD_FLAGS.
var version = ""

// TokenStore loads and saves the access tokens. Load() should return empty
// strings (not an error) if no credentials are stored.
type TokenStore interface {
	Save(token, secret string) error
	Load() (token, secret string, err error)
}

// Logger is the logging interface used by Client.
type Logger interface {
	Printf(string, ...interface{})
	Print(...interface{})
}

type nullLogger struct{}

func (l nullLogger) Printf(format string, args ...interface{}) {}
func (l nullLogger) Print(args ...interface{})                 {}

var _ Logger = nullLogger{}

// Client implements a subset of the Goodreads API.
type Client struct {
	APIKey    string     // Goodreads API key
	APISecret string     // Goodreads API secret
	Store     TokenStore // Persistent store for access tokens
	Log       Logger     // Library logger

	token       *oauth.AccessToken
	apiClient   *http.Client
	lastRequest time.Time
}

// New creates a new Client. It calls store.Load() and passes through any error it returns.
func New(apiKey, apiSecret string, store TokenStore) (*Client, error) {
	c := &Client{
		APIKey:    apiKey,
		APISecret: apiSecret,
		Store:     store,
		Log:       nullLogger{},
	}
	token, secret, err := store.Load()
	if err != nil {
		return nil, errors.Wrap(err, "load OAuth credentials from store")
	}

	if token != "" && secret != "" {
		c.token = &oauth.AccessToken{Token: token, Secret: secret}
	}

	return c, nil
}
