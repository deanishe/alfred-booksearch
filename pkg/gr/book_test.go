// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package gr

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/fxtlabs/date"
	"github.com/stretchr/testify/assert"
)

// TestParseSearchResults parse search results
func TestParseSearchResults(t *testing.T) {
	t.Parallel()

	books, err := unmarshalSearchResults(readFile("dresden_files.xml", t))
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	assert.Equal(t, expectedDresden, books, "unexpected Books")
}

// TestParseAuthorBooks parse list of author's books
func TestParseAuthorBooks(t *testing.T) {
	t.Parallel()
	books, meta, err := unmarshalAuthorBooks(readFile("jim_butcher.xml", t))
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	assert.Equal(t, 1, meta.Start, "unexpected meta.Start")
	assert.Equal(t, 30, meta.End, "unexpected meta.End")
	assert.Equal(t, 162, meta.Total, "unexpected meta.Total")
	assert.Equal(t, expectedButcher, books, "unexpected Books")
}

// TestParseBooks parses book details
func TestParseBooks(t *testing.T) {
	t.Parallel()
	tests := []string{"forged.xml", "shockwave.xml"}

	for i, filename := range tests {
		i, filename := i, filename
		t.Run(filename, func(t *testing.T) {
			book, err := unmarshalBookDetails(readFile(filename, t))
			if err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			assert.Equal(t, books[i], book, "unexpected book")
		})
	}
}

// TestParseTitle parses book titles into title + series
func TestParseTitle(t *testing.T) {
	t.Parallel()
	tests := []struct {
		s      string
		title  string
		series Series
	}{
		{"Storm Front (The Dresden Files, #1)",
			"Storm Front",
			Series{Title: "The Dresden Files", Position: 1}},
		{"Side Jobs: Stories from the Dresden Files (The Dresden Files, #12.5)",
			"Side Jobs: Stories from the Dresden Files",
			Series{Title: "The Dresden Files", Position: 12.5}},
		{"Glass World (Undying Mercenaries Series Book 13)",
			"Glass World",
			Series{Title: "Undying Mercenaries", Position: 13}},
		{"Shockwave (Star Kingdom #1)",
			"Shockwave",
			Series{Title: "Star Kingdom", Position: 1}},
	}

	for _, td := range tests {
		td := td
		t.Run(td.s, func(t *testing.T) {
			title, series := parseTitle(td.s)
			assert.Equal(t, td.title, title, "unexpected title: %q", title)
			assert.Equal(t, td.series, series, "unexpected series: %v", series)
		})
	}
}

func readFile(filename string, t *testing.T) []byte {
	data, err := ioutil.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatalf("read %s: %v", filename, err)
	}
	return data
}

var (
	books = []Book{
		{
			ID:            50740363,
			WorkID:        75767304,
			ISBN:          "0356511146",
			ISBN13:        "9780356511146",
			Title:         "Forged (Alex Verus, #11)",
			TitleNoSeries: "Forged",
			Series:        Series{Title: "Alex Verus", Position: 11, ID: 71196},
			Author:        Author{Name: "Benedict Jacka", ID: 849723, URL: "https://www.goodreads.com/author/show/849723"},
			PubDate:       date.New(2020, time.November, 24),
			Rating:        4.3,
			Description:   `Alex Verus faces his dark side in this return to the bestselling urban fantasy series about a London-based mage.<br /><br />To protect his friends, Mage Alex Verus has had to change--and embrace his dark side. But the life mage Anne has changed too, and made a bond with a dangerous power. She's going after everyone she's got a grudge against--and it's a long list.<br /><br />In the meantime, Alex has to deal with his arch-enemy, Levistus. The Council's death squads are hunting Alex as well as Anne, and the only way for Alex to stop them is to end his long war with Levistus and the Council, by whatever means necessary. It will take everything Alex has to stay a step ahead of the Council and stop Anne from letting the world burn.`,
			URL:           "https://www.goodreads.com/book/show/50740363",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1591956617l/50740363._SX98_.jpg",
		},
		{
			ID:            45353889,
			WorkID:        70095418,
			ISBN:          "",
			ISBN13:        "",
			Title:         "Shockwave (Star Kingdom #1)",
			TitleNoSeries: "Shockwave",
			Series:        Series{Title: "Star Kingdom", Position: 1, ID: 261857},
			Author:        Author{Name: "Lindsay Buroker", ID: 4512224, URL: "https://www.goodreads.com/author/show/4512224"},
			PubDate:       date.New(2019, time.May, 8),
			Rating:        4.18,
			Description:   `<b>What if being a hero was encoded in your genes?<br /><br />And nobody told you?</b><br /><br />Casmir Dabrowski would laugh if someone asked him that. After all, he had to build a robot to protect himself from bullies when he was in school.<br /><br />Fortunately, life is a little better these days. He's an accomplished robotics engineer, a respected professor, and he almost never gets picked on in the lunchroom. But he's positive heroics are for other people.<br /><br />Until robot assassins stride onto campus and try to kill him.<br /><br />Forced to flee the work he loves and the only home he's ever known, Casmir catches the first ship into space, where he hopes to buy time to figure out who wants him dead and why. If he can't, he'll never be able to return home.<br /><br />But he soon finds himself entangled with bounty hunters, mercenaries, and pirates, including the most feared criminal in the Star Kingdom: Captain Tenebris Rache.<br /><br />Rache could snap his spine with one cybernetically enhanced finger, but he may be the only person with the answer Casmir desperately needs:<br /><br />What in his genes is worth killing for?`,
			URL:           "https://www.goodreads.com/book/show/45353889",
			ImageURL:      "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		},
	}

	expectedDresden = []Book{
		{
			ID:            47212,
			WorkID:        1137060,
			Title:         "Storm Front (The Dresden Files, #1)",
			TitleNoSeries: "Storm Front",
			Series:        Series{Title: "The Dresden Files", Position: 1},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2000, time.April, 1),
			Rating:        4.02,
			URL:           "https://www.goodreads.com/book/show/47212",
			ImageURL:      "https://images.gr-assets.com/books/1419456275m/47212.jpg",
		},
		{
			ID:            6585201,
			WorkID:        6778696,
			Title:         "Changes (The Dresden Files, #12)",
			TitleNoSeries: "Changes",
			Series:        Series{Title: "The Dresden Files", Position: 12},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2010, time.April, 6),
			Rating:        4.54,
			URL:           "https://www.goodreads.com/book/show/6585201",
			ImageURL:      "https://images.gr-assets.com/books/1304027244m/6585201.jpg",
		},
		{
			ID:            91477,
			WorkID:        855288,
			Title:         "Fool Moon (The Dresden Files, #2)",
			TitleNoSeries: "Fool Moon",
			Series:        Series{Title: "The Dresden Files", Position: 2},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2001, time.January, 1),
			Rating:        4.03,
			URL:           "https://www.goodreads.com/book/show/91477",
			ImageURL:      "https://images.gr-assets.com/books/1507307616m/91477.jpg",
		},
		{
			ID:            7779059,
			WorkID:        10351697,
			Title:         "Side Jobs: Stories from the Dresden Files (The Dresden Files, #12.5)",
			TitleNoSeries: "Side Jobs: Stories from the Dresden Files",
			Series:        Series{Title: "The Dresden Files", Position: 12.5},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2010, time.October, 26),
			Rating:        4.25,
			URL:           "https://www.goodreads.com/book/show/7779059",
			ImageURL:      "https://images.gr-assets.com/books/1269115846m/7779059.jpg",
		},
		{
			ID:            91476,
			WorkID:        803205,
			Title:         "Grave Peril (The Dresden Files, #3)",
			TitleNoSeries: "Grave Peril",
			Series:        Series{Title: "The Dresden Files", Position: 3},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2001, time.September, 1),
			Rating:        4.18,
			URL:           "https://www.goodreads.com/book/show/91476",
			ImageURL:      "https://images.gr-assets.com/books/1266470209m/91476.jpg",
		},
		{
			ID:            91478,
			WorkID:        912988,
			Title:         "Summer Knight (The Dresden Files, #4)",
			TitleNoSeries: "Summer Knight",
			Series:        Series{Title: "The Dresden Files", Position: 4},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2002, time.September, 3),
			Rating:        4.3,
			URL:           "https://www.goodreads.com/book/show/91478",
			ImageURL:      "https://images.gr-assets.com/books/1345557469m/91478.jpg",
		},
		{
			ID:            91479,
			WorkID:        2183,
			Title:         "Death Masks (The Dresden Files, #5)",
			TitleNoSeries: "Death Masks",
			Series:        Series{Title: "The Dresden Files", Position: 5},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2003, time.August, 1),
			Rating:        4.32,
			URL:           "https://www.goodreads.com/book/show/91479",
			ImageURL:      "https://images.gr-assets.com/books/1345557713m/91479.jpg",
		},
		{
			ID:            99383,
			WorkID:        227172,
			Title:         "Blood Rites (The Dresden Files, #6)",
			TitleNoSeries: "Blood Rites",
			Series:        Series{Title: "The Dresden Files", Position: 6},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2004, time.August, 1),
			Rating:        4.33,
			URL:           "https://www.goodreads.com/book/show/99383",
			ImageURL:      "https://images.gr-assets.com/books/1345557965m/99383.jpg",
		},
		{
			ID:            17683,
			WorkID:        6614452,
			Title:         "Dead Beat (The Dresden Files, #7)",
			TitleNoSeries: "Dead Beat",
			Series:        Series{Title: "The Dresden Files", Position: 7},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2005, time.May, 3),
			Rating:        4.44,
			URL:           "https://www.goodreads.com/book/show/17683",
			ImageURL:      "https://images.gr-assets.com/books/1345667776m/17683.jpg",
		},
		{
			ID:            91475,
			WorkID:        1254936,
			Title:         "White Night (The Dresden Files, #9)",
			TitleNoSeries: "White Night",
			Series:        Series{Title: "The Dresden Files", Position: 9},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2007, time.April, 3),
			Rating:        4.41,
			URL:           "https://www.goodreads.com/book/show/91475",
			ImageURL:      "https://images.gr-assets.com/books/1309552288m/91475.jpg",
		},
		{
			ID:            91474,
			WorkID:        576222,
			Title:         "Proven Guilty (The Dresden Files, #8)",
			TitleNoSeries: "Proven Guilty",
			Series:        Series{Title: "The Dresden Files", Position: 8},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2006, time.February, 1),
			Rating:        4.42,
			URL:           "https://www.goodreads.com/book/show/91474",
			ImageURL:      "https://images.gr-assets.com/books/1345667469m/91474.jpg",
		},
		{
			ID:            927979,
			WorkID:        2054834,
			Title:         "Small Favor (The Dresden Files, #10)",
			TitleNoSeries: "Small Favor",
			Series:        Series{Title: "The Dresden Files", Position: 10},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2008, time.April, 1),
			Rating:        4.44,
			URL:           "https://www.goodreads.com/book/show/927979",
			ImageURL:      "https://images.gr-assets.com/books/1298085176m/927979.jpg",
		},
		{
			ID:            3475161,
			WorkID:        3516480,
			Title:         "Turn Coat (The Dresden Files, #11)",
			TitleNoSeries: "Turn Coat",
			Series:        Series{Title: "The Dresden Files", Position: 11},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2009, time.April, 7),
			Rating:        4.45,
			URL:           "https://www.goodreads.com/book/show/3475161",
			ImageURL:      "https://images.gr-assets.com/books/1304027128m/3475161.jpg",
		},
		{
			ID:            8058301,
			WorkID:        12731936,
			Title:         "Ghost Story (The Dresden Files,  #13)",
			TitleNoSeries: "Ghost Story",
			Series:        Series{Title: "The Dresden Files", Position: 13},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2011, time.January, 1),
			Rating:        4.25,
			URL:           "https://www.goodreads.com/book/show/8058301",
			ImageURL:      "https://images.gr-assets.com/books/1329104700m/8058301.jpg",
		},
		{
			ID:            12216302,
			WorkID:        17189468,
			Title:         "Cold Days (The Dresden Files, #14)",
			TitleNoSeries: "Cold Days",
			Series:        Series{Title: "The Dresden Files", Position: 14},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2012, time.November, 27),
			Rating:        4.51,
			URL:           "https://www.goodreads.com/book/show/12216302",
			ImageURL:      "https://images.gr-assets.com/books/1345145377m/12216302.jpg",
		},
		{
			ID:            19486421,
			WorkID:        23811929,
			Title:         "Skin Game (The Dresden Files, #15)",
			TitleNoSeries: "Skin Game",
			Series:        Series{Title: "The Dresden Files", Position: 15},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2014, time.May, 27),
			Rating:        4.56,
			URL:           "https://www.goodreads.com/book/show/19486421",
			ImageURL:      "https://images.gr-assets.com/books/1387236318m/19486421.jpg",
		},
		{
			ID:            2575572,
			WorkID:        2588518,
			Title:         "Backup (The Dresden Files, #10.4)",
			TitleNoSeries: "Backup",
			Series:        Series{Title: "The Dresden Files", Position: 10.4},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2008, time.October, 31),
			Rating:        4.12,
			URL:           "https://www.goodreads.com/book/show/2575572",
			ImageURL:      "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		},
		{
			ID:            22249640,
			WorkID:        40515430,
			Title:         "Peace Talks (The Dresden Files, #16)",
			TitleNoSeries: "Peace Talks",
			Series:        Series{Title: "The Dresden Files", Position: 16},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.Date{},
			Rating:        4.44,
			URL:           "https://www.goodreads.com/book/show/22249640",
			ImageURL:      "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		},
		{
			ID:            4271488,
			WorkID:        4319010,
			Title:         "Vignette (The Dresden Files, #5.5)",
			TitleNoSeries: "Vignette",
			Series:        Series{Title: "The Dresden Files", Position: 5.5},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2008, 1, 1),
			Rating:        4.06,
			URL:           "https://www.goodreads.com/book/show/4271488",
			ImageURL:      "https://images.gr-assets.com/books/1476229870m/4271488.jpg",
		},
		{
			ID:            12183815,
			WorkID:        17155691,
			Title:         "Brief Cases (The Dresden Files, #15.1)",
			TitleNoSeries: "Brief Cases",
			Series:        Series{Title: "The Dresden Files", Position: 15.1},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2018, time.June, 5),
			Rating:        4.41,
			URL:           "https://www.goodreads.com/book/show/12183815",
			ImageURL:      "https://images.gr-assets.com/books/1513644037m/12183815.jpg",
		},
	}

	expectedButcher = []Book{
		{
			ID:            47212,
			ISBN:          "0451457811",
			ISBN13:        "9780451457813",
			Title:         "Storm Front (The Dresden Files, #1)",
			TitleNoSeries: "Storm Front",
			Series:        Series{Title: "The Dresden Files", Position: 1},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2000, time.April, 1),
			Rating:        4.02,
			URL:           "https://www.goodreads.com/book/show/47212",
			ImageURL:      "https://images.gr-assets.com/books/1419456275m/47212.jpg",
			Description:   `<b>HARRY DRESDEN — WIZARD</b><br /><br /><i>Lost Items Found. Paranormal Investigations. Consulting. Advice. Reasonable Rates. No Love Potions, Endless Purses, or Other Entertainment.</i><br /><br />Harry Dresden is the best at what he does. Well, technically, he's the <i>only</i> at what he does. So when the Chicago P.D. has a case that transcends mortal creativity or capability, they come to him for answers. For the "everyday" world is actually full of strange and magical things—and most don't play well with humans. That's where Harry comes in. Takes a wizard to catch a—well, whatever. There's just one problem. Business, to put it mildly, stinks.<br /><br />So when the police bring him in to consult on a grisly double murder committed with black magic, Harry's seeing dollar signs. But where there's black magic, there's a black mage behind it. And now that mage knows Harry's name. And that's when things start to get interesting.<br /><br />Magic - it can get a guy killed.`,
		},
		{
			ID:            91477,
			ISBN:          "0451458125",
			ISBN13:        "9780451458124",
			Title:         "Fool Moon (The Dresden Files, #2)",
			TitleNoSeries: "Fool Moon",
			Series:        Series{Title: "The Dresden Files", Position: 2},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2001, time.January, 9),
			Rating:        4.03,
			URL:           "https://www.goodreads.com/book/show/91477",
			ImageURL:      "https://images.gr-assets.com/books/1507307616m/91477.jpg",
			Description:   `<b>Harry Dresden--Wizard</b><br />Lost Items Found. Paranormal Investigations. Consulting. Advice. Reasonable Rates. No Love Potions, Endless Purses, or Other Entertainment.<br /><br />Business has been slow. Okay, business has been dead. And not even of the undead variety. You would think Chicago would have a little more action for the only professional wizard in the phone book. But lately, Harry Dresden hasn't been able to dredge up any kind of work--magical <i>or</i> mundane.<br /><br />But just when it looks like he can't afford his next meal, a murder comes along that requires his particular brand of supernatural expertise.<br /><br />A brutally mutilated corpse. Strange-looking paw prints. A full moon. Take three guesses--and the first two don't count...`,
		},
		{
			ID:            91476,
			ISBN:          "0451458443",
			ISBN13:        "9780451458445",
			Title:         "Grave Peril (The Dresden Files, #3)",
			TitleNoSeries: "Grave Peril",
			Series:        Series{Title: "The Dresden Files", Position: 3},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2001, time.September, 4),
			Rating:        4.18,
			URL:           "https://www.goodreads.com/book/show/91476",
			ImageURL:      "https://images.gr-assets.com/books/1266470209m/91476.jpg",
			Description:   `<i>An alternative cover edition with a different page count exists <a href="https://www.goodreads.com/book/show/13511897.here" title="here" rel="nofollow">here</a>.</i><br /><br />Harry Dresden - Wizard<br />Lost Items Found. Paranormal Investigations. Consulting. Advice. Reasonable Rates. No Love Potions, Endless Purses, or Other Entertainment.<br /><br />Harry Dresden has faced some pretty terrifying foes during his career. Giant scorpions. Oversexed vampires. Psychotic werewolves. It comes with the territory when you're the only professional wizard in the Chicago-area phone book.<br /><br />But in all Harry's years of supernatural sleuthing, he's never faced anything like this: The spirit world has gone postal. All over Chicago, ghosts are causing trouble - and not just of the door-slamming, boo-shouting variety. These ghosts are tormented, violent, and deadly. Someone - or <i>something</i> - is purposely stirring them up to wreak unearthly havoc. But why? And why do so many of the victims have ties to Harry? If Harry doesn't figure it out soon, he could wind up a ghost himself....`,
		},
		{
			ID:            91478,
			ISBN:          "0451458923",
			ISBN13:        "9780451458926",
			Title:         "Summer Knight (The Dresden Files, #4)",
			TitleNoSeries: "Summer Knight",
			Series:        Series{Title: "The Dresden Files", Position: 4},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2002, time.September, 3),
			Rating:        4.3,
			URL:           "https://www.goodreads.com/book/show/91478",
			ImageURL:      "https://images.gr-assets.com/books/1345557469m/91478.jpg",
			Description:   `<i>For the 1st printing edition of this ISBN, see <a href="https://www.goodreads.com/book/show/10446473.here" title="here" rel="nofollow">here</a>.</i><br /><br />HARRY DRESDEN -- WIZARD<br /><br />Lost items found. Paranormal Investigations. Consulting. Advice. Reasonable Rates.<br />No Love Potions, Endless Purses, or Other Entertainment<br /><br />Ever since his girlfriend left town to deal with her newly acquired taste for blood, Harry Dresden has been down and out in Chicago. He can't pay his rent. He's alienating his friends. He can't even recall the last time he took a shower.<br /><br />The only professional wizard in the phone book has become a desperate man.<br /><br />And just when it seems things can't get any worse, in saunters the Winter Queen of Faerie. She has an offer Harry can't refuse if he wants to free himself of the supernatural hold his faerie godmother has over him--and hopefully end his run of bad luck. All he has to do is find out who murdered the Summer Queen's right-hand man, the Summer Knight, and clear the Winter Queen's name.<br /><br />It seems simple enough, but Harry knows better than to get caught in the middle of faerie politics. Until he finds out that the fate of the entire world rests on his solving this case. No pressure or anything...`,
		},
		{
			ID:            91479,
			ISBN:          "0451459407",
			ISBN13:        "9780451459404",
			Title:         "Death Masks (The Dresden Files, #5)",
			TitleNoSeries: "Death Masks",
			Series:        Series{Title: "The Dresden Files", Position: 5},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2003, time.August, 5),
			Rating:        4.32,
			URL:           "https://www.goodreads.com/book/show/91479",
			ImageURL:      "https://images.gr-assets.com/books/1345557713m/91479.jpg",
			Description:   `Harry Dresden, Chicago's only practicing professional wizard, should be happy that business is pretty good for a change. But now he's getting more than he bargained for:<br /><br />A duel with the Red Court of Vampires' champion, who must kill Harry to end the war between vampires and wizards...<br /><br />Professional hit men using Harry for target practice...<br /><br />The missing Shroud of Turin...<br /><br />A handless and headless corpse the Chicago police need identified...<br /><br />Not to mention the return of Harry's ex-girlfriend Susan, who's still struggling with her semi-vampiric nature. And who seems to have a new man in her life.<br /><br />Some days, it just doesn't pay to get out of bed. No matter how much you're charging.`,
		},
		{
			ID:            99383,
			ISBN:          "0451459873",
			ISBN13:        "9780451459879",
			Title:         "Blood Rites (The Dresden Files, #6)",
			TitleNoSeries: "Blood Rites",
			Series:        Series{Title: "The Dresden Files", Position: 6},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2004, time.August, 3),
			Rating:        4.34,
			URL:           "https://www.goodreads.com/book/show/99383",
			ImageURL:      "https://images.gr-assets.com/books/1345557965m/99383.jpg",
			Description:   `For Harry Dresden, Chicago's only professional wizard, there have been worse assignments than going undercover on the set of an adult film. Dodging flaming monkey poo, for instance. Or going toe-to-leaf with a walking plant monster. Still, there is something more troubling than usual about his newest case. The film's producer believes he's the target of a sinister entropy curse, but it's the women around him who are dying, in increasingly spectacular ways.<br /><br />Harry is doubly frustrated because he got involved with this bizarre mystery only as a favor to Thomas, his flirtatious, self-absorbed vampire acquaintance of dubious integrity. Thomas has a personal stake in the case Harry can't quite figure out, until his investigation leads him straight to Thomas' oversexed vampire family. Harry is about to discover that Thomas' family tree has been hiding a shocking secret; a revelation that will change Harry's life forever.`,
		},
		{
			ID:            17683,
			ISBN:          "045146091X",
			ISBN13:        "9780451460912",
			Title:         "Dead Beat (The Dresden Files, #7)",
			TitleNoSeries: "Dead Beat",
			Series:        Series{Title: "The Dresden Files", Position: 7},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2006, time.January, 1),
			Rating:        4.44,
			URL:           "https://www.goodreads.com/book/show/17683",
			ImageURL:      "https://images.gr-assets.com/books/1345667776m/17683.jpg",
			Description:   `Paranormal investigations are Harry Dresden’s business and Chicago is his beat, as he tries to bring law and order to a world of wizards and monsters that exists alongside everyday life. And though most inhabitants of the Windy City don’t believe in magic, the Special Investigations Department of the Chicago PD knows better. <br /><br />Karrin Murphy is the head of S. I. and Harry’s good friend. So when a killer vampire threatens to destroy Murphy’s reputation unless Harry does her bidding, he has no choice. The vampire wants the Word of Kemmler (whatever that is) and all the power that comes with it. Now, Harry is in a race against time—and six merciless necromancers—to find the Word before Chicago experiences a Halloween night to wake the dead...`,
		},
		{
			ID:            91475,
			ISBN:          "0451461401",
			ISBN13:        "9780451461407",
			Title:         "White Night (The Dresden Files, #9)",
			TitleNoSeries: "White Night",
			Series:        Series{Title: "The Dresden Files", Position: 9},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2007, time.April, 3),
			Rating:        4.41,
			URL:           "https://www.goodreads.com/book/show/91475",
			ImageURL:      "https://images.gr-assets.com/books/1309552288m/91475.jpg",
			Description:   `<b>The inspiration for the Sci Fi channel television series</b> <br /><br /> In Chicago, someone has been killing practitioners of magic, those incapable of becoming full-fledged wizards. Shockingly, all the evidence points to Harry Dresden's half-brother, Thomas, as the murderer. Determined to clear his sibling's name, Harry uncovers a conspiracy within the White Council of Wizards that threatens not only him, but his nearest and dearest, too...`,
		},
		{
			ID:            91474,
			ISBN:          "0451461037",
			ISBN13:        "9780451461032",
			Title:         "Proven Guilty (The Dresden Files, #8)",
			TitleNoSeries: "Proven Guilty",
			Series:        Series{Title: "The Dresden Files", Position: 8},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2007, time.February, 6),
			Rating:        4.42,
			URL:           "https://www.goodreads.com/book/show/91474",
			ImageURL:      "https://images.gr-assets.com/books/1345667469m/91474.jpg",
			Description:   `There's no love lost between Harry Dresden, the only wizard in the Chicago phone book, and the White Council of Wizards, who find him brash and undisciplined. But war with the vampires has thinned their ranks, so the Council has drafted Harry as a Warden and assigned him to look into rumors of black magic in the Windy City.<br /><br />As Harry adjusts to his new role, another problem arrives in the form of the tattooed and pierced daughter of an old friend, all grown-up and already in trouble. Her boyfriend is the only suspect in what looks like a supernatural assault straight out of a horror film. Malevolent entities that feed on fear are loose in Chicago, but it's all in a day's work for a wizard, his faithful dog, and a talking skull named Bob....`,
		},
		{
			ID:            6585201,
			ISBN:          "045146317X",
			ISBN13:        "9780451463173",
			Title:         "Changes (The Dresden Files, #12)",
			TitleNoSeries: "Changes",
			Series:        Series{Title: "The Dresden Files", Position: 12},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2010, time.April, 6),
			Rating:        4.52,
			URL:           "https://www.goodreads.com/book/show/6585201",
			ImageURL:      "https://images.gr-assets.com/books/1304027244m/6585201.jpg",
			Description:   `Long ago, Susan Rodriguez was Harry Dresden's lover-until she was attacked by his enemies, leaving her torn between her own humanity and the bloodlust of the vampiric Red Court. Susan then disappeared to South America, where she could fight both her savage gift and those who cursed her with it. <br /><br />Now Arianna Ortega, Duchess of the Red Court, has discovered a secret Susan has long kept, and she plans to use it-against Harry. To prevail this time, he may have no choice but to embrace the raging fury of his own untapped dark power. Because Harry's not fighting to save the world... <br /><br />He's fighting to save his child.`,
		},
		{
			ID:            927979,
			ISBN:          "0451461894",
			ISBN13:        "9780451461896",
			Title:         "Small Favor (The Dresden Files, #10)",
			TitleNoSeries: "Small Favor",
			Series:        Series{Title: "The Dresden Files", Position: 10},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2008, time.April, 1),
			Rating:        4.44,
			URL:           "https://www.goodreads.com/book/show/927979",
			ImageURL:      "https://images.gr-assets.com/books/1298085176m/927979.jpg",
			Description:   `<b>THE <i>New York Times</i> Bestseller</b><br /><br />Harry Dresden's life finally seems to be calming down -- until a shadow from the past returns. Mab, monarch of the Sidhe Winter Court, calls in an old favor from Harry -- one small favor that will trap him between a nightmarish foe and an equally deadly ally, and that will strain his skills -- and loyalties -- to their very limits.`,
		},
		{
			ID:            29396,
			ISBN:          "044101268X",
			ISBN13:        "9780441012688",
			Title:         "Furies of Calderon (Codex Alera, #1)",
			TitleNoSeries: "Furies of Calderon",
			Series:        Series{Title: "Codex Alera", Position: 1},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2005, time.June, 28),
			Rating:        4.13,
			URL:           "https://www.goodreads.com/book/show/29396",
			ImageURL:      "https://images.gr-assets.com/books/1329104514m/29396.jpg",
			Description:   `For a thousand years, the people of Alera have united against the aggressive and threatening races that inhabit the world, using their unique bond with the furies - elementals of earth, air, fire, water, and metal. <br /><br />But now, Gaius Sextus, First Lord of Alera, grows old and lacks an heir. Ambitious High Lords plot and maneuver to place their Houses in positions of power, and a war of succession looms on the horizon. <br /><br />Far from city politics in the Calderon Valley, the boy Tavi struggles with his lack of furycrafting. At fifteen, he has no wind fury to help him fly, no fire fury to light his lamps. Yet as the Alerans' most savage enemy - the Marat - return to the Valley, he will discover that his destiny is much greater than he could ever imagine. <br /><br />Caught in a storm of deadly wind furies, Tavi saves the life of a runaway slave named Amara. But she is actually a spy for Gaius Sextus, sent to the Valley to gather intelligence on traitors to the Crown, who may be in league with the barbaric Marat horde. And when the Valley erupts in chaos - when rebels war with loyalists and furies clash with furies - Amara will find Tavi's courage and resourcefulness to be a power greater than any fury - one that could turn the tides of war.`,
		},
		{
			ID:            3475161,
			ISBN:          "0451462564",
			ISBN13:        "9780451462565",
			Title:         "Turn Coat (The Dresden Files, #11)",
			TitleNoSeries: "Turn Coat",
			Series:        Series{Title: "The Dresden Files", Position: 11},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2009, time.April, 7),
			Rating:        4.45,
			URL:           "https://www.goodreads.com/book/show/3475161",
			ImageURL:      "https://images.gr-assets.com/books/1304027128m/3475161.jpg",
			Description:   `Accused of treason against the Wizards of the White Council, Warden Morgan goes in search of Harry Dresden in a desperate attempt to clear his name and stop the deadly punishment from taking place in this latest thrilling addition to the Dresden Files series.`,
		},
		{
			ID:            12216302,
			ISBN:          "0451464400",
			ISBN13:        "9780451464408",
			Title:         "Cold Days (The Dresden Files, #14)",
			TitleNoSeries: "Cold Days",
			Series:        Series{Title: "The Dresden Files", Position: 14},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2012, time.November, 27),
			Rating:        4.51,
			URL:           "https://www.goodreads.com/book/show/12216302",
			ImageURL:      "https://images.gr-assets.com/books/1345145377m/12216302.jpg",
			Description:   `You can't keep a good wizard down - even when he wants to stay that way.<br /><br />For years, Harry Dresden has been Chicago's only professional wizard, but a bargain made in desperation with the Queen of Air and Darkness has forced him into a new job: professional killer.<br /><br />Mab, the mother of wicked faeries, has restored the mostly-dead wizard to health, and dispatches him upon his first mission - to bring death to an immortal. Even as he grapples with the impossible task, Dresden learns of a looming danger to Demonreach, the living island hidden upon Lake Michigan, a place whose true purpose and dark potential have the potential to destroy billions and to land Dresden in the deepest trouble he has ever known - even deeper than being dead. How messed up is that?<br /><br />Beset by his new enemies and hounded by the old, Dresden has only twenty four hours to reconnect with his old allies, prevent a cataclysm and do the impossible - all while the power he bargained to get - but never meant to keep - lays siege to his very soul.<br /><br />Magic. It can get a guy killed.`,
		},
		{
			ID:            8058301,
			ISBN:          "045146379X",
			ISBN13:        "9780451463791",
			Title:         "Ghost Story (The Dresden Files,  #13)",
			TitleNoSeries: "Ghost Story",
			Series:        Series{Title: "The Dresden Files", Position: 13},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2011, time.July, 26),
			Rating:        4.25,
			URL:           "https://www.goodreads.com/book/show/8058301",
			ImageURL:      "https://images.gr-assets.com/books/1329104700m/8058301.jpg",
			Description:   `When we last left the mighty wizard detective Harry Dresden, he wasn't doing well. In fact, he had been murdered by an unknown assassin.<br /><br />But being dead doesn't stop him when his friends are in danger. Except now he has no body, and no magic to help him. And there are also several dark spirits roaming the Chicago shadows who owe Harry some payback of their own.<br /><br />To save his friends—and his own soul—Harry will have to pull off the ultimate trick without any magic...`,
		},
		{
			ID:            19486421,
			ISBN:          "0451464397",
			ISBN13:        "9780451464392",
			Title:         "Skin Game (The Dresden Files, #15)",
			TitleNoSeries: "Skin Game",
			Series:        Series{Title: "The Dresden Files", Position: 15},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2014, time.May, 27),
			Rating:        4.56,
			URL:           "https://www.goodreads.com/book/show/19486421",
			ImageURL:      "https://images.gr-assets.com/books/1387236318m/19486421.jpg",
			Description:   `Harry Dresden, Chicago's only professional wizard, is about to have a very bad day…<br /><br />Because as Winter Knight to the Queen of Air and Darkness, Harry never knows what the scheming Mab might want him to do. Usually, it’s something awful.<br /><br />He doesn’t know the half of it…<br /><br />Mab has just traded Harry’s skills to pay off one of her debts. And now he must help a group of supernatural villains—led by one of Harry’s most dreaded and despised enemies, Nicodemus Archleone—to break into the highest-security vault in town, so that they can then access the highest-security vault in the Nevernever. <br /><br />It's a smash and grab job to recover the literal Holy Grail from the vaults of the greatest treasure hoard in the supernatural world—which belongs to the one and only Hades, Lord of the freaking Underworld and generally unpleasant character. Worse, Dresden suspects that there is another game afoot that no one is talking about. And he's dead certain that Nicodemus has no intention of allowing any of his crew to survive the experience. Especially Harry.<br /><br />Dresden's always been tricky, but he's going to have to up his backstabbing game to survive this mess—assuming his own allies don’t end up killing him before his enemies get the chance…`,
		},
		{
			ID:            133664,
			ISBN:          "0441013406",
			ISBN13:        "9780441013401",
			Title:         "Academ's Fury (Codex Alera, #2)",
			TitleNoSeries: "Academ's Fury",
			Series:        Series{Title: "Codex Alera", Position: 2},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2010, time.February, 1),
			Rating:        4.28,
			URL:           "https://www.goodreads.com/book/show/133664",
			ImageURL:      "https://images.gr-assets.com/books/1381026900m/133664.jpg",
			Description:   `For centuries, the people of Alera have relied on the power of the furies to protect them from outside invaders. But the gravest threat might be closer than they think.<br /><br />Tavi has escaped the Calderon Valley and the mysterious attack of the Marat on his homeland. But he is far from safe, as trying to keep up the illusion of being a student while secretly training as one of the First Lord's spies is a dangerous game. And he has not yet learned to use the furies, making him especially vulnerable.<br /><br />When the attack comes it's on two fronts. A sudden strike threatens the First Lord's life and threatens to plunge the land into civil war. While in the Calderon Valley, the threat faced from the Marat is dwarfed by an ancient menace. And Tavi must learn to harness the furies if he has any chance of fighting the greatest threat Alera has ever known . . .`,
		},
		{
			ID:            346087,
			ISBN:          "0441015271",
			ISBN13:        "9780441015276",
			Title:         "Captain's Fury (Codex Alera, #4)",
			TitleNoSeries: "Captain's Fury",
			Series:        Series{Title: "Codex Alera", Position: 4},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2007, time.December, 4),
			Rating:        4.38,
			URL:           "https://www.goodreads.com/book/show/346087",
			ImageURL:      "https://images.gr-assets.com/books/1315083292m/346087.jpg",
			Description:   `After two years of bitter conflict with the hordes of invading Canim, Tavi of Calderon, now Captain of the First Aleran Legion, realizes that a peril far greater than the Canim exists--the mysterious threat that drove the savage Canim to flee their homeland. <br /><br />Now, Tavi must find a way to overcome the centuries-old animosities between Aleran and Cane if an alliance is to be forged against their mutual enemy. And he must lead his legion in defiance of the law, against friend and foe--or no one will have a chance of survival . . .`,
		},
		{
			ID:            29394,
			ISBN:          "0441014348",
			ISBN13:        "9780441014347",
			Title:         "Cursor's Fury (Codex Alera, #3)",
			TitleNoSeries: "Cursor's Fury",
			Series:        Series{Title: "Codex Alera", Position: 3},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2006, time.December, 5),
			Rating:        4.37,
			URL:           "https://www.goodreads.com/book/show/29394",
			ImageURL:      "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
			Description:   `The power-hungry High Lord of Kalare has launched a rebellion against the aging First Lord, Gaius Sextus, who with the loyal forces of Alera must fight beside the unlikeliest of allies-the equally contentious High Lord of Aquitaine. <br /><br />Meanwhile, young Tavi of Calderon joins a newly formed legion under an assumed name even as the ruthless Kalare unites with the Canim, bestial enemies of the realm whose vast numbers spell certain doom for Alera. <br /><br />When treachery from within destroys the army's command structure, Tavi finds himself leading an inexperienced, poorly equipped legion-the only force standing between the Canim horde and the war-torn realm.`,
		},
		{
			ID:            2903736,
			ISBN:          "0441016383",
			ISBN13:        "9780441016389",
			Title:         "Princeps' Fury (Codex Alera, #5)",
			TitleNoSeries: "Princeps' Fury",
			Series:        Series{Title: "Codex Alera", Position: 5},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2008, time.November, 25),
			Rating:        4.36,
			URL:           "https://www.goodreads.com/book/show/2903736",
			ImageURL:      "https://images.gr-assets.com/books/1315082776m/2903736.jpg",
			Description:   `Tavi of Calderon, now recognized as Princeps Gaius Octavian and heir to the crown, has achieved a fragile alliance with Alera’s oldest foes, the savage Canim. But when Tavi and his legions guide the Canim safely to their lands, his worst fears are realized.<br /><br />The dreaded Vord—the enemy of Aleran and Cane alike—have spent the last three years laying waste to the Canim homeland. And when the Alerans are cut off from their ships, they find themselves with no choice but to fight shoulder to shoulder if they are to survive.<br /><br />For a thousand years, Alera and her furies have withstood every enemy, and survived every foe.<br /><br />The thousand years are over…`,
		},
		{
			ID:            6316821,
			ISBN:          "044101769X",
			ISBN13:        "9780441017690",
			Title:         "First Lord's Fury (Codex Alera, #6)",
			TitleNoSeries: "First Lord's Fury",
			Series:        Series{Title: "Codex Alera", Position: 6},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2009, time.November, 24),
			Rating:        4.38,
			URL:           "https://www.goodreads.com/book/show/6316821",
			ImageURL:      "https://images.gr-assets.com/books/1327903582m/6316821.jpg",
			Description:   `For Gaius Octavian, life has been one long battle. Now, the end of all he fought for is close at hand. The brutal, dreaded Vord are on the march against Alera. And perhaps for the final time, Gaius Octavian and his legions must stand against the enemies of his people. And it will take all his intelligence, ingenuity, and furycraft to save their world from eternal darkness.`,
		},
		{
			ID:            7779059,
			ISBN:          "045146365X",
			ISBN13:        "9780451463654",
			Title:         "Side Jobs: Stories from the Dresden Files (The Dresden Files, #12.5)",
			TitleNoSeries: "Side Jobs: Stories from the Dresden Files",
			Series:        Series{Title: "The Dresden Files", Position: 12.5},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2010, time.October, 26),
			Rating:        4.25,
			URL:           "https://www.goodreads.com/book/show/7779059",
			ImageURL:      "https://images.gr-assets.com/books/1269115846m/7779059.jpg",
			Description: `Here, together for the first time, are the shorter from Jim Butcher's DRESDEN FILES series — a compendium of cases that Harry and his cadre of allies managed to close in record time. The tales range from the deadly serious to the absurdly hilarious. Also included is a new, never-before-published novella that takes place after the cliff-hanger ending of the new April 2010 hardcover, <i>Changes</i>.<br /><br />Contains:<br />+ "Restoration of Faith"<br />+ "Vignette"<br />+ "Something Borrowed" -- from <em>
  <a href="https://www.goodreads.com/book/show/84156.My_Big_Fat_Supernatural_Wedding" title="My Big Fat Supernatural Wedding" rel="nofollow">My Big Fat Supernatural Wedding</a>
</em><br />+ "It's My Birthday Too" -- from <em>
  <a href="https://www.goodreads.com/book/show/140098.Many_Bloody_Returns" title="Many Bloody Returns" rel="nofollow">Many Bloody Returns</a>
</em><br />+ "Heorot" -- from <em>
  <a href="https://www.goodreads.com/book/show/1773616.My_Big_Fat_Supernatural_Honeymoon" title="My Big Fat Supernatural Honeymoon" rel="nofollow">My Big Fat Supernatural Honeymoon</a>
</em><br />+ "Day Off" -- from <em>
  <a href="https://www.goodreads.com/book/show/2871256.Blood_Lite" title="Blood Lite" rel="nofollow">Blood Lite</a>
</em><br />+ "<a href="https://www.goodreads.com/book/show/2575572.Backup" title="Backup" rel="nofollow">Backup</a>" -- novelette from Thomas' point of view, originally published by Subterranean Press<br />+ "The Warrior" -- novelette from <em>
  <a href="https://www.goodreads.com/book/show/3475145.Mean_Streets" title="Mean Streets" rel="nofollow">Mean Streets</a>
</em><br />+ "Last Call" -- from <em>
  <a href="https://www.goodreads.com/book/show/6122181.Strange_Brew" title="Strange Brew" rel="nofollow">Strange Brew</a>
</em><br />+ "Love Hurts" -- from <em>
  <a href="https://www.goodreads.com/book/show/7841656.Songs_of_Love_and_Death" title="Songs of Love and Death" rel="nofollow">Songs of Love and Death</a>
</em><br />+ <em>Aftermath</em> -- all-new novella from Murphy's point of view, set forty-five minutes after the end of <em>
  <a href="https://www.goodreads.com/book/show/6585201.Changes" title="Changes" rel="nofollow">Changes</a>
</em>`,
		},
		{
			ID:            24876258,
			ISBN:          "0451466802",
			ISBN13:        "9780451466808",
			Title:         "The Aeronaut's Windlass (The Cinder Spires, #1)",
			TitleNoSeries: "The Aeronaut's Windlass",
			Series:        Series{Title: "The Cinder Spires", Position: 1},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2015, time.September, 29),
			Rating:        4.18,
			URL:           "https://www.goodreads.com/book/show/24876258",
			ImageURL:      "https://images.gr-assets.com/books/1425415066m/24876258.jpg",
			Description:   `Jim Butcher, the #1 New York Times bestselling author of The Dresden Files and the Codex Alera novels, conjures up a new series set in a fantastic world of noble families, steam-powered technology, and magic-wielding warriors…<br /><br />Since time immemorial, the Spires have sheltered humanity, towering for miles over the mist-shrouded surface of the world. Within their halls, aristocratic houses have ruled for generations, developing scientific marvels, fostering trade alliances, and building fleets of airships to keep the peace.<br /><br />Captain Grimm commands the merchant ship, <i>Predator</i>. Fiercely loyal to Spire Albion, he has taken their side in the cold war with Spire Aurora, disrupting the enemy’s shipping lines by attacking their cargo vessels. But when the <i>Predator</i> is severely damaged in combat, leaving captain and crew grounded, Grimm is offered a proposition from the Spirearch of Albion—to join a team of agents on a vital mission in exchange for fully restoring <i>Predator</i> to its fighting glory.<br /><br />And even as Grimm undertakes this dangerous task, he will learn that the conflict between the Spires is merely a premonition of things to come. Humanity’s ancient enemy, silent for more than ten thousand years, has begun to stir once more. And death will follow in its wake…`,
		},
		{
			ID:            2637138,
			ISBN:          "0345507460",
			ISBN13:        "9780345507464",
			Title:         "Welcome to the Jungle",
			TitleNoSeries: "Welcome to the Jungle",
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2008, time.October, 21),
			Rating:        4.10,
			URL:           "https://www.goodreads.com/book/show/2637138",
			ImageURL:      "https://images.gr-assets.com/books/1320418408m/2637138.jpg",
			Description:   `When the supernatural world spins out of control, when the police can’t handle what goes bump in the night, when monsters come screaming out of nightmares and into the mean streets, there’s just one man to call: Harry Dresden, the only professional wizard in the Chicago phone book. A police consultant and private investigator, Dresden has to walk the dangerous line between the world of night and the light of day.<br /><br />Now Harry Dresden is investigating a brutal mauling at the Lincoln Park Zoo that has left a security guard dead and many questions unanswered. As an investigator of the supernatural, he senses that there’s more to this case than a simple animal attack, and as Dresden searches for clues to figure out who is really behind the crime, he finds himself next on the victim list, and being hunted by creatures that won’t leave much more than a stain if they catch him.<br /><br />Written exclusively for comics by Jim Butcher, The Dresden Files: Welcome to the Jungle is a brand-new story that’s sure to enchant readers with a blend of gripping mystery and fantastic adventure.`,
		},
		{
			ID:            4961959,
			ISBN:          "0345506391",
			ISBN13:        "9780345506399",
			Title:         "Jim Butcher's The Dresden Files: Storm Front, Volume 1: The Gathering Storm",
			TitleNoSeries: "Jim Butcher's The Dresden Files: Storm Front, Volume 1: The Gathering Storm",
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2009, time.June, 2),
			Rating:        4.36,
			URL:           "https://www.goodreads.com/book/show/4961959",
			ImageURL:      "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
			Description:   `A graphic novel based on the bestselling Harry Dresden books by Jim Butcher!<br /><br />If circumstances surrounding a crime defy the ordinary and evidence points to a suspect who is anything but human, the men and women of the Chicago Police Department call in the one guy who can handle bizarre and often brutal phenomena. Harry Dresden is a wizard who knows firsthand that the everyday world is actually full of strange and magical things—most of which don't play well with humans.<br /><br />Now the cops have turned to Dresden to investigate a horrifying double murder that was committed with black magic. Never one to turn down a paycheck, Dresden also takes on another case—to find a missing husband who has quite likely been dabbling in sorcery. As Dresden tries to solve the seemingly unrelated cases, he is confronted with all the Windy City can blow at him, from the mob to mages and all creatures in between.`,
		},
		{
			ID:            12183815,
			ISBN:          "0451492129",
			ISBN13:        "9780451492128",
			Title:         "Brief Cases (The Dresden Files, #15.1)",
			TitleNoSeries: "Brief Cases",
			Series:        Series{Title: "The Dresden Files", Position: 15.1},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2018, time.June, 5),
			Rating:        4.42,
			URL:           "https://www.goodreads.com/book/show/12183815",
			ImageURL:      "https://images.gr-assets.com/books/1513644037m/12183815.jpg",
			Description:   `<i>Brief Cases</i> is the sequel anthology of <i>Side Jobs</i>, and will be released before <i>Peace Talks</i><br /><br />Set to include the following stories:<br /><br />An exclusive novellette from the perspective of Maggie and Mouse.<br /><br />“Curses” — from The Naked City, edited by Ellen Datlow<br />Takes place between Small Favor and Turn Coat.<br /><br />“AAAA Wizardry” — from the Dresden Files RPG, published by Evil Hat<br />Harry teaches a group of young Wardens his procedure for dealing with supernatural nasties.<br /><br />“Even Hand” — originally from Dark and Stormy Knights, edited by Pat Elrod. Reprinted in Beyond the Pale, edited by Henry Herz.<br />Gentleman Johnnie Marcone clashes with a rival supernatural power. Told from Marcone’s point of view.<br />Takes place between Turn Coat and Changes.<br /><br />“B is for Bigfoot” — from Under My Hat: Tales From the Cauldron, edited by Jonathan Strahan. Republished in Working for Bigfoot.<br />Takes place between Fool Moon and Grave Peril.<br /><br />“I Was A Teenage Bigfoot” — from Blood Lite 3: Aftertaste, edited by Kevin J. Anderson. Republished in Working for Bigfoot.<br />Takes place circa Dead Beat.<br /><br />“Bigfoot on Campus” — from Hex Appeal, edited by P.N. Elrod. Republished in Working for Bigfoot.<br />Takes place between Turn Coat and Changes.<br /><br />“Bombshells” — Molly-POV novella from Dangerous Women, edited by George R. R. Martin and Gardner Duzois.<br />Molly teams up with Justine and Andi to thwart a Fomor plot.<br />Takes place between Ghost Story and Cold Days.<br /><br />“Jury Duty” — short story for Unbound, edited by Shawn Speakman.<br />Harry endures Jury Duty.<br />Set after Skin Game.<br /><br />“Cold Case” — short story from Shadowed Souls, edited by Jim Butcher and Kerrie Hughes.<br />In Molly’s first job in her new role, she teams up with Ramirez to take on a Lovecraft-esque cult.<br />Takes place shortly after Cold Days.<br /><br />“Day One” — short story for Unfettered II, edited by Shawn Speakman.<br />Butters’ first mission.<br />Set after Skin Game.<br /><br />“A Fistful of Warlocks” — short story for Straight Outta Tombstone, edited by David Boop.<br />Luccio takes on necromancers in the Wild West.<br />Set long before the events of the series.`,
		},
		{
			ID:            3475145,
			ISBN:          "0451462491",
			ISBN13:        "9780451462497",
			Title:         "Mean Streets",
			TitleNoSeries: "Mean Streets",
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2009, time.January, 1),
			Rating:        4.01,
			URL:           "https://www.goodreads.com/book/show/3475145",
			ImageURL:      "https://images.gr-assets.com/books/1303858861m/3475145.jpg",
			Description:   `Jim Butcher delivers a hard-boiled tale in which Harry Dresden’s latest case may be his last.<br /><br />Nightside dweller John Taylor is hired by a woman to find something she lost—her memory—in a noir tale from Simon R. Green.<br /><br />Kat Richardson’s Greywalker finds herself in too deep when a “simple job” goes bad and Harper Blaine is enmeshed in a tangle of dark secrets and revenge from beyond the grave.<br /><br />For centuries, the being that we know as Noah lived among us. Now he is dead, and fallen-angel-turned-detective Remy Chandler has been hired to find out who killed him in a whodunit by Thomas E. Sniegoski.`,
		},
		{
			ID:            2575572,
			ISBN:          "1596061820",
			ISBN13:        "9781596061828",
			Title:         "Backup (The Dresden Files, #10.4)",
			TitleNoSeries: "Backup",
			Series:        Series{Title: "The Dresden Files", Position: 10.4},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2008, time.October, 1),
			Rating:        4.12,
			URL:           "https://www.goodreads.com/book/show/2575572",
			ImageURL:      "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
			Description:   `Let's get something clear right up front. I'm not Harry Dresden. Harry's a wizard. A genuine, honest-to-goodness wizard. He's Gandalf on crack and an IV of Red Bull, with a big leather coat and a .44 revolver in his pocket. He'll spit in the eye of gods and demons alike if he thinks it needs to be done, and to hell with the consequences--and yet somehow my little brother manages to remain a decent human being. I'll be damned if I know how. But then, I'll be damned regardless. My name is Thomas Raith, and I'm a monster. So begins "Backup" a twelve thousand word novelette set in Jim Butcher's ultra-popular Dresden Files series. This time Harry's in trouble he knows nothing about, and it's up to his big brother Thomas to track him down and solve those little life-threatening difficulties without his little brother even noticing.<br /><br /><i>Also in <a href="https://www.goodreads.com/book/show/7779059._Side_Jobs" title=" Side Jobs" rel="nofollow"> Side Jobs</a>.</i>`,
		},
		{
			ID:            25807691,
			Title:         "Working for Bigfoot (The Dresden Files, #15.5)",
			TitleNoSeries: "Working for Bigfoot",
			Series:        Series{Title: "The Dresden Files", Position: 15.5},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.Date{},
			Rating:        4.26,
			URL:           "https://www.goodreads.com/book/show/25807691",
			ImageURL:      "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
			Description:   `"B is for Bigfoot" takes place between <i>Fool Moon</i> and <i>Grave Peril</i>. <br />"I Was a Teenage Bigfoot" takes place circa <i>Deadbeat</i>. <br />"Bigfoot on Campus" takes place between <i>Turn Coat</i> and <i>Changes</i>.<br /><br />Chicago wizard-for-hire Harry Dresden is used to mysterious clients with long hair and legs up to here. But when it turns out the long hair covers every square inch of his latest client’s body, and the legs contribute to a nine-foot height, even the redoubtable detective realizes he’s treading new ground. Strength of a River in His Shoulders is one of the legendary forest people, a Bigfoot, and he has a problem that only Harry can solve. His son Irwin is a scion, the child of a supernatural creature and a human. He’s a good kid, but the extraordinary strength of his magical aura has a way of attracting trouble.<br /><br />In the three novellas that make up Working For Bigfoot, collected together for the first time here, readers encounter Dresden at different points in his storied career, and in Irwin’s life. As a middle-schooler, in “B Is for Bigfoot,” Irwin attracts the unwelcome attention of a pair of bullying brothers who are more than they seem, and when Harry steps in, it turns out they have a mystical guardian of their own. At a fancy private high school in “I Was a Teenage Bigfoot,” Harry is called in when Irwin grows ill for the first time, and it’s not just a case of mono. Finally, Irwin is all grown up and has a grown-up’s typical problems as a freshman in college in “Bigfoot on Campus,” or would have if typical included vampires.`,
		},
		{
			ID:            15732549,
			Title:         "Restoration of Faith (The Dresden Files, #0.2)",
			TitleNoSeries: "Restoration of Faith",
			Series:        Series{Title: "The Dresden Files", Position: 0.2},
			Author:        Author{ID: 10746, Name: "Jim Butcher", URL: "https://www.goodreads.com/author/show/10746"},
			PubDate:       date.New(2004, time.January, 1),
			Rating:        4.00,
			URL:           "https://www.goodreads.com/book/show/15732549",
			ImageURL:      "https://images.gr-assets.com/books/1483723101m/15732549.jpg",
			Description:   `A short story of the Dresden Files.<br /><br />(Also included in <a href="https://www.goodreads.com/book/show/7779059.Side_Jobs" title="Side Jobs" rel="nofollow">Side Jobs</a>)`,
		},
	}
)
