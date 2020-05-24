// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-08-08

package cli

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
	"go.deanishe.net/alfred-booksearch/pkg/gr"
)

// show books in a series
func runSeries() {
	updateStatus()
	if !authorisedStatus() {
		return
	}

	id := opts.SeriesID

	if id == 0 {
		var (
			book gr.Book
			key  = "books/" + cachefileID(opts.BookID)
		)
		if !wf.Cache.Exists(key) {
			if !wf.IsRunning(bookJob) {
				checkErr(runJob(bookJob, "-savebook"))
			}
			wf.Rerun(rerunInterval)
			wf.NewItem("Loading Series…").
				Subtitle("Results will appear momentarily").
				Icon(spinnerIcon())
			wf.SendFeedback()
			return
		}

		checkErr(wf.Cache.LoadJSON(key, &book))
		id = book.Series.ID
	}

	if id == 0 {
		wf.Fatal("Book is not part of a series")
	}

	// log.Printf("[series] seriesID=%d", id)
	wf.Var("last_action", "series")
	wf.Var("last_query", opts.Query)
	wf.Var("SERIES_ID", fmt.Sprintf("%d", id))

	var (
		key    = "series/" + cachefileID(id)
		icons  = newIconCache(iconCacheDir)
		mods   = LoadModifiers()
		series gr.Series
		rerun  = wf.IsRunning(seriesJob)
	)

	if wf.Cache.Expired(key, opts.MaxCache.Default) {
		rerun = true
		checkErr(runJob(seriesJob, "-saveseries", fmt.Sprintf("%d", id)))
	}

	if wf.Cache.Exists(key) {
		checkErr(wf.Cache.LoadJSON(key, &series))
		log.Printf("loaded series %q from cache", series.Title)
	} else {
		wf.NewItem("Loading Series…").
			Subtitle("Results will appear momentarily").
			Icon(spinnerIcon())
	}

	log.Printf("[series] %d book(s) in series %q", len(series.Books), series.Title)

	if opts.QueryEmpty() {
		wf.Configure(aw.SuppressUIDs(true))
	}

	for _, b := range series.Books {
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

// save series list to cache
func runCacheSeries() {
	wf.Configure(aw.TextErrors(true))
	if !opts.Authorised() {
		return
	}

	var (
		id, _  = strconv.ParseInt(opts.Query, 10, 64)
		key    = "series/" + cachefileID(id)
		series gr.Series
		err    error
	)

	series, err = api.Series(id)
	checkErr(err)

	util.MustExist(filepath.Dir(filepath.Join(wf.CacheDir(), key)))
	checkErr(wf.Cache.StoreJSON(key, series))
}

func runCacheBook() {
	wf.Configure(aw.TextErrors(true))
	if !opts.Authorised() {
		return
	}
	_, err := bookDetails(opts.BookID)
	checkErr(err)
}
