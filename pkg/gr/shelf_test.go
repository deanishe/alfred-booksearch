// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package gr

import (
	"testing"
	"time"

	"github.com/fxtlabs/date"
	"github.com/stretchr/testify/assert"
)

// Books from shelf
func TestParseShelf(t *testing.T) {
	t.Parallel()

	books, meta, err := unmarshalShelf(readFile("currently-reading.xml", t))
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	assert.Equal(t, expectedCurrentlyReading, books, "unexpected Books")
	assert.Equal(t, PageData{Start: 1, End: 5, Total: 5}, meta, "unexpected meta")
}

var (
	expectedCurrentlyReading = []Book{
		{
			ID:            31379281,
			Title:         "Deep Cover Jack (Hunt For Reacher #4)",
			TitleNoSeries: "Deep Cover Jack",
			Series:        Series{Title: "Hunt For Reacher", Position: 4},
			Author:        Author{Name: "Diane Capri", ID: 5070259, URL: "https://www.goodreads.com/author/show/5070259"},
			Rating:        4.08,
			URL:           "https://www.goodreads.com/book/show/31379281",
			ImageURL:      "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
			Description:   "In the thrilling follow-up to the ITW Thriller Award Finalist (“Jack and Joe”), FBI Special Agents Kim Otto and Carlos Gaspar will wait no longer. They head to Houston to find Susan Duffy, one of Jack Reacher’s known associates, determined to get answers. But Duffy’s left town, headed for trouble. Otto and Gaspar are right behind her, and powerful enemies with their backs against the wall will have everything to lose.",
		},
		{
			ID:            10383597,
			ISBN:          "0771041411",
			ISBN13:        "9780771041419",
			Title:         "Arguably: Selected Essays",
			TitleNoSeries: "Arguably: Selected Essays",
			PubDate:       date.New(2011, time.January, 1),
			Author:        Author{Name: "Christopher Hitchens", ID: 3956, URL: "https://www.goodreads.com/author/show/3956"},
			Rating:        4.2,
			URL:           "https://www.goodreads.com/book/show/10383597",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1426386037l/10383597._SX98_.jpg",
			Description:   `The first new book of essays by Christopher Hitchens since 2004, <i>Arguably</i> offers an indispensable key to understanding the passionate and skeptical spirit of one of our most dazzling writers, widely admired for the clarity of his style, a result of his disciplined and candid thinking. <br /><br />Topics range from ruminations on why Charles Dickens was among the best of writers and the worst of men to the haunting science fiction of J.G. Ballard; from the enduring legacies of Thomas Jefferson and George Orwell to the persistent agonies of anti-Semitism and jihad. Hitchens even looks at the recent financial crisis and argues for the enduring relevance of Karl Marx. <br /><br />The book forms a bridge between the two parallel enterprises of culture and politics. It reveals how politics justifies itself by culture, and how the latter prompts the former. In this fashion, <i>Arguably</i> burnishes Christopher Hitchens' credentials as (to quote Christopher Buckley) our "greatest living essayist in the English language."`,
		},
		{
			ID:            61886,
			ISBN:          "0007133618",
			ISBN13:        "9780007133611",
			Title:         "The Curse of Chalion (World of the Five Gods, #1)",
			TitleNoSeries: "The Curse of Chalion",
			PubDate:       date.New(2003, time.February, 3),
			Series:        Series{Title: "World of the Five Gods", Position: 1},
			Author:        Author{Name: "Lois McMaster Bujold", ID: 16094, URL: "https://www.goodreads.com/author/show/16094"},
			Rating:        4.15,
			URL:           "https://www.goodreads.com/book/show/61886",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1322571773l/61886._SX98_.jpg",
			Description:   `A man broken in body and spirit, Cazaril, has returned to the noble household he once served as page, and is named, to his great surprise, as the secretary-tutor to the beautiful, strong-willed sister of the impetuous boy who is next in line to rule. <br /><br />It is an assignment Cazaril dreads, for it will ultimately lead him to the place he fears most, the royal court of Cardegoss, where the powerful enemies, who once placed him in chains, now occupy lofty positions. In addition to the traitorous intrigues of villains, Cazaril and the Royesse Iselle, are faced with a sinister curse that hangs like a sword over the entire blighted House of Chalion and all who stand in their circle. Only by employing the darkest, most forbidden of magics, can Cazaril hope to protect his royal charge—an act that will mark the loyal, damaged servant as a tool of the miraculous, and trap him, flesh and soul, in a maze of demonic paradox, damnation, and death.`,
		},
		{
			ID:            14497,
			ISBN:          "0060557818",
			ISBN13:        "9780060557812",
			Title:         "Neverwhere (London Below, #1)",
			TitleNoSeries: "Neverwhere",
			PubDate:       date.New(2003, time.September, 2),
			Series:        Series{Title: "London Below", Position: 1},
			Author:        Author{Name: "Neil Gaiman", ID: 1221698, URL: "https://www.goodreads.com/author/show/1221698"},
			Rating:        4.17,
			URL:           "https://www.goodreads.com/book/show/14497",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1348747943l/14497._SX98_.jpg",
			Description:   `Under the streets of London there's a place most people could never even dream of. A city of monsters and saints, murderers and angels, knights in armour and pale girls in black velvet. This is the city of the people who have fallen between the cracks.<br /><br />Richard Mayhew, a young businessman, is going to find out more than enough about this other London. A single act of kindness catapults him out of his workday existence and into a world that is at once eerily familiar and utterly bizarre. And a strange destiny awaits him down here, beneath his native city: Neverwhere.`,
		},
		{
			ID:            18656030,
			Title:         "Cibola Burn (The Expanse, #4)",
			TitleNoSeries: "Cibola Burn",
			PubDate:       date.New(2014, time.June, 17),
			Series:        Series{Title: "The Expanse", Position: 4},
			Author:        Author{Name: "James S.A. Corey", ID: 4192148, URL: "https://www.goodreads.com/author/show/4192148"},
			Rating:        4.17,
			URL:           "https://www.goodreads.com/book/show/18656030",
			ImageURL:      "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1405023040l/18656030._SX98_.jpg",
			Description:   `<b>The fourth novel in James S.A. Corey’s New York Times bestselling Expanse series</b><br /><br />The gates have opened the way to thousands of habitable planets, and the land rush has begun. Settlers stream out from humanity's home planets in a vast, poorly controlled flood, landing on a new world. Among them, the Rocinante, haunted by the vast, posthuman network of the protomolecule as they investigate what destroyed the great intergalactic society that built the gates and the protomolecule.<br /><br />But Holden and his crew must also contend with the growing tensions between the settlers and the company which owns the official claim to the planet. Both sides will stop at nothing to defend what's theirs, but soon a terrible disease strikes and only Holden - with help from the ghostly Detective Miller - can find the cure.`,
		},
	}
)
