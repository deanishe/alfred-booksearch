// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-07-14

package gr

import (
	"context"
	"net/http"
	"os/exec"
	"time"

	"github.com/deanishe/awgo/util"
	"github.com/mrjones/oauth"
	"github.com/pkg/errors"
)

const (
	authServerURL = "localhost:53233"

	oauthTokenURL     = "https://www.goodreads.com/oauth/request_token"
	oauthAuthoriseURL = "https://www.goodreads.com/oauth/authorize"
	oauthAccessURL    = "https://www.goodreads.com/oauth/access_token"
)

// retrieve OAuth token from disk or Goodreads API
func (c *Client) getAuthToken() (*oauth.AccessToken, error) {
	if c.token != nil {
		return c.token, nil
	}

	var err error
	if c.token, err = c.authoriseWorkflow(); err != nil {
		return nil, errors.Wrap(err, "get OAuth token")
	}

	if err := c.Store.Save(c.token.Token, c.token.Secret); err != nil {
		return nil, errors.Wrap(err, "save OAuth token to store")
	}

	// initialise API client
	if c.apiClient, err = c.oauthConsumer().MakeHttpClient(c.token); err != nil {
		return nil, errors.Wrap(err, "make HTTP client")
	}

	return c.token, err
}

func (c *Client) oauthConsumer() *oauth.Consumer {
	consumer := oauth.NewConsumer(c.APIKey, c.APISecret, oauth.ServiceProvider{
		RequestTokenUrl:   oauthTokenURL,
		AuthorizeTokenUrl: oauthAuthoriseURL,
		AccessTokenUrl:    oauthAccessURL,
	})
	consumer.AdditionalAuthorizationUrlParams["name"] = "goodreads"
	return consumer
}

// Authorise initiates authorisation flow.
func (c *Client) Authorise() error {
	_, err := c.getAuthToken()
	return err
}

// execute OAuth authorisation flow
func (c *Client) authoriseWorkflow() (*oauth.AccessToken, error) {
	type response struct {
		token *oauth.AccessToken
		err   error
	}

	var (
		consumer = c.oauthConsumer()
		ch       = make(chan response)
		mux      = http.NewServeMux()
		srv      = &http.Server{
			Addr:         authServerURL,
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 5,
			Handler:      mux,
		}
		tokens = map[string]*oauth.RequestToken{}
		rtoken *oauth.RequestToken
		reqURL string
		err    error
	)

	rtoken, reqURL, err = consumer.GetRequestTokenAndUrl("http://" + authServerURL + "/")
	if err != nil {
		return nil, errors.Wrap(err, "get OAuth request token")
	}
	tokens[rtoken.Token] = rtoken

	if _, err := util.RunCmd(exec.Command("/usr/bin/open", reqURL)); err != nil {
		return nil, errors.Wrap(err, "open OAuth endpoint")
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var (
			values           = r.URL.Query()
			verificationCode = values.Get("oauth_verifier")
			key              = values.Get("oauth_token")
			token            *oauth.AccessToken
		)

		rtoken := tokens[key]
		if rtoken == nil {
			ch <- response{err: errors.New("no request token")}
			return
		}

		c.Log.Print("[oauth] authorising request token ...")
		token, err = consumer.AuthorizeToken(rtoken, verificationCode)
		if err != nil {
			return
		}

		c.Log.Printf("OAuth AccessToken: %#v", token)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
		ch <- response{token: token}
	})

	// start server
	go func() {
		c.Log.Printf("[oauth] starting local webserver on %s ...", authServerURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ch <- response{err: err}
		}
	}()

	// automatically close server after 3 minutes
	timeout := time.AfterFunc(time.Minute*3, func() {
		c.Log.Print("[oauth] automatically stopping server after timeout")
		if err := srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
			c.Log.Printf("[oauth] shutdown: %v", err)
			ch <- response{err: err}
			return
		}
		ch <- response{err: errors.New("OAuth server timeout exceeded")}
	})

	r := <-ch
	timeout.Stop()

	if err := srv.Shutdown(context.Background()); err != nil {
		if err != http.ErrServerClosed {
			return nil, errors.Wrap(err, "OAuth server")
		}
	}

	return r.token, r.err
}

// AuthedClient returns an HTTP client that has Goodreads OAuth tokens.
func (c *Client) AuthedClient() (*http.Client, error) {
	if c.apiClient != nil {
		return c.apiClient, nil
	}

	token, err := c.getAuthToken()
	if err != nil {
		return nil, err
	}

	c.apiClient, err = c.oauthConsumer().MakeHttpClient(token)
	if err != nil {
		return nil, errors.Wrap(err, "make HTTP client")
	}

	return c.apiClient, nil
}
