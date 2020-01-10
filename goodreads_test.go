// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/fxtlabs/date"
	"github.com/stretchr/testify/assert"
)

// Books from search
func TestParseSearchResults(t *testing.T) {
	t.Parallel()

	rs, err := unmarshalSearchResults(readFile("dresden_files.xml", t))
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	assert.Equal(t, expectedDresden, rs, "unexpected Books")
}

// Books for an author
func TestParseAuthorBooks(t *testing.T) {
	t.Parallel()
	rs, meta, err := unmarshalAuthorBooks(readFile("jim_butcher.xml", t))
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	assert.Equal(t, 1, meta.Start, "unexpected meta.Start")
	assert.Equal(t, 30, meta.End, "unexpected meta.End")
	assert.Equal(t, 162, meta.Total, "unexpected meta.Total")
	assert.Equal(t, expectedButcher, rs, "unexpected Books")
}

func readFile(filename string, t *testing.T) []byte {
	data, err := ioutil.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatalf("read %s: %v", filename, err)
	}
	return data
}

var (
	expectedDresden = []Book{
		Book{
			ID:       47212,
			Title:    "Storm Front (The Dresden Files, #1)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2000, time.April, 1),
			Rating:   4.02,
			URL:      "https://www.goodreads.com/book/show/47212",
			ImageURL: "https://images.gr-assets.com/books/1419456275m/47212.jpg",
		},
		Book{
			ID:       6585201,
			Title:    "Changes (The Dresden Files, #12)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2010, time.April, 6),
			Rating:   4.54,
			URL:      "https://www.goodreads.com/book/show/6585201",
			ImageURL: "https://images.gr-assets.com/books/1304027244m/6585201.jpg",
		},
		Book{
			ID:       91477,
			Title:    "Fool Moon (The Dresden Files, #2)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2001, time.January, 1),
			Rating:   4.03,
			URL:      "https://www.goodreads.com/book/show/91477",
			ImageURL: "https://images.gr-assets.com/books/1507307616m/91477.jpg",
		},
		Book{
			ID:       7779059,
			Title:    "Side Jobs: Stories from the Dresden Files (The Dresden Files, #12.5)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2010, time.October, 26),
			Rating:   4.25,
			URL:      "https://www.goodreads.com/book/show/7779059",
			ImageURL: "https://images.gr-assets.com/books/1269115846m/7779059.jpg",
		},
		Book{
			ID:       91476,
			Title:    "Grave Peril (The Dresden Files, #3)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2001, time.September, 1),
			Rating:   4.18,
			URL:      "https://www.goodreads.com/book/show/91476",
			ImageURL: "https://images.gr-assets.com/books/1266470209m/91476.jpg",
		},
		Book{
			ID:       91478,
			Title:    "Summer Knight (The Dresden Files, #4)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2002, time.September, 3),
			Rating:   4.3,
			URL:      "https://www.goodreads.com/book/show/91478",
			ImageURL: "https://images.gr-assets.com/books/1345557469m/91478.jpg",
		},
		Book{
			ID:       91479,
			Title:    "Death Masks (The Dresden Files, #5)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2003, time.August, 1),
			Rating:   4.32,
			URL:      "https://www.goodreads.com/book/show/91479",
			ImageURL: "https://images.gr-assets.com/books/1345557713m/91479.jpg",
		},
		Book{
			ID:       99383,
			Title:    "Blood Rites (The Dresden Files, #6)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2004, time.August, 1),
			Rating:   4.33,
			URL:      "https://www.goodreads.com/book/show/99383",
			ImageURL: "https://images.gr-assets.com/books/1345557965m/99383.jpg",
		},
		Book{
			ID:       17683,
			Title:    "Dead Beat (The Dresden Files, #7)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2005, time.May, 3),
			Rating:   4.44,
			URL:      "https://www.goodreads.com/book/show/17683",
			ImageURL: "https://images.gr-assets.com/books/1345667776m/17683.jpg",
		},
		Book{
			ID:       91475,
			Title:    "White Night (The Dresden Files, #9)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2007, time.April, 3),
			Rating:   4.41,
			URL:      "https://www.goodreads.com/book/show/91475",
			ImageURL: "https://images.gr-assets.com/books/1309552288m/91475.jpg",
		},
		Book{
			ID:       91474,
			Title:    "Proven Guilty (The Dresden Files, #8)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2006, time.February, 1),
			Rating:   4.42,
			URL:      "https://www.goodreads.com/book/show/91474",
			ImageURL: "https://images.gr-assets.com/books/1345667469m/91474.jpg",
		},
		Book{
			ID:       927979,
			Title:    "Small Favor (The Dresden Files, #10)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2008, time.April, 1),
			Rating:   4.44,
			URL:      "https://www.goodreads.com/book/show/927979",
			ImageURL: "https://images.gr-assets.com/books/1298085176m/927979.jpg",
		},
		Book{
			ID:       3475161,
			Title:    "Turn Coat (The Dresden Files, #11)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2009, time.April, 7),
			Rating:   4.45,
			URL:      "https://www.goodreads.com/book/show/3475161",
			ImageURL: "https://images.gr-assets.com/books/1304027128m/3475161.jpg",
		},
		Book{
			ID:       8058301,
			Title:    "Ghost Story (The Dresden Files,  #13)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2011, time.January, 1),
			Rating:   4.25,
			URL:      "https://www.goodreads.com/book/show/8058301",
			ImageURL: "https://images.gr-assets.com/books/1329104700m/8058301.jpg",
		},
		Book{
			ID:       12216302,
			Title:    "Cold Days (The Dresden Files, #14)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2012, time.November, 27),
			Rating:   4.51,
			URL:      "https://www.goodreads.com/book/show/12216302",
			ImageURL: "https://images.gr-assets.com/books/1345145377m/12216302.jpg",
		},
		Book{
			ID:       19486421,
			Title:    "Skin Game (The Dresden Files, #15)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2014, time.May, 27),
			Rating:   4.56,
			URL:      "https://www.goodreads.com/book/show/19486421",
			ImageURL: "https://images.gr-assets.com/books/1387236318m/19486421.jpg",
		},
		Book{
			ID:       2575572,
			Title:    "Backup (The Dresden Files, #10.4)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2008, time.October, 31),
			Rating:   4.12,
			URL:      "https://www.goodreads.com/book/show/2575572",
			ImageURL: "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		},
		Book{
			ID:       22249640,
			Title:    "Peace Talks (The Dresden Files, #16)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.Date{},
			Rating:   4.44,
			URL:      "https://www.goodreads.com/book/show/22249640",
			ImageURL: "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		},
		Book{
			ID:       4271488,
			Title:    "Vignette (The Dresden Files, #5.5)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2008, 1, 1),
			Rating:   4.06,
			URL:      "https://www.goodreads.com/book/show/4271488",
			ImageURL: "https://images.gr-assets.com/books/1476229870m/4271488.jpg",
		},
		Book{
			ID:       12183815,
			Title:    "Brief Cases (The Dresden Files, #15.1)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2018, time.June, 5),
			Rating:   4.41,
			URL:      "https://www.goodreads.com/book/show/12183815",
			ImageURL: "https://images.gr-assets.com/books/1513644037m/12183815.jpg",
		},
	}

	expectedButcher = []Book{
		Book{
			ID:       47212,
			Title:    "Storm Front (The Dresden Files, #1)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2000, time.April, 1),
			Rating:   4.02,
			URL:      "https://www.goodreads.com/book/show/47212",
			ImageURL: "https://images.gr-assets.com/books/1419456275m/47212.jpg",
		},
		Book{
			ID:       91477,
			Title:    "Fool Moon (The Dresden Files, #2)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2001, time.January, 9),
			Rating:   4.03,
			URL:      "https://www.goodreads.com/book/show/91477",
			ImageURL: "https://images.gr-assets.com/books/1507307616m/91477.jpg",
		},
		Book{
			ID:       91476,
			Title:    "Grave Peril (The Dresden Files, #3)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2001, time.September, 4),
			Rating:   4.18,
			URL:      "https://www.goodreads.com/book/show/91476",
			ImageURL: "https://images.gr-assets.com/books/1266470209m/91476.jpg",
		},
		Book{
			ID:       91478,
			Title:    "Summer Knight (The Dresden Files, #4)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2002, time.September, 3),
			Rating:   4.3,
			URL:      "https://www.goodreads.com/book/show/91478",
			ImageURL: "https://images.gr-assets.com/books/1345557469m/91478.jpg",
		},
		Book{
			ID:       91479,
			Title:    "Death Masks (The Dresden Files, #5)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2003, time.August, 5),
			Rating:   4.32,
			URL:      "https://www.goodreads.com/book/show/91479",
			ImageURL: "https://images.gr-assets.com/books/1345557713m/91479.jpg",
		},
		Book{
			ID:       99383,
			Title:    "Blood Rites (The Dresden Files, #6)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2004, time.August, 3),
			Rating:   4.34,
			URL:      "https://www.goodreads.com/book/show/99383",
			ImageURL: "https://images.gr-assets.com/books/1345557965m/99383.jpg",
		},
		Book{
			ID:       17683,
			Title:    "Dead Beat (The Dresden Files, #7)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2006, time.January, 1),
			Rating:   4.44,
			URL:      "https://www.goodreads.com/book/show/17683",
			ImageURL: "https://images.gr-assets.com/books/1345667776m/17683.jpg",
		},
		Book{
			ID:       91475,
			Title:    "White Night (The Dresden Files, #9)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2007, time.April, 3),
			Rating:   4.41,
			URL:      "https://www.goodreads.com/book/show/91475",
			ImageURL: "https://images.gr-assets.com/books/1309552288m/91475.jpg",
		},
		Book{
			ID:       91474,
			Title:    "Proven Guilty (The Dresden Files, #8)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2007, time.February, 6),
			Rating:   4.42,
			URL:      "https://www.goodreads.com/book/show/91474",
			ImageURL: "https://images.gr-assets.com/books/1345667469m/91474.jpg",
		},
		Book{
			ID:       6585201,
			Title:    "Changes (The Dresden Files, #12)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2010, time.April, 6),
			Rating:   4.52,
			URL:      "https://www.goodreads.com/book/show/6585201",
			ImageURL: "https://images.gr-assets.com/books/1304027244m/6585201.jpg",
		},
		Book{
			ID:       927979,
			Title:    "Small Favor (The Dresden Files, #10)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2008, time.April, 1),
			Rating:   4.44,
			URL:      "https://www.goodreads.com/book/show/927979",
			ImageURL: "https://images.gr-assets.com/books/1298085176m/927979.jpg",
		},
		Book{
			ID:       29396,
			Title:    "Furies of Calderon (Codex Alera, #1)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2005, time.June, 28),
			Rating:   4.13,
			URL:      "https://www.goodreads.com/book/show/29396",
			ImageURL: "https://images.gr-assets.com/books/1329104514m/29396.jpg",
		},
		Book{
			ID:       3475161,
			Title:    "Turn Coat (The Dresden Files, #11)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2009, time.April, 7),
			Rating:   4.45,
			URL:      "https://www.goodreads.com/book/show/3475161",
			ImageURL: "https://images.gr-assets.com/books/1304027128m/3475161.jpg",
		},
		Book{
			ID:       12216302,
			Title:    "Cold Days (The Dresden Files, #14)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2012, time.November, 27),
			Rating:   4.51,
			URL:      "https://www.goodreads.com/book/show/12216302",
			ImageURL: "https://images.gr-assets.com/books/1345145377m/12216302.jpg",
		},
		Book{
			ID:       8058301,
			Title:    "Ghost Story (The Dresden Files,  #13)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2011, time.July, 26),
			Rating:   4.25,
			URL:      "https://www.goodreads.com/book/show/8058301",
			ImageURL: "https://images.gr-assets.com/books/1329104700m/8058301.jpg",
		},
		Book{
			ID:       19486421,
			Title:    "Skin Game (The Dresden Files, #15)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2014, time.May, 27),
			Rating:   4.56,
			URL:      "https://www.goodreads.com/book/show/19486421",
			ImageURL: "https://images.gr-assets.com/books/1387236318m/19486421.jpg",
		},
		Book{
			ID:       133664,
			Title:    "Academ's Fury (Codex Alera, #2)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2010, time.February, 1),
			Rating:   4.28,
			URL:      "https://www.goodreads.com/book/show/133664",
			ImageURL: "https://images.gr-assets.com/books/1381026900m/133664.jpg",
		},
		Book{
			ID:       346087,
			Title:    "Captain's Fury (Codex Alera, #4)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2007, time.December, 4),
			Rating:   4.38,
			URL:      "https://www.goodreads.com/book/show/346087",
			ImageURL: "https://images.gr-assets.com/books/1315083292m/346087.jpg",
		},
		Book{
			ID:       29394,
			Title:    "Cursor's Fury (Codex Alera, #3)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2006, time.December, 5),
			Rating:   4.37,
			URL:      "https://www.goodreads.com/book/show/29394",
			ImageURL: "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		},
		Book{
			ID:       2903736,
			Title:    "Princeps' Fury (Codex Alera, #5)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2008, time.November, 25),
			Rating:   4.36,
			URL:      "https://www.goodreads.com/book/show/2903736",
			ImageURL: "https://images.gr-assets.com/books/1315082776m/2903736.jpg",
		},
		Book{
			ID:       6316821,
			Title:    "First Lord's Fury (Codex Alera, #6)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2009, time.November, 24),
			Rating:   4.38,
			URL:      "https://www.goodreads.com/book/show/6316821",
			ImageURL: "https://images.gr-assets.com/books/1327903582m/6316821.jpg",
		},
		Book{
			ID:       7779059,
			Title:    "Side Jobs: Stories from the Dresden Files (The Dresden Files, #12.5)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2010, time.October, 26),
			Rating:   4.25,
			URL:      "https://www.goodreads.com/book/show/7779059",
			ImageURL: "https://images.gr-assets.com/books/1269115846m/7779059.jpg",
		},
		Book{
			ID:       24876258,
			Title:    "The Aeronaut's Windlass (The Cinder Spires, #1)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2015, time.September, 29),
			Rating:   4.18,
			URL:      "https://www.goodreads.com/book/show/24876258",
			ImageURL: "https://images.gr-assets.com/books/1425415066m/24876258.jpg",
		},
		Book{
			ID:       2637138,
			Title:    "Welcome to the Jungle",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2008, time.October, 21),
			Rating:   4.10,
			URL:      "https://www.goodreads.com/book/show/2637138",
			ImageURL: "https://images.gr-assets.com/books/1320418408m/2637138.jpg",
		},
		Book{
			ID:       4961959,
			Title:    "Jim Butcher's The Dresden Files: Storm Front, Volume 1: The Gathering Storm",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2009, time.June, 2),
			Rating:   4.36,
			URL:      "https://www.goodreads.com/book/show/4961959",
			ImageURL: "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		},
		Book{
			ID:       12183815,
			Title:    "Brief Cases (The Dresden Files, #15.1)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2018, time.June, 5),
			Rating:   4.42,
			URL:      "https://www.goodreads.com/book/show/12183815",
			ImageURL: "https://images.gr-assets.com/books/1513644037m/12183815.jpg",
		},
		Book{
			ID:       3475145,
			Title:    "Mean Streets",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2009, time.January, 1),
			Rating:   4.01,
			URL:      "https://www.goodreads.com/book/show/3475145",
			ImageURL: "https://images.gr-assets.com/books/1303858861m/3475145.jpg",
		},
		Book{
			ID:       2575572,
			Title:    "Backup (The Dresden Files, #10.4)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2008, time.October, 1),
			Rating:   4.12,
			URL:      "https://www.goodreads.com/book/show/2575572",
			ImageURL: "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		},
		Book{
			ID:       25807691,
			Title:    "Working for Bigfoot (The Dresden Files, #15.5)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.Date{},
			Rating:   4.26,
			URL:      "https://www.goodreads.com/book/show/25807691",
			ImageURL: "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		}, Book{
			ID:       15732549,
			Title:    "Restoration of Faith (The Dresden Files, #0.2)",
			Author:   Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:  date.New(2004, time.January, 1),
			Rating:   4.00,
			URL:      "https://www.goodreads.com/book/show/15732549",
			ImageURL: "https://images.gr-assets.com/books/1483723101m/15732549.jpg",
		},
	}
)
