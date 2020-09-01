// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package gr

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/fxtlabs/date"
	"github.com/pkg/errors"
)

const (
	shelfURL      = "https://www.goodreads.com/review/list.xml?v=2&id=%d&shelf=%s&page=%d&per_page=50&sort=position"
	shelvesURL    = "https://www.goodreads.com/shelf/list.xml?user_id=%d&page=%d"
	shelfAddURL   = "https://www.goodreads.com/shelf/add_to_shelf.xml"
	shelvesAddURL = "https://www.goodreads.com/shelf/add_books_to_shelves.xml"
)

// Shelf is a user's bookshelf/list.
type Shelf struct {
	ID       int64
	Name     string
	URL      string
	Size     int    // number of books on shelf
	Books    []Book // not populated in shelf list
	Selected bool
}

// String implements Stringer.
func (s Shelf) String() string {
	return fmt.Sprintf(`Shelf{ID: %d, Name: %q, Size: %d}`, s.ID, s.Name, s.Size)
}

// Title is the formatted Shelf name.
func (s Shelf) Title() string {
	switch s.Name {
	case "read":
		return "Read"
	case "currently-reading":
		return "Currently Reading"
	case "to-read":
		return "Want to Read"
	default:
		return s.Name
	}
}

// ShelvesByName sorts a slice of Shelf structs by name.
type ShelvesByName []Shelf

// Implement sort.Interface
func (s ShelvesByName) Len() int           { return len(s) }
func (s ShelvesByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ShelvesByName) Less(i, j int) bool { return s[i].Name < s[j].Name }

// UserShelf returns the books on the specified shelf.
func (c *Client) UserShelf(userID int64, name string, page int) ([]Book, PageData, error) {
	if page == 0 {
		page = 1
	}

	var (
		u    = fmt.Sprintf(shelfURL, userID, name, page)
		data []byte
		err  error
	)

	if data, err = c.apiRequest(u); err != nil {
		return nil, PageData{}, errors.Wrap(err, "fetch shelf")
	}

	return unmarshalShelf(data)
}

func unmarshalShelf(data []byte) ([]Book, PageData, error) {
	var (
		books []Book
		meta  PageData
		err   error
	)
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
			} `xml:"review>book"`
			XMLName xml.Name `xml:"reviews"`
		}
	}{}
	if err = xml.Unmarshal(data, &v); err != nil {
		return nil, PageData{}, errors.Wrap(err, "unmarshal shelf")
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
			TitleNoSeries: r.TitleNoSeries,
			Series:        series,
			Description:   r.Description,
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

	return books, meta, nil
}

// UserShelves retrieve user's shelves. Only basic shelf info (ID, name, book count) is returned.
func (c *Client) UserShelves(userID int64, page int) (shelves []Shelf, meta PageData, err error) {
	if page == 0 {
		page = 1
	}

	var (
		u    = fmt.Sprintf(shelvesURL, userID, page)
		data []byte
	)

	if data, err = c.apiRequest(u); err != nil {
		return
	}

	if shelves, meta, err = unmarshalShelves(data); err == nil {
		for i, s := range shelves {
			URL := fmt.Sprintf("https://www.goodreads.com/review/list/%d", userID)
			u, _ := url.Parse(URL)
			v := u.Query()
			v.Set("shelf", s.Name)
			u.RawQuery = v.Encode()
			s.URL = u.String()
			shelves[i] = s
		}
	}
	return
}

func unmarshalShelves(data []byte) ([]Shelf, PageData, error) {
	var (
		shelves []Shelf
		meta    PageData
		err     error
	)
	v := struct {
		List struct {
			Start   int `xml:"start,attr"`
			End     int `xml:"end,attr"`
			Total   int `xml:"total,attr"`
			Shelves []struct {
				ID        int64  `xml:"id"`
				Name      string `xml:"name"`
				BookCount int    `xml:"book_count"`
			} `xml:"user_shelf"`
			XMLName xml.Name `xml:"shelves"`
		}
	}{}
	if err = xml.Unmarshal(data, &v); err != nil {
		return nil, PageData{}, errors.Wrap(err, "parse shelves data")
	}

	meta.Start = v.List.Start
	meta.End = v.List.End
	meta.Total = v.List.Total

	for _, r := range v.List.Shelves {
		s := Shelf{
			ID:   r.ID,
			Name: r.Name,
			Size: r.BookCount,
		}
		shelves = append(shelves, s)
	}

	return shelves, meta, nil
}

// AddToShelves adds a book to the specified shelves.
func (c *Client) AddToShelves(bookID int64, shelves []string) error {
	var (
		u, _ = url.Parse(shelvesAddURL)
		v    = u.Query()
	)
	v.Set("shelves", strings.Join(shelves, ","))
	v.Set("bookids", fmt.Sprintf("%d", bookID))
	u.RawQuery = v.Encode()

	if _, err := c.apiRequest(u.String(), "POST"); err != nil {
		return err
	}

	return nil
}

// AddToShelf adds a book to the specified shelf.
func (c *Client) AddToShelf(bookID int64, shelf string) error {
	if err := c.addRemoveShelf(bookID, shelf, false); err != nil {
		return errors.Wrap(err, "add to shelf")
	}
	return nil
}

// RemoveFromShelf removes a book from the specified shelf.
func (c *Client) RemoveFromShelf(bookID int64, shelf string) error {
	if err := c.addRemoveShelf(bookID, shelf, true); err != nil {
		return errors.Wrap(err, "remove from shelf")
	}
	return nil
}

func (c *Client) addRemoveShelf(bookID int64, shelf string, remove bool) error {
	var (
		u, _ = url.Parse(shelfAddURL)
		v    = u.Query()
	)
	v.Set("name", shelf)
	v.Set("book_id", fmt.Sprintf("%d", bookID))
	if remove {
		v.Set("a", "remove")
	}
	u.RawQuery = v.Encode()

	if _, err := c.apiRequest(u.String(), "POST"); err != nil {
		return err
	}

	return nil
}
