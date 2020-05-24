// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package cli

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"

	"go.deanishe.net/alfred-booksearch/pkg/gr"
)

// Search Goodreads
func runSearch() {
	updateStatus()
	if !authorisedStatus() {
		return
	}

	if opts.LastRequest != "" {
		data, _ := opts.LastRequestParsed.MarshalText()
		wf.Var("LAST_REQUEST", string(data))
	}

	// Search for books
	log.Printf("query=%q, sinceLastRequest=%v", opts.Query, time.Since(opts.LastRequestParsed))

	if opts.QueryTooShort() {
		wf.NewItem("Query Too Short").
			Subtitle("Keep typing…")
		wf.SendFeedback()
		return
	}

	wf.Var("last_action", "search")
	wf.Var("last_query", opts.Query)

	var (
		icons = newIconCache(iconCacheDir)
		books = cachingSearch(opts.Query)
		mods  = LoadModifiers()
	)

	for _, b := range books {
		bookItem(b, icons, mods)
	}

	wf.WarnEmpty("No Books Found", "Try a different query?")

	if icons.HasQueue() {
		var err error
		if err = icons.Close(); err == nil {
			err = runJob(iconsJob, "-icons")
		}
		logIfError(err, "cache icons: %v")
	}

	if wf.IsRunning(iconsJob) {
		wf.Rerun(rerunInterval)
	}
	wf.SendFeedback()
}

func cachingSearch(query string) (results []gr.Book) {
	key := "queries/" + cachefile(hash(query), ".json")
	reload := func() (interface{}, error) {
		earliest := opts.LastRequestParsed.Add(time.Second * 1)
		now := time.Now()
		if earliest.After(now) {
			d := earliest.Sub(now)
			log.Printf("[throttled] waiting for %v ...", d)
			time.Sleep(d)
		}
		data, _ := time.Now().MarshalText()
		wf.Var("LAST_REQUEST", string(data))
		return api.Search(query)
	}

	util.MustExist(filepath.Dir(filepath.Join(wf.CacheDir(), key)))
	if err := wf.Cache.LoadOrStoreJSON(key, opts.MaxCache.Search, reload, &results); err != nil {
		panic(err)
	}
	return
}

// return an aw.Item for Book.
func bookItem(b gr.Book, icons *iconCache, mods []Modifier) *aw.Item {
	var date, subtitle, rating string

	if !b.PubDate.IsZero() {
		date = fmt.Sprintf(" (%s)", b.PubDate.Format("2006"))
	}
	if b.Rating != 0 {
		rating = fmt.Sprintf(" ⭑ %0.2f", b.Rating)
	}

	subtitle = b.Author.Name + date + rating

	it := wf.NewItem(b.Title).
		Subtitle(subtitle).
		Match(b.Title+" "+b.Author.Name).
		Arg("-script", opts.DefaultScript).
		Copytext(b.Title).
		Valid(true).
		UID(fmt.Sprintf("%d", b.ID)).
		Icon(icons.BookIcon(b)).
		Var("hide_alfred", "true").
		Var("query", opts.Query).
		Var("passvars", "true").
		Var("action", "")

	if b.Description != "" {
		it.Largetype(b.DescriptionText())
	}

	for k, v := range bookVariables(b) {
		it.Var(k, v)
	}

	it.NewModifier(aw.ModCmd).
		Subtitle("All Actions…").
		Arg("-noop").
		Icon(iconMore).
		Var("hide_alfred", "").
		Var("action", "scripts").
		Var("query", "")

	for _, m := range mods {
		it.NewModifier(m.Keys...).
			Subtitle(m.Script.Name).
			Arg("-script", m.Script.Name).
			Icon(m.Script.Icon)
	}

	return it
}
