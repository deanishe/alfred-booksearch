// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-07-18

package gr

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	rxEntityDec  = regexp.MustCompile(`&#\d+;`)
	rxWhitespace = regexp.MustCompile(`\s+`)
	mdPolicy     *bluemonday.Policy
	textPolicy   *bluemonday.Policy
)

func init() {
	mdPolicy = bluemonday.NewPolicy()
	mdPolicy.AllowElements("br", "i", "em", "strong", "b", "p")
	textPolicy = bluemonday.NewPolicy()
	textPolicy.AllowElements("br")
}

// HTML2Markdown converts HTML to Markdown.
func HTML2Markdown(s string) string {
	tags := []struct {
		find, repl string
	}{
		{"<p>", ""},
		{"</p>", "\n\n"},
		{"<br /><br />", "\n\n"},
		{"<br/><br/>", "\n\n"},
		{"<br />", "  \n"},
		{"<br/>", "  \n"},
		{"<i>", "*"},
		{"</i>", "*"},
		{"<em>", "*"},
		{"</em>", "*"},
		{"<strong>", "**"},
		{"</strong>", "**"},
		{"<b>", "**"},
		{"</b>", "**"},
	}
	s = mdPolicy.Sanitize(s)
	s = decodeEntities(tidyText(s))
	for _, t := range tags {
		s = strings.Replace(s, t.find, t.repl, -1)
	}
	return s
}

var rxBR = regexp.MustCompile(`<br ?/>`)

// HTML2Text converts HTML to plaintext.
func HTML2Text(s string) string {
	s = textPolicy.Sanitize(s)
	s = decodeEntities(tidyText(s))
	return rxBR.ReplaceAllString(s, "\n")
}

// convert HTML &#NN; entities to text.
func decodeEntities(s string) string {
	for {
		m := rxEntityDec.FindStringIndex(s)
		if m == nil {
			break
		}
		i, _ := strconv.ParseInt(s[m[0]+2:m[1]-1], 10, 64)
		s = s[0:m[0]] + string(i) + s[m[1]:]
	}
	return strings.TrimSpace(s)
}

// replace ASCII dashes and ellipses with Unicode versions; collapse whitespace.
func tidyText(s string) string {
	for _, pat := range []string{"....", ". . . .", ". . .", "..."} {
		s = strings.Replace(s, pat, "…", -1)
	}
	s = strings.Replace(s, "--", "—", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = rxWhitespace.ReplaceAllString(s, " ")
	return s
}
