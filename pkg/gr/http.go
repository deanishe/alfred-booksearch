// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-07-18

package gr

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	userAgent  string
	httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   30 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
		},
	}
)

func init() {
	userAgent = "Alfred Booksearch Workflow " + version + " (+https://github.com/deanishe/alfred-booksearch)"
}

// retrieve URL with standard HTTP client.
func (c *Client) httpGet(URL string) ([]byte, error) {
	return c.httpRequest(URL, httpClient)
}

// retrieve URL with authorised API client.
func (c *Client) apiRequest(URL string, method ...string) ([]byte, error) {
	var (
		client *http.Client
		data   []byte
		err    error
	)
	if client, err = c.AuthedClient(); err != nil {
		return nil, errors.Wrap(err, "get API client")
	}

	d := time.Since(c.lastRequest)
	if d < time.Second {
		d = time.Second - d
		c.Log.Printf("[api] pausing %v until next request ...", d)
		time.Sleep(d)
	}

	if data, err = c.httpRequest(URL, client, method...); err != nil {
		return nil, err
	}
	c.lastRequest = time.Now()

	return data, nil
}

// retrieve URL with given client.
func (c *Client) httpRequest(URL string, client *http.Client, method ...string) ([]byte, error) {
	var (
		meth = "GET"
		req  *http.Request
		r    *http.Response
		data []byte
		err  error
	)
	if len(method) > 0 {
		meth = method[0]
	}
	c.Log.Printf("[http] retrieving %q ...", cleanURL(URL))

	if req, err = http.NewRequest(strings.ToUpper(meth), URL, nil); err != nil {
		return nil, errors.Wrap(err, "build HTTP request")
	}
	req.Header.Set("User-Agent", userAgent)

	if r, err = client.Do(req); err != nil {
		return nil, errors.Wrap(err, "retrieve URL")
	}
	defer r.Body.Close()
	c.Log.Printf("[%d] %s", r.StatusCode, cleanURL(URL))

	if r.StatusCode > 299 {
		return nil, errors.Wrap(fmt.Errorf("%s: %s", URL, r.Status), "retrieve URL")
	}

	if data, err = ioutil.ReadAll(r.Body); err != nil {
		return nil, errors.Wrap(err, "read HTTP response")
	}

	return data, nil
}

func cleanURL(URL string) string {
	if u, err := url.Parse(URL); err == nil {
		v := u.Query()
		if v.Get("key") != "" {
			v.Set("key", "xxx")
			u.RawQuery = v.Encode()
			return u.String()
		}
	}
	return URL
}
