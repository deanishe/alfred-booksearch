// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-07-18

package cli

import (
	"fmt"
	"log"
	"sort"
	"time"

	aw "github.com/deanishe/awgo"

	"go.deanishe.net/alfred-booksearch/pkg/gr"
	"go.deanishe.net/fuzzy"
)

const shelvesKey = "shelves.json"

// Show books on shelf
func runShelf() {
	updateStatus()
	if !authorisedStatus() {
		return
	}

	var (
		shelf gr.Shelf
		key   = "shelves/" + opts.ShelfName + ".json"
		rerun = wf.IsRunning(shelfJob)
	)

	if wf.Cache.Expired(key, opts.MaxCache.Shelf) {
		rerun = true
		checkErr(runJob(shelfJob, "-saveshelf"))
	}

	if wf.Cache.Exists(key) {
		checkErr(wf.Cache.LoadJSON(key, &shelf))
	} else {
		wf.NewItem("Loading Books…").
			Subtitle("Results will appear momentarily").
			Icon(spinnerIcon())
	}

	wf.Var("last_action", "shelf")
	wf.Var("hide_alfred", "")
	wf.Var("last_query", opts.Query)

	var (
		icons = newIconCache(iconCacheDir)
		mods  = LoadModifiers()
	)

	log.Printf("query=%q", opts.Query)

	// show books in list order if there's no query
	if opts.QueryEmpty() {
		wf.Configure(aw.SuppressUIDs(true))
	}

	for _, b := range shelf.Books {
		it := bookItem(b, icons, mods)

		it.NewModifier(aw.ModCtrl).
			Subtitle("Remove from Shelf").
			Arg("-remove", shelf.Name).
			Icon(iconDelete).
			Var("action", "shelf").
			Var("query", opts.Query).
			Var("passvars", "true")
	}

	// add alternate actions
	if len(opts.Query) > 2 {
		wf.NewItem("Reload").
			Subtitle("Reload shelf from server").
			Arg("-reload").
			UID("reload").
			Icon(iconReload).
			Valid(true).
			Var("last_query", "")

		// wf.NewItem("Shelves").
		// 	Subtitle("Go back to shelves list").
		// 	Match("shelves").
		// 	Arg("-noop").
		// 	UID("shelves").
		// 	Icon(iconBook).
		// 	Valid(true).
		// 	Var("action", "shelves").
		// 	Var("last_query", "").
		// 	Var("last_action", "")
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

// Show user's shelves
func runShelves() {
	updateStatus()
	if !authorisedStatus() {
		return
	}

	if opts.UserID == 0 {
		wf.NewItem("User ID is not Set").
			Subtitle("Please deauthorise and re-authorise the workflow").
			Valid(false).
			Icon(iconError)
		wf.SendFeedback()
		return
	}

	wf.Var("last_action", "shelves")
	wf.Var("last_query", opts.Query)

	var (
		shelves []gr.Shelf
		rerun   = wf.IsRunning(shelvesJob)
	)

	if wf.Cache.Expired(shelvesKey, opts.MaxCache.Shelf) {
		rerun = true
		checkErr(runJob(shelvesJob, "-saveshelves"))
	}

	if wf.Cache.Exists(shelvesKey) {
		checkErr(wf.Cache.LoadJSON(shelvesKey, &shelves))
		log.Printf("loaded %d shelves from cache", len(shelves))
	} else {
		wf.NewItem("Loading Shelves…").
			Subtitle("Results will appear momentarily").
			Icon(spinnerIcon())
	}

	if opts.QueryEmpty() {
		wf.Configure(aw.SuppressUIDs(true))
	}

	for _, shelf := range shelves {
		id := fmt.Sprintf("%d", shelf.ID)
		it := wf.NewItem(shelf.Title()).
			Subtitle(fmt.Sprintf("%d book(s)", shelf.Size)).
			UID(id).
			Valid(true).
			Icon(iconShelf).
			Var("SHELF_ID", id).
			Var("SHELF_NAME", shelf.Name).
			Var("SHELF_TITLE", shelf.Title()).
			Var("action", "shelf").
			Var("passvars", "true")

		it.NewModifier(aw.ModCmd).
			Subtitle(fmt.Sprintf("Open “%s” on goodreads.com", shelf.Title())).
			Valid(true).
			Arg("-open", shelf.URL).
			Var("action", "").
			Var("hide_alfred", "true")
	}

	// add alternate actions
	if len(opts.Query) > 2 {
		wf.NewItem("Reload").
			Subtitle("Reload shelves from server").
			Arg("-reloadshelves").
			UID("reload").
			Icon(iconReload).
			Valid(true).
			Var("action", "shelves").
			Var("hide_alfred", "")
	}

	addNavActions("shelves")

	if !opts.QueryEmpty() {
		wf.Filter(opts.Query)
	}

	log.Printf("query=%q", opts.Query)

	if rerun {
		wf.Rerun(rerunInterval)
	}
	wf.WarnEmpty("No Matching Lists", "Try a different query?")
	wf.SendFeedback()
}

// add book to shelf
func runAddToShelves() {
	wf.Configure(aw.TextErrors(true))
	if !opts.Authorised() {
		return
	}
	log.Printf("adding book %d to shelves %v", opts.BookID, opts.Args)
	if err := api.AddToShelves(opts.BookID, opts.Args); err != nil {
		notifyError("Add to Shelves Failed", err)
		log.Fatalf("[ERROR] add to shelf %q: %v", opts.Query, err)
	}

	msg := "Added to 1 shelf"
	if len(opts.Args) > 1 {
		msg = fmt.Sprintf("Added to %d shelves", len(opts.Args))
	}
	v := aw.NewArgVars()
	v.Var("notification_title", opts.BookTitle).
		Var("notification_text", msg)

	// deselect shelves
	for _, s := range opts.Args {
		v.Var("shelf_"+s, "")
	}
	checkErr(v.Send())
}

// remove book from shelf
func runRemoveFromShelf() {
	wf.Configure(aw.TextErrors(true))
	if !opts.Authorised() {
		return
	}
	log.Printf("removing book %d from shelf %q", opts.BookID, opts.Query)
	if err := api.RemoveFromShelf(opts.BookID, opts.Query); err != nil {
		notifyError("Remove from Shelf Failed", err)
		log.Fatalf("[ERROR] remove from shelf %q: %v", opts.Query, err)
	}

	title := opts.ShelfTitle
	if title == "" {
		title = opts.Query
	}
	notify(opts.BookTitle, fmt.Sprintf("Removed from “%s”", title))

	// remove book from cache
	var (
		key     = "shelves/" + opts.ShelfName + ".json"
		cleaned []gr.Book
		shelf   gr.Shelf
	)
	if !wf.Cache.Exists(key) {
		return
	}
	if err := wf.Cache.LoadJSON(key, &shelf); err != nil {
		log.Fatalf("[ERROR] load cached shelf: %v", err)
	}

	for _, b := range shelf.Books {
		if b.ID != opts.BookID {
			cleaned = append(cleaned, b)
		}
	}
	shelf.Books = cleaned
	if err := wf.Cache.StoreJSON(key, shelf); err != nil {
		log.Fatalf("[ERROR] cache shelf: %v", err)
	}
}

// update cached shelf
func runReloadShelf() {
	wf.Configure(aw.TextErrors(true))
	if !opts.Authorised() {
		return
	}
	checkErr(runJob(shelfJob, "-saveshelf"))
}

// update cached shelves
func runReloadShelves() {
	wf.Configure(aw.TextErrors(true))
	if !opts.Authorised() {
		return
	}
	checkErr(runJob(shelvesJob, "-saveshelves"))
}

// choose shelves to add a book to
func runSelectShelves() {
	if opts.FlagSelectShelf {
		wf.Configure(aw.TextErrors(true))
		var (
			key   = fmt.Sprintf("shelf_" + opts.Query)
			value = "true"
		)
		if wf.Config.GetBool(key) {
			value = "false"
		}
		log.Printf("[shelves] shelf=%s, selected=%s", opts.Query, value)
		checkErr(aw.NewArgVars().Var("shelf_"+opts.Query, value).Send())
		return
	}

	updateStatus()
	if !authorisedStatus() {
		return
	}

	var (
		shelves []gr.Shelf
		rerun   = wf.IsRunning(shelvesJob)
	)

	if wf.Cache.Expired(shelvesKey, opts.MaxCache.Shelf) {
		rerun = true
		checkErr(runJob(shelvesJob, "-saveshelves"))
	}

	if wf.Cache.Exists(shelvesKey) {
		checkErr(wf.Cache.LoadJSON(shelvesKey, &shelves))
		log.Printf("loaded %d shelves from cache", len(shelves))
	} else {
		wf.NewItem("Loading Shelves…").
			Subtitle("Results will appear momentarily").
			Icon(spinnerIcon())
	}

	wf.Var("hide_alfred", "").Var("passvars", "true")

	var (
		selected   = selectShelves(shelves)
		args       = append([]string{"-add"}, selected...)
		lastAction = wf.Config.Get("last_action")
		lastQuery  = wf.Config.Get("last_query")
		msg        = "Add to 1 shelf"
	)

	if len(selected) != 1 {
		msg = fmt.Sprintf("Add to %d shelves", len(selected))
	}

	if !opts.QueryEmpty() {
		shelves = filterShelves(shelves, opts.Query)
	}

	for _, shelf := range shelves {
		icon := iconShelf
		if shelf.Selected {
			icon = iconShelfSelected
		}

		it := wf.NewItem(shelf.Title()).
			Subtitle(fmt.Sprintf("%d book(s)", shelf.Size)).
			Arg("-selection", "-select", shelf.Name).
			Valid(true).
			Icon(icon).
			Var("action", "select").
			Var("query", "")

		if len(selected) > 0 {
			it.NewModifier(aw.ModCmd).
				Subtitle(msg).
				Valid(true).
				Arg(args...).
				Icon(iconSave).
				Var("action", lastAction).
				Var("query", lastQuery).
				Var("last_query", "").
				Var("last_action", "")
		}
	}

	log.Printf("query=%q", opts.Query)

	if rerun {
		wf.Rerun(rerunInterval)
	}
	wf.WarnEmpty("No Matching Lists", "Try a different query?")
	wf.SendFeedback()
}

// cache a specific shelf
func runCacheShelf() {
	wf.Configure(aw.TextErrors(true))
	if !opts.Authorised() {
		return
	}

	var (
		key       = "shelves/" + opts.ShelfName + ".json"
		page      = 1
		pageCount int
		shelf     = gr.Shelf{ID: opts.ShelfID, Name: opts.ShelfName}
		books     []gr.Book
		meta      gr.PageData
		last      time.Time
		err       error

		writePartial = !wf.Cache.Exists(key)
	)

	log.Printf("[shelves] fetching shelf %q ...", opts.ShelfName)

	for {
		if pageCount > 0 && page > pageCount {
			break
		}

		if !last.IsZero() && time.Since(last) < time.Second {
			delay := time.Second - time.Since(last)
			log.Printf("[shelves] pausing %v till next request ...", delay)
			time.Sleep(delay)
		}
		last = time.Now()

		books, meta, err = api.UserShelf(opts.UserID, opts.ShelfName, page)
		checkErr(err)

		if pageCount == 0 {
			pageCount = meta.Total / 50
			if meta.Total%50 > 0 {
				pageCount++
			}
		}

		shelf.Books = append(shelf.Books, books...)
		shelf.Size = meta.Total
		if writePartial {
			checkErr(wf.Cache.StoreJSON(key, shelf))
		}
		log.Printf("[shelves] cached page %d/%d, %d book(s)", page, pageCount, len(books))
		page++
	}

	checkErr(wf.Cache.StoreJSON(key, shelf))
}

// cache list of user's shelves
func runCacheShelves() {
	wf.Configure(aw.TextErrors(true))
	if !opts.Authorised() {
		return
	}

	var (
		page         = 1
		pageCount    int
		shelves, res []gr.Shelf
		meta         gr.PageData
		last         time.Time
		writePartial = !wf.Cache.Exists(shelvesKey)
		err          error
	)

	log.Println("[shelves] fetching users shelves ...")

	for {
		if pageCount > 0 && page > pageCount {
			break
		}

		if !last.IsZero() && time.Since(last) < time.Second {
			delay := time.Second - time.Since(last)
			log.Printf("[shelves] pausing %v till next request ...", delay)
			time.Sleep(delay)
		}
		last = time.Now()

		res, meta, err = api.UserShelves(opts.UserID, page)
		checkErr(err)

		if pageCount == 0 {
			pageCount = meta.Total / 15
			if meta.Total%15 > 0 {
				pageCount++
			}
		}

		shelves = append(shelves, res...)
		if writePartial {
			checkErr(wf.Cache.StoreJSON(shelvesKey, shelves))
		}
		log.Printf("[shelves] cached page %d/%d, %d shelves", page, pageCount, len(shelves))
		page++
	}

	checkErr(wf.Cache.StoreJSON(shelvesKey, shelves))
}

// bySelection sorts shelves by selection status
type bySelection []gr.Shelf

// Implement sort.Interface
func (s bySelection) Len() int      { return len(s) }
func (s bySelection) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s bySelection) Less(i, j int) bool {
	a, b := s[i], s[j]
	if b.Selected && !a.Selected {
		return true
	}
	return false
}
func (s bySelection) Keywords(i int) string { return s[i].Title() }

func filterShelves(shelves []gr.Shelf, query string) []gr.Shelf {
	groups := make([][]gr.Shelf, 2)
	for i, r := range fuzzy.New(bySelection(shelves)).Sort(query) {
		if !r.Match {
			continue
		}
		s := shelves[i]
		var n int
		if s.Selected {
			n = 1
		}
		groups[n] = append(groups[n], s)
	}
	var matches []gr.Shelf
	for _, g := range groups {
		matches = append(matches, g...)
	}
	return matches
}

func selectShelves(shelves []gr.Shelf) (names []string) {
	for i, s := range shelves {
		s.Selected = wf.Config.GetBool("shelf_" + s.Name)
		log.Printf("[shelves] name=%q, selected=%v", s.Name, s.Selected)
		shelves[i] = s
		if s.Selected {
			names = append(names, s.Name)
		}
	}
	sort.Stable(bySelection(shelves))
	return
}
