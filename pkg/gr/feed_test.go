// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package gr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFeed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		books []Book
	}{
		{"to-read", expectedToRead},
		{"fantasy", expectedFantasy},
	}

	for _, td := range tests {
		td := td
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			feed, err := unmarshalFeed(readFile(td.name+".rss", t))
			require.Nil(t, err, "unmarshal feed %s.xml", td.name)
			assert.Equal(t, td.name, feed.Name)
			assert.Equal(t, td.books, feed.Books)
		})
	}
}

func TestParseFeedURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		URL, userID, feedKey string
	}{
		{"https://www.goodreads.com/review/list_rss/123456?key=7yg3Z3aVn-TWgH8Q_GGZDx&shelf=%23ALL%23", "123456", "7yg3Z3aVn-TWgH8Q_GGZDx"},
		{"https://www.goodreads.com/review/list_rss/7220456?key=eHfwY8fE_unud_BMzPT-uj2&shelf=to-read", "7220456", "eHfwY8fE_unud_BMzPT-uj2"},
	}

	for _, td := range tests {
		td := td
		t.Run(td.URL, func(t *testing.T) {
			t.Parallel()
			uid, key, err := parseFeedURL(td.URL)
			assert.Nil(t, err, "parse feed URL %q", td.URL)
			assert.Equal(t, td.userID, uid, "unexpected user_id")
			assert.Equal(t, td.feedKey, key, "unexpected feed_key")
		})
	}
}

var (
	expectedToRead = []Book{
		{
			ID:       109502,
			ImageURL: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1410138674l/109502.jpg",
		},
		{
			ID:       22477307,
			ImageURL: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1418771663l/22477307.jpg",
		},
		{
			ID:       25135194,
			ImageURL: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1432827094l/25135194._SY475_.jpg",
		},
		{
			ID:       50740363,
			ImageURL: "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		},
		{
			ID:       1268479,
			ImageURL: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1240256182l/1268479.jpg",
		},
		{
			ID:       42592353,
			ImageURL: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1562549322l/42592353._SY475_.jpg",
		},
	}

	expectedFantasy = []Book{
		{
			ID:       32337902,
			ImageURL: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1481987017l/32337902.jpg",
		},
		{
			ID:       16096968,
			ImageURL: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1487946539l/16096968.jpg",
		},
		{
			ID:       26030742,
			ImageURL: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1456696097l/26030742._SY475_.jpg",
		},
		{
			ID:       37769892,
			ImageURL: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1514589702l/37769892._SY475_.jpg",
		},
	}
)
