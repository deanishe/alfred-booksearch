// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-08-08

package gr

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fxtlabs/date"
	"github.com/pkg/errors"
)

const seriesURL = "https://www.goodreads.com/series/%d?format=xml&key=%s"

// Series is a Goodreads series.
type Series struct {
	Title    string  // series title
	Position float64 // position of current book in series
	ID       int64   // only set in book details
	Books    []Book  // only set by series endpoint
}

// String returns series name and book position.
func (s Series) String() string {
	if s.Title == "" {
		return ""
	}
	return fmt.Sprintf("%s #%.1f", s.Title, s.Position)
}

// Series fetches the full details of a book.
func (c *Client) Series(id int64) (Series, error) {
	var (
		u    = fmt.Sprintf(seriesURL, id, c.APIKey)
		data []byte
		err  error
	)

	if data, err = c.apiRequest(u); err != nil {
		return Series{}, errors.Wrap(err, "fetch series")
	}

	return unmarshalSeries(data)
}

func unmarshalSeries(data []byte) (Series, error) {
	v := struct {
		Series struct {
			ID          int64  `xml:"id"`
			Title       string `xml:"title"`
			Description string `xml:"description"`
			Works       []struct {
				WorkID        int64   `xml:"work>id"`
				BookID        int64   `xml:"work>best_book>id"`
				Title         string  `xml:"work>best_book>title"`
				TitleNoSeries string  `xml:"work>original_title"`
				AuthorID      int64   `xml:"work>best_book>author>id"`
				AuthorName    string  `xml:"work>best_book>author>name"`
				ImageURL      string  `xml:"work>best_book>image_url"`
				Year          int     `xml:"work>original_publication_year"`
				Month         int     `xml:"work>original_publication_month"`
				Day           int     `xml:"work>original_publication_day"`
				Rating        float64 `xml:"work>average_rating"`
				Position      string  `xml:"user_position"`
			} `xml:"series_works>series_work"`
		} `xml:"series"`
	}{}

	if err := xml.Unmarshal(data, &v); err != nil {
		return Series{}, err
	}
	series := Series{
		ID:    v.Series.ID,
		Title: strings.TrimSpace(v.Series.Title),
	}

	for _, w := range v.Series.Works {
		var pos float64
		if f, err := strconv.ParseFloat(w.Position, 64); err == nil {
			pos = f
		}
		b := Book{
			ID:            w.BookID,
			WorkID:        w.WorkID,
			Title:         strings.TrimSpace(w.Title),
			TitleNoSeries: strings.TrimSpace(w.TitleNoSeries),
			Series:        Series{ID: series.ID, Title: series.Title, Position: pos},
			Author:        Author{ID: w.AuthorID, Name: w.AuthorName, URL: fmt.Sprintf("https://www.goodreads.com/author/show/%d", w.AuthorID)},
			Rating:        w.Rating,
			URL:           fmt.Sprintf("https://www.goodreads.com/book/show/%d", w.BookID),
			ImageURL:      w.ImageURL,
		}

		if b.TitleNoSeries == "" {
			b.TitleNoSeries, _ = parseTitle(b.Title)
		}

		if w.Month == 0 {
			w.Month = 1
		}
		if w.Day == 0 {
			w.Day = 1
		}
		if w.Year != 0 {
			b.PubDate = date.New(w.Year, time.Month(w.Month), w.Day)
		}

		series.Books = append(series.Books, b)
	}

	return series, nil
}
