// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package gr

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fxtlabs/date"
	"github.com/pkg/errors"
)

const (
	apiURL    = "https://www.goodreads.com/search/index.xml?q=%s"
	authorURL = "https://www.goodreads.com/author/list.xml?id=%d&page=%d"
	bookURL   = "https://www.goodreads.com/book/show/%d.xml?key=%s"
)

var (
	errEmptyQuery = errors.New("empty query")
)

// Book is an entry from an RSS or API feed.
// The two types of feeds contain different information, so not
// every field is always set.
type Book struct {
	ID     int64
	WorkID int64  // not in search results
	ISBN   string // not in search results
	ISBN13 string // not in search results

	Title         string
	TitleNoSeries string
	Series        Series
	Author        Author
	PubDate       date.Date
	Rating        float64 // average rating or user rating (RSS feeds)
	Description   string  // HTML, not in search results

	URL      string // Book's page on goodreads.com
	ImageURL string // URL of cover image
}

// HasSeries returns true if Book belongs to a series.
func (b Book) HasSeries() bool { return b.Series.Title != "" }

// DescriptionText return Book description as plaintext.
func (b Book) DescriptionText() string { return HTML2Text(b.Description) }

// DescriptionMarkdown return Book description as Markdown.
func (b Book) DescriptionMarkdown() string { return HTML2Markdown(b.Description) }

// String implements Stringer.
func (b Book) String() string {
	return fmt.Sprintf(`Book{ID: %d, Title: %q,	Series: %q, Author: %q}`, b.ID, b.TitleNoSeries, b.Series, b.Author)
}

// Data returns template data.
func (b Book) Data() map[string]string {
	data := map[string]string{
		"TITLE":                b.Title,
		"TITLE_NO_SERIES":      b.TitleNoSeries,
		"SERIES":               b.Series.Title,
		"SERIES_ID":            fmt.Sprintf("%d", b.Series.ID),
		"BOOK_ID":              fmt.Sprintf("%d", b.ID),
		"ISBN":                 b.ISBN,
		"WORK_ID":              fmt.Sprintf("%d", b.WorkID),
		"DESCRIPTION":          b.DescriptionText(),
		"DESCRIPTION_HTML":     b.Description,
		"DESCRIPTION_MARKDOWN": b.DescriptionMarkdown(),
		"AUTHOR":               b.Author.Name,
		"AUTHOR_ID":            fmt.Sprintf("%d", b.Author.ID),
		"AUTHOR_URL":           b.Author.URL,
		"YEAR":                 b.PubDate.Format("2006"),
		"RATING":               fmt.Sprintf("%f", b.Rating),
		"BOOK_URL":             b.URL,
		"IMAGE_URL":            b.ImageURL,
	}

	// remove empty/unset variabels
	out := map[string]string{}
	for k, v := range data {
		if v == "" || v == "0" {
			continue
		}
		out[k] = v
	}
	return out
}

// Author is the author of a book.
type Author struct {
	ID   int64  `xml:"id"` // not available in feeds
	Name string `xml:"name"`
	URL  string // not available in feeds
}

// String returns author's name.
func (a Author) String() string { return a.Name }

// Search API for books.
func (c *Client) Search(query string) (books []Book, err error) {
	var (
		u    = urlForQuery(url.QueryEscape(query))
		data []byte
	)
	if u == "" {
		err = errEmptyQuery
		return
	}
	if data, err = c.apiRequest(u); err != nil {
		return
	}

	return unmarshalSearchResults(data)
}

// BookDetails fetches the full details of a book.
func (c *Client) BookDetails(id int64) (Book, error) {
	var (
		u    = fmt.Sprintf(bookURL, id, c.APIKey)
		data []byte
		err  error
	)

	if data, err = c.apiRequest(u); err != nil {
		return Book{}, errors.Wrap(err, "fetch book details")
	}

	return unmarshalBookDetails(data)
}

func unmarshalBookDetails(data []byte) (Book, error) {
	v := struct {
		Book struct {
			ID     int64  `xml:"id"`
			WorkID int64  `xml:"work>id"`
			ISBN   string `xml:"isbn"`
			ISBN13 string `xml:"isbn13"`

			Title          string   `xml:"title"`
			TitleNoSeries  string   `xml:"work>original_title"`
			SeriesName     string   `xml:"series_works>series_work>series>title"`
			SeriesPosition float64  `xml:"series_works>series_work>user_position"`
			SeriesID       int64    `xml:"series_works>series_work>series>id"`
			Authors        []Author `xml:"authors>author"`
			Year           int      `xml:"work>original_publication_year"`
			Month          int      `xml:"work>original_publication_month"`
			Day            int      `xml:"work>original_publication_day"`
			Rating         float64  `xml:"average_rating"`
			Description    string   `xml:"description"`

			ImageURL string `xml:"image_url"`
		} `xml:"book"`
	}{}

	if err := xml.Unmarshal(data, &v); err != nil {
		log.Printf("[book] raw=%s", string(data))
		return Book{}, err
	}

	title := v.Book.TitleNoSeries
	if title == "" {
		title, _ = parseTitle(v.Book.Title)
	}

	b := Book{
		ID:            v.Book.ID,
		WorkID:        v.Book.WorkID,
		ISBN:          v.Book.ISBN,
		ISBN13:        v.Book.ISBN13,
		Title:         v.Book.Title,
		TitleNoSeries: title,
		Series: Series{
			Title:    strings.TrimSpace(v.Book.SeriesName),
			ID:       v.Book.SeriesID,
			Position: v.Book.SeriesPosition,
		},
		Rating:      v.Book.Rating,
		URL:         fmt.Sprintf("https://www.goodreads.com/book/show/%d", v.Book.ID),
		Description: v.Book.Description,
		ImageURL:    v.Book.ImageURL,
	}

	if len(v.Book.Authors) > 0 {
		b.Author = v.Book.Authors[0]
		b.Author.URL = fmt.Sprintf("https://www.goodreads.com/author/show/%d", b.Author.ID)
	}

	if v.Book.Month == 0 {
		v.Book.Month = 1
	}
	if v.Book.Day == 0 {
		v.Book.Day = 1
	}
	if v.Book.Year != 0 {
		b.PubDate = date.New(v.Book.Year, time.Month(v.Book.Month), v.Book.Day)
	}

	return b, nil
}

func urlForQuery(query string) string {
	if query == "" {
		return ""
	}
	return fmt.Sprintf(apiURL, query)
}

func unmarshalSearchResults(data []byte) (books []Book, err error) {
	v := struct {
		Works []struct {
			ID     int64  `xml:"best_book>id"`
			WorkID int64  `xml:"id"`
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
		title, series := parseTitle(r.Title)
		b := Book{
			ID:            r.ID,
			WorkID:        r.WorkID,
			Title:         r.Title,
			TitleNoSeries: title,
			Series:        series,
			Author:        r.Author,
			Rating:        r.Rating,
			URL:           fmt.Sprintf("https://www.goodreads.com/book/show/%d", r.ID),
			ImageURL:      r.ImageURL,
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

// PageData contains pagination data.
type PageData struct {
	Start int
	End   int
	Total int
}

func (pd PageData) String() string {
	return fmt.Sprintf("pageData{Start: %d, End: %d, Total: %d}", pd.Start, pd.End, pd.Total)
}

// AuthorBooks returns books for specified author.
func (c *Client) AuthorBooks(id int64, page int) (books []Book, meta PageData, err error) {
	if page == 0 {
		page = 1
	}
	var (
		u    = urlForAuthor(id, page)
		data []byte
	)
	if u == "" {
		err = errEmptyQuery
		return
	}

	if data, err = c.apiRequest(u); err != nil {
		return
	}

	return unmarshalAuthorBooks(data)
}

func urlForAuthor(id int64, page int) string {
	if page == 0 {
		page = 1
	}
	return fmt.Sprintf(authorURL, id, page)
}

func unmarshalAuthorBooks(data []byte) (books []Book, meta PageData, err error) {
	v := struct {
		List struct {
			Start int `xml:"start,attr"`
			End   int `xml:"end,attr"`
			Total int `xml:"total,attr"`

			Books []struct {
				ID            int64  `xml:"id"`
				ISBN          string `xml:"isbn"`
				ISBN13        string `xml:"isbn13"`
				Title         string `xml:"title"`
				TitleNoSeries string `xml:"title_without_series"`
				Description   string `xml:"description"`
				Year          int    `xml:"publication_year"`
				Month         int    `xml:"publication_month"`
				Day           int    `xml:"publication_day"`

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
		_, series := parseTitle(r.Title)
		b := Book{
			ID:            r.ID,
			ISBN:          r.ISBN,
			ISBN13:        r.ISBN13,
			Title:         r.Title,
			Series:        series,
			Description:   r.Description,
			TitleNoSeries: r.TitleNoSeries,
			Rating:        r.Rating,
			URL:           fmt.Sprintf("https://www.goodreads.com/book/show/%d", r.ID),
			ImageURL:      r.ImageURL,
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

// regexes to match book titles with embedded series info
var (
	seriesRegexes = []*regexp.Regexp{
		regexp.MustCompile(`^(.+)\s\((.+?),?\s+#([0-9.]+)\)$`),
		regexp.MustCompile(`^(.+)\s\((.+?) Series Book ([0-9.]+)\)$`),
	}
)

// extract title & series from book title with embedded series info.
func parseTitle(s string) (title string, series Series) {
	for _, rx := range seriesRegexes {
		values := rx.FindAllStringSubmatch(s, -1)
		if len(values) == 1 {
			title = strings.TrimSpace(values[0][1])
			series.Title = strings.TrimSpace(values[0][2])
			series.Position, _ = strconv.ParseFloat(values[0][3], 64)
			return
		}
	}

	title = s
	return
}
