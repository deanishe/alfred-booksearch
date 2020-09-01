// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-07-18

package cli

import (
	"log"
	"path/filepath"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"

	"go.deanishe.net/alfred-booksearch/pkg/gr"
)

// show books by author
func runAuthor() {
	if !authorisedStatus() {
		return
	}

	wf.Var("last_action", "author")
	wf.Var("last_query", opts.Query)

	var (
		books []gr.Book
		key   = "authors/" + cachefileID(opts.AuthorID)
		rerun = wf.IsRunning(booksJob)
	)

	if wf.Cache.Expired(key, opts.MaxCache.Search) {
		rerun = true
		if err := runJob(booksJob, "-savebooks"); err != nil {
			wf.FatalError(err)
		}
	}

	if wf.Cache.Exists(key) {
		checkErr(wf.Cache.LoadJSON(key, &books))
		log.Printf("loaded %d book(s) from cache", len(books))
	} else {
		wf.NewItem("Loading Books…").
			Subtitle("Results will appear momentarily").
			Icon(spinnerIcon())
	}

	// Search for books
	log.Printf("authorName=%q, authorID=%d, sinceLastRequest=%v", opts.AuthorName, opts.AuthorID, time.Since(opts.LastRequestParsed))

	var (
		icons = newIconCache(iconCacheDir)
		mods  = LoadModifiers()
	)

	for _, b := range books {
		bookItem(b, icons, mods)
	}

	addNavActions()

	if !opts.QueryEmpty() {
		wf.Filter(opts.Query)
	}

	wf.WarnEmpty("No Matching Books", "Try a different query?")

	if icons.HasQueue() {
		var err error
		if err = icons.Close(); err == nil {
			err = runJob(iconsJob, "-icons")
		}
		logIfError(err, "cache icons: %v")
	}

	if rerun || wf.IsRunning(iconsJob) {
		wf.Rerun(rerunInterval)
	}

	wf.SendFeedback()
}

// cache books by a given author
func runCacheAuthorList() {
	wf.Configure(aw.TextErrors(true))
	if !opts.Authorised() {
		return
	}

	var (
		key        = "authors/" + cachefileID(opts.AuthorID)
		page       = 1
		pageCount  int
		books, res []gr.Book
		meta       gr.PageData
		last       time.Time
		// Whether to write partial result sets or wait until everything
		// has been downloaded.
		writePartial bool
		err          error
	)
	util.MustExist(filepath.Dir(filepath.Join(wf.CacheDir(), key)))
	log.Printf("[authors] caching books by %q (%d) ...", opts.AuthorName, opts.AuthorID)
	// log.Printf("[authors] cache: %s", key)

	writePartial = !wf.Cache.Exists(key)

	for {
		if pageCount > 0 && page > pageCount {
			break
		}

		if !last.IsZero() && time.Since(last) < time.Second {
			delay := time.Second - time.Since(last)
			log.Printf("[authors] pausing %v till next request ...", delay)
			time.Sleep(delay)
		}
		last = time.Now()

		res, meta, err = api.AuthorBooks(opts.AuthorID, page)
		checkErr(err)

		if pageCount == 0 {
			n := meta.Total
			if n > opts.MaxBooks {
				n = opts.MaxBooks
			}
			pageCount = n / 30
			if n%30 > 0 {
				pageCount++
			}
		}
		books = append(books, res...)
		if writePartial {
			checkErr(wf.Cache.StoreJSON(key, books))
		}
		log.Printf("[authors] cached page %d/%d, %d book(s) for %q", page, pageCount, len(books), opts.AuthorName)
		page++
	}

	checkErr(wf.Cache.StoreJSON(key, books))
}
