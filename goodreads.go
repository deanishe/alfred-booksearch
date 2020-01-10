// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/fxtlabs/date"
)

const (
	apiURL    = "https://www.goodreads.com/search/index.xml?key=%s&q=%s"
	authorURL = "https://www.goodreads.com/author/list.xml?key=%s&id=%s&page=%d"
)

var (
	errEmptyQuery = errors.New("empty query")
	timeout       = 10 * time.Second
	client        *http.Client
)

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: timeout,
			}).Dial,
			TLSHandshakeTimeout:   timeout,
			ResponseHeaderTimeout: timeout,
			ExpectContinueTimeout: timeout,
		},
	}
}

type Author struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
	URL  string
}

// Eq checks whether two Authors have the same values.
func (a Author) Eq(other Author) bool {
	if a.ID != other.ID {
		log.Printf("[Author.ID] %d != %d", a.ID, other.ID)
		return false
	}
	if a.Name != other.Name {
		log.Printf("[Author.Name] %q != %q", a.Name, other.Name)
		return false
	}
	if a.URL != other.URL {
		log.Printf("[Author.URL] %q != %q", a.URL, other.URL)
		return false
	}
	return true
}

type Book struct {
	ID int

	Title   string
	Author  Author
	PubDate date.Date
	Rating  float64

	URL      string
	ImageURL string
}

// String implements Stringer.
func (b Book) String() string {
	return fmt.Sprintf(`Book{
	ID: %d,

	Title:   %q,
	Author:  %v,
	PubDate: %v,
	Rating:  %0.2f,

	URL:      %q,
	ImageURL: %q,
}`, b.ID, b.Title, b.Author, b.PubDate, b.Rating, b.URL, b.ImageURL)
}

// Eq checks whether two Books have the same values.
func (b Book) Eq(other Book) bool {
	if b.ID != other.ID {
		log.Printf("[ID] %d != %d", b.ID, other.ID)
		return false
	}
	if b.Title != other.Title {
		log.Printf("[Title] %q != %q", b.Title, other.Title)
		return false
	}
	if !b.Author.Eq(other.Author) {
		return false
	}
	if (b.PubDate.IsZero() && !other.PubDate.IsZero()) || (!b.PubDate.IsZero() && other.PubDate.IsZero()) {
		return false
	}
	if !b.PubDate.IsZero() {
		if !b.PubDate.Equal(other.PubDate) {
			log.Printf("[PubDate] %v != %v", b.PubDate, other.PubDate)
			return false
		}
	}
	if b.Rating != other.Rating {
		log.Printf("[Rating] %f != %f", b.Rating, other.Rating)
		return false
	}
	if b.URL != other.URL {
		log.Printf("[URL] %q != %q", b.URL, other.URL)
		return false
	}
	if b.ImageURL != other.ImageURL {
		log.Printf("[ImageURL] %q != %q", b.ImageURL, other.ImageURL)
		return false
	}
	return true
}

// Data returns template data.
func (b Book) Data() map[string]string {
	titleNoSeries := trimSeries(b.Title)
	return map[string]string{
		"Title":               b.Title,
		"TitleNoSeries":       titleNoSeries,
		"Author":              b.Author.Name,
		"AuthorID":            fmt.Sprintf("%d", b.Author.ID),
		"AuthorURL":           b.Author.URL,
		"AuthorTitle":         b.Author.Name + " " + b.Title,
		"AuthorTitleNoSeries": b.Author.Name + " " + titleNoSeries,
		"Year":                b.PubDate.Format("2006"),
		"Rating":              fmt.Sprintf("%f", b.Rating),
		"URL":                 b.URL,
		"ImageURL":            b.ImageURL,
	}
}

func urlForQuery(query, apiKey string) string {
	if query == "" {
		return ""
	}
	return fmt.Sprintf(apiURL, apiKey, query)
}

func urlForAuthor(id, apiKey string, page int) string {
	if page == 0 {
		page = 1
	}
	return fmt.Sprintf(authorURL, apiKey, id, page)
}

func search(query, apiKey string) (books []Book, err error) {
	var (
		u    = urlForQuery(url.QueryEscape(query), url.QueryEscape(apiKey))
		r    *http.Response
		data []byte
	)
	if u == "" {
		err = errEmptyQuery
		return
	}
	log.Printf("retrieving %q ...", strings.Replace(u, apiKey, "XXXXX", 1))
	if r, err = client.Get(u); err != nil {
		return
	}
	defer r.Body.Close()
	if r.StatusCode > 299 {
		err = errors.New(r.Status)
		return
	}

	if data, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}

	return unmarshalSearchResults(data)
}

type pageData struct {
	Start int
	End   int
	Total int
}

func authorBooks(id, apiKey string, page int) (books []Book, meta pageData, err error) {
	if page == 0 {
		page = 1
	}
	var (
		u    = urlForAuthor(url.QueryEscape(id), url.QueryEscape(apiKey), page)
		r    *http.Response
		data []byte
	)
	if u == "" {
		err = errEmptyQuery
		return
	}
	log.Printf("retrieving %q ...", strings.Replace(u, apiKey, "XXXXX", 1))
	if r, err = client.Get(u); err != nil {
		return
	}
	defer r.Body.Close()
	if r.StatusCode > 299 {
		err = errors.New(r.Status)
		return
	}

	if data, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}

	return unmarshalAuthorBooks(data)
}

func unmarshalSearchResults(data []byte) (books []Book, err error) {
	v := struct {
		Works []struct {
			ID     int    `xml:"best_book>id"`
			Title  string `xml:"best_book>title"`
			Author Author `xml:"best_book>author"`
			Year   int    `xml:"original_publication_year"`
			Month  int    `xml:"original_publication_month"`
			Day    int    `xml:"original_publication_day"`

			Rating   float64 `xml:"average_rating"`
			ImageURL string  `xml:"best_book>image_url"`
		} `xml:"search>results>work"`
	}{}
	if err = xml.Unmarshal(data, &v); err != nil {
		return
	}

	for _, r := range v.Works {
		b := Book{
			ID:       r.ID,
			Title:    r.Title,
			Author:   r.Author,
			Rating:   r.Rating,
			URL:      fmt.Sprintf("https://www.goodreads.com/book/show/%d", r.ID),
			ImageURL: r.ImageURL,
		}
		if r.Month == 0 {
			r.Month = 1
		}
		if r.Day == 0 {
			r.Day = 1
		}

		if r.Year != 0 {
			b.PubDate = date.New(r.Year, time.Month(r.Month), r.Day)
		}
		b.Author.URL = fmt.Sprintf("https://www.goodreads.com/author/show/%d", b.Author.ID)
		books = append(books, b)
	}

	return
}

func unmarshalAuthorBooks(data []byte) (books []Book, meta pageData, err error) {
	v := struct {
		List struct {
			Start int `xml:"start,attr"`
			End   int `xml:"end,attr"`
			Total int `xml:"total,attr"`

			Books []struct {
				ID    int    `xml:"id"`
				Title string `xml:"title"`
				Year  int    `xml:"publication_year"`
				Month int    `xml:"publication_month"`
				Day   int    `xml:"publication_day"`

				Authors []Author `xml:"authors>author"`

				Rating   float64 `xml:"average_rating"`
				ImageURL string  `xml:"image_url"`
			} `xml:"book"`
		} `xml:"author>books"`
	}{}
	if err = xml.Unmarshal(data, &v); err != nil {
		return
	}

	meta.Start = v.List.Start
	meta.End = v.List.End
	meta.Total = v.List.Total

	for _, r := range v.List.Books {
		b := Book{
			ID:       r.ID,
			Title:    r.Title,
			Rating:   r.Rating,
			URL:      fmt.Sprintf("https://www.goodreads.com/book/show/%d", r.ID),
			ImageURL: r.ImageURL,
		}

		if len(r.Authors) > 0 {
			b.Author = r.Authors[0]
			b.Author.URL = fmt.Sprintf("https://www.goodreads.com/author/show/%d", b.Author.ID)
		}

		if r.Month == 0 {
			r.Month = 1
		}
		if r.Day == 0 {
			r.Day = 1
		}

		if r.Year != 0 {
			b.PubDate = date.New(r.Year, time.Month(r.Month), r.Day)
		}
		books = append(books, b)
	}

	return
}

var seriesRegex = regexp.MustCompile(`^(.+)\s\(.+?#\d+\)$`)

// remove series from book title
func trimSeries(title string) string {
	values := seriesRegex.FindAllStringSubmatch(title, -1)
	if values == nil || len(values) == 0 {
		return title
	}

	return values[0][1]
}
