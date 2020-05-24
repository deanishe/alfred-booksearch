// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-08-08

package gr

import (
	"testing"
	"time"

	"github.com/fxtlabs/date"
	"github.com/stretchr/testify/assert"
)

// TestParseSeries parses series XML
func TestParseSeries(t *testing.T) {
	t.Parallel()

	series, err := unmarshalSeries(readFile("alex_verus.xml", t))
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	assert.Equal(t, expectedVerus, series)
}

var expectedVerus = Series{
	Title: "Alex Verus",
	ID:    71196,
	Books: []Book{
		{
			ID:            11737387,
			WorkID:        16686573,
			Title:         "Fated (Alex Verus, #1)",
			TitleNoSeries: "Fated",
			Series:        Series{ID: 71196, Title: "Alex Verus", Position: 1},
			Author:        Author{ID: 849723, Name: "Benedict Jacka", URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2012, time.February, 1),
			URL:           "https://www.goodreads.com/book/show/11737387",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1330906653l/11737387._SX98_.jpg",
		},
		{
			ID:            13274082,
			WorkID:        18478073,
			Title:         "Cursed (Alex Verus, #2)",
			TitleNoSeries: "Cursed",
			Series:        Series{ID: 71196, Title: "Alex Verus", Position: 2},
			Author:        Author{ID: 849723, Name: "Benedict Jacka", URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2012, time.May, 29),
			URL:           "https://www.goodreads.com/book/show/13274082",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1330971845l/13274082._SX98_.jpg",
		},
		{
			ID:            13542616,
			WorkID:        19062926,
			Title:         "Taken (Alex Verus, #3)",
			TitleNoSeries: "Taken",
			Series:        Series{ID: 71196, Title: "Alex Verus", Position: 3},
			Author:        Author{ID: 849723, Name: "Benedict Jacka", URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2012, time.August, 28),
			URL:           "https://www.goodreads.com/book/show/13542616",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1346617379l/13542616._SX98_.jpg",
		},
		{
			ID:            16072988,
			WorkID:        21867224,
			Title:         "Chosen (Alex Verus, #4)",
			TitleNoSeries: "Chosen",
			Series:        Series{ID: 71196, Title: "Alex Verus", Position: 4},
			Author:        Author{ID: 849723, Name: "Benedict Jacka", URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2013, time.August, 27),
			URL:           "https://www.goodreads.com/book/show/16072988",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1365983616l/16072988._SX98_.jpg",
		},
		{
			ID:            18599601,
			WorkID:        26365716,
			Title:         "Hidden (Alex Verus, #5)",
			TitleNoSeries: "Hidden",
			Series:        Series{ID: 71196, Title: "Alex Verus", Position: 5},
			Author:        Author{ID: 849723, Name: "Benedict Jacka", URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2014, time.September, 2),
			URL:           "https://www.goodreads.com/book/show/18599601",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1386933935l/18599601._SX98_.jpg",
		},
		{
			ID:            23236738,
			WorkID:        42780952,
			Title:         "Veiled (Alex Verus, #6)",
			TitleNoSeries: "Veiled",
			Series:        Series{ID: 71196, Title: "Alex Verus", Position: 6},
			Author:        Author{ID: 849723, Name: "Benedict Jacka", URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2015, time.August, 4),
			URL:           "https://www.goodreads.com/book/show/23236738",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1421862439l/23236738._SX98_.jpg",
		},
		{
			ID:            23236743,
			WorkID:        42780954,
			Title:         "Burned (Alex Verus, #7)",
			TitleNoSeries: "Burned",
			Series:        Series{ID: 71196, Title: "Alex Verus", Position: 7},
			Author:        Author{ID: 849723, Name: "Benedict Jacka", URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2016, time.April, 5),
			URL:           "https://www.goodreads.com/book/show/23236743",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1453058973l/23236743._SX98_.jpg",
		},
		{
			ID:            29865319,
			WorkID:        76176665,
			Title:         "Bound (Alex Verus, #8)",
			TitleNoSeries: "Bound",
			Series:        Series{ID: 71196, Title: "Alex Verus", Position: 8},
			Author:        Author{ID: 849723, Name: "Benedict Jacka", URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2017, time.April, 4),
			URL:           "https://www.goodreads.com/book/show/29865319",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1474725377l/29865319._SX98_.jpg",
		},
		{
			ID:            36068567,
			WorkID:        57651776,
			Title:         "Marked (Alex Verus, #9)",
			TitleNoSeries: "Marked",
			Series:        Series{ID: 71196, Title: "Alex Verus", Position: 9},
			Author:        Author{ID: 849723, Name: "Benedict Jacka", URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2018, time.July, 3),
			URL:           "https://www.goodreads.com/book/show/36068567",
			ImageURL:      "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		},
		{
			ID:            43670629,
			WorkID:        63342276,
			Title:         "Fallen (Alex Verus, #10)",
			TitleNoSeries: "Fallen",
			Series:        Series{ID: 71196, Title: "Alex Verus", Position: 10},
			Author:        Author{ID: 849723, Name: "Benedict Jacka", URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2019, time.September, 24),
			URL:           "https://www.goodreads.com/book/show/43670629",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1554395647l/43670629._SX98_.jpg",
		},
		{
			ID:            50740363,
			WorkID:        75767304,
			Title:         "Forged (Alex Verus, #11)",
			TitleNoSeries: "Forged",
			Series:        Series{ID: 71196, Title: "Alex Verus", Position: 11},
			Author:        Author{ID: 849723, Name: "Benedict Jacka", URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2020, time.November, 24),
			URL:           "https://www.goodreads.com/book/show/50740363",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1591956617l/50740363._SX98_.jpg",
		},
		{
			ID:            20980464,
			WorkID:        40357662,
			Title:         "The Alex Verus Novels, Books 1-4",
			TitleNoSeries: "The Alex Verus Novels, Books 1-4",
			Series:        Series{ID: 71196, Title: "Alex Verus", Position: 0},
			Author:        Author{ID: 849723, Name: "Benedict Jacka", URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2014, time.March, 4),
			URL:           "https://www.goodreads.com/book/show/20980464",
			ImageURL:      "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		},
	},
}
