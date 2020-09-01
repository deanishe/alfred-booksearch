// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package gr

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// Base URL of RSS feeds.
const rssURL = "https://www.goodreads.com/review/list_rss/"

// Feed is a Goodreads RSS feed. It's only used to retrieve cover images
// (a lot of covers are missing from API responses), so the Books contained
// by a Feed only have ID and ImageURL set.
type Feed struct {
	Name  string
	Books []Book
}

// Extract user ID and user feed key from a feed URL. Currently unused.
func parseFeedURL(URL string) (userID, feedKey string, err error) {
	var u *url.URL
	if u, err = url.Parse(URL); err != nil {
		return
	}

	userID = path.Base(u.Path)
	feedKey = u.Query().Get("key")
	return
}

// FetchFeed retrieves and parses a Goodreads RSS feed.
func (c *Client) FetchFeed(userID int64, shelf string) (Feed, error) {
	var (
		data []byte
		err  error
	)

	u, _ := url.Parse(fmt.Sprintf("%s%d", rssURL, userID))
	v := u.Query()
	v.Set("shelf", shelf)
	u.RawQuery = v.Encode()

	if data, err = c.httpGet(u.String()); err != nil {
		return Feed{}, errors.Wrap(err, "retrive feed")
	}

	return unmarshalFeed(data)
}

// Parse RSS feed data.
func unmarshalFeed(data []byte) (Feed, error) {
	var (
		feed Feed
		err  error
	)
	v := struct {
		Name  string `xml:"channel>title"`
		Items []struct {
			ID             int64  `xml:"book_id"`
			ImageURL       string `xml:"book_image_url"`
			ImageURLMedium string `xml:"book_medium_image_url"`
			ImageURLLarge  string `xml:"book_large_image_url"`
		} `xml:"channel>item"`
	}{}

	if err = xml.Unmarshal(data, &v); err != nil {
		return Feed{}, errors.Wrap(err, "unmarshal feed")
	}

	feed.Name = parseFeedTitle(v.Name)

	for _, r := range v.Items {
		b := Book{
			ID:       r.ID,
			ImageURL: r.ImageURL,
		}

		if r.ImageURLLarge != "" {
			b.ImageURL = r.ImageURLLarge
		} else if r.ImageURLMedium != "" {
			b.ImageURL = r.ImageURLMedium
		}

		feed.Books = append(feed.Books, b)
	}

	return feed, nil
}

func parseFeedTitle(s string) string {
	i := strings.Index(s, "bookshelf: ")
	if i < 0 {
		return s
	}
	return s[i+11:]
}
