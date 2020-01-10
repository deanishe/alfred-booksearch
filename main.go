// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/keychain"
	"github.com/deanishe/awgo/update"
	"github.com/deanishe/awgo/util"
)

const (
	repo            = "deanishe/alfred-goodreads"
	readmeURL       = "https://github.com/deanishe/alfred-goodreads/"
	issueTrackerURL = "https://github.com/deanishe/alfred-goodreads/issues"
	apiKeyURL       = "https://www.goodreads.com/api/keys"
	updateJob       = "update"
	iconsJob        = "icons"
	booksJob        = "booklist"
)

var (
	iconCacheDir      string
	searchCacheDir    string
	maxCacheAge       = 24 * time.Hour
	minCacheAge       = 1 * time.Minute
	maxIconCacheAge   = 672 * time.Hour // 28 days
	maxBooksPerAuthor = 100
	minBooksPerAuthor = 30

	iconAuthor          = &aw.Icon{Value: "icons/author.png"}
	iconBook            = &aw.Icon{Value: "icons/book.png"}
	iconDelete          = &aw.Icon{Value: "icons/delete.png"}
	iconDocs            = &aw.Icon{Value: "icons/docs.png"}
	iconError           = &aw.Icon{Value: "icons/error.png"}
	iconHelp            = &aw.Icon{Value: "icons/help.png"}
	iconIssue           = &aw.Icon{Value: "icons/issue.png"}
	iconLink            = &aw.Icon{Value: "icons/link.png"}
	iconOK              = &aw.Icon{Value: "icons/ok.png"}
	iconSpinner         = &aw.Icon{Value: "icons/spinner.png"}
	iconUpdateAvailable = &aw.Icon{Value: "icons/update-available.png"}
	iconUpdateOK        = &aw.Icon{Value: "icons/update-ok.png"}
	iconURL             = &aw.Icon{Value: "icons/url.png"}
	iconWarning         = &aw.Icon{Value: "icons/warning.png"}
	// iconDefault         = &aw.Icon{Value: "icon.png"}

	wf   *aw.Workflow
	opts = &options{
		MaxBooks:          maxBooksPerAuthor,
		MaxCacheAge:       maxCacheAge,
		LastRequestParsed: time.Time{},
	}
)

type options struct {
	MaxBooks       int           // How many books to load per author
	MaxCacheAge    time.Duration // How long search results are cached
	APIKey         string        // Goodreads API key
	MinQueryLength int           // Minimum length of search query

	// Time of last request to Goodreads API.
	// Requests are throttled to 1/sec.
	LastRequest       string
	LastRequestParsed time.Time

	// Workflow data
	BookID     string
	BookTitle  string
	AuthorID   string
	AuthorName string

	// Alternate actions
	FlagAuthor      bool `env:"-"`
	FlagCacheAuthor bool `env:"-"`
	FlagIcons       bool `env:"-"`
	FlagCheck       bool `env:"-"`
	FlagConf        bool `env:"-"`
	FlagAPIKey      bool `env:"-"`
	FlagDelKey      bool `env:"-"`
	FlagSaveKey     bool `env:"-"`
	FlagHelp        bool `env:"-"`

	// Search query. Populated from first argument.
	Query string `env:"-"`
}

// QueryEmpty returns true if query is empty.
func (opts *options) QueryEmpty() bool { return strings.TrimSpace(opts.Query) == "" }

// QueryTooShort returns true if query is empty.
func (opts *options) QueryTooShort() bool {
	return strings.TrimSpace(opts.Query) == ""
}

func init() {
	aw.IconError = iconError
	aw.IconWarning = iconWarning

	flag.BoolVar(&opts.FlagAuthor, "author", false, "list books for author")
	flag.BoolVar(&opts.FlagCacheAuthor, "savebooks", false, "cache all books by author")
	flag.BoolVar(&opts.FlagCheck, "check", false, "check for a new version")
	flag.BoolVar(&opts.FlagConf, "conf", false, "show workflow configuration")
	flag.BoolVar(&opts.FlagIcons, "icons", false, "download queued icons")
	flag.BoolVar(&opts.FlagAPIKey, "apikey", false, "set Goodreads API key")
	flag.BoolVar(&opts.FlagSaveKey, "savekey", false, "save API key to Keychain")
	flag.BoolVar(&opts.FlagSaveKey, "delkey", false, "delete API key from Keychain")
	flag.BoolVar(&opts.FlagHelp, "h", false, "show this message and exit")

	wf = aw.New(
		aw.HelpURL(issueTrackerURL),
		update.GitHub(repo),
	)
	iconCacheDir = filepath.Join(wf.CacheDir(), "icons")
	searchCacheDir = filepath.Join(wf.CacheDir(), "queries")
}

func main() {
	wf.Run(run)
}

func runHelp() {
	wf.Configure(aw.TextErrors(true))
	flag.Usage()
}

func runAPIKey() {
	if opts.Query == "" {
		wf.NewItem("Enter API Key").
			Subtitle("Enter your Goodreads.com API key")
	} else {
		wf.NewItem(fmt.Sprintf("Set API Key to %q", opts.Query)).
			Subtitle("↩	to save Goodreads.com API key").
			Arg(opts.Query).
			Valid(true).
			Var("command", "setkey")
	}

	wf.NewItem("Goodreads.com API Keys").
		Subtitle("Open Goodreads API page in browser").
		Arg(apiKeyURL).
		Valid(true).
		Icon(iconURL).
		Var("command", "open")

	wf.SendFeedback()
}

func runSaveKey() {
	wf.Configure(aw.TextErrors(true))
	var (
		key = opts.Query
		v   = aw.NewArgVars()
		kc  = keychain.New(wf.BundleID())
		err error
	)

	if err = kc.Set("api_key", key); err != nil {
		wf.Fatalf("add to keychain: %v", err)
	}

	v.Var("command", "config")
	v.Var("API_KEY", key)
	_ = v.Send()

	log.Printf("keychain: saved API key")
}

func runDelKey() {
	wf.Configure(aw.TextErrors(true))
	var (
		v   = aw.NewArgVars()
		kc  = keychain.New(wf.BundleID())
		err error
	)

	if err = kc.Delete("api_key"); err != nil {
		if err != keychain.ErrNotFound {
			wf.Fatalf("delete from keychain: %v", err)
		}
	}

	v.Var("command", "config")
	v.Var("API_KEY", "")
	_ = v.Send()

	log.Printf("keychain: deleted API key")
}

func runCacheBookList() {
	wf.Configure(aw.TextErrors(true))

	var (
		key        = "authors/" + cachefile(hash(opts.AuthorID), ".json")
		page       = 1
		pageCount  int
		books, res []Book
		meta       pageData
		last       time.Time
		// Whether to write partial result sets or wait until everything
		// has been downloaded.
		writePartial bool
		err          error
	)
	util.MustExist(filepath.Dir(filepath.Join(wf.CacheDir(), key)))
	log.Printf("caching books by %q (%s) ...", opts.AuthorName, opts.AuthorID)
	log.Printf("cache: %s", key)

	writePartial = !wf.Cache.Exists(key)

	for {

		if pageCount > 0 && page > pageCount {
			break
		}

		if !last.IsZero() && time.Now().Sub(last) < time.Second {
			delay := time.Second - time.Now().Sub(last)
			log.Printf("pausing %v till next request ...", delay)
			time.Sleep(delay)
		}
		last = time.Now()

		res, meta, err = authorBooks(opts.AuthorID, opts.APIKey, page)
		if err != nil {
			wf.FatalError(err)
		}
		if pageCount == 0 {
			n := meta.Total
			if n > opts.MaxBooks {
				n = opts.MaxBooks
			}
			pageCount = n / 30
			r := n % 30
			if r > 0 {
				pageCount++
			}
		}
		books = append(books, res...)
		if writePartial {
			if err := wf.Cache.StoreJSON(key, books); err != nil {
				wf.FatalError(err)
			}
		}
		log.Printf("cached page %d/%d, %d book(s) for %q", page, pageCount, len(books), opts.AuthorName)
		page++
	}

	if err := wf.Cache.StoreJSON(key, books); err != nil {
		wf.FatalError(err)
	}
}

func runAuthor() {
	var (
		books   []Book
		key     = "authors/" + cachefile(hash(opts.AuthorID), ".json")
		running = wf.IsRunning(booksJob)
		rerun   = running
	)

	if wf.Cache.Expired(key, opts.MaxCacheAge) {
		rerun = true
		if !running {
			if err := wf.RunInBackground(booksJob, exec.Command(os.Args[0], "-savebooks")); err != nil {
				wf.FatalError(err)
			}
		}
	}

	if wf.Cache.Exists(key) {
		if err := wf.Cache.LoadJSON(key, &books); err != nil {
			wf.FatalError(err)
		}
		log.Printf("loaded %d book(s) from cache", len(books))
	} else {
		wf.NewItem("Loading Books…").
			Subtitle("Results will appear momentarily").
			Icon(spinnerIcon())
	}

	// Search for books
	log.Printf("authorName=%q, authorID=%s, sinceLastRequest=%v", opts.AuthorName, opts.AuthorID, time.Now().Sub(opts.LastRequestParsed))

	icons := newIconCache(iconCacheDir)
	mods := LoadModifiers()

	for i, b := range books {
		var (
			title    = b.Title
			subtitle = fmt.Sprintf("%s (%s) ⭑ %0.2f", b.Author.Name, b.PubDate.Format("2006"), b.Rating)
			icon     = icons.BookIcon(b)
			authorID = fmt.Sprintf("%d", b.Author.ID)
		)

		it := wf.NewItem(title).
			Subtitle(subtitle).
			Arg(b.URL).
			Copytext(b.Title).
			Valid(true).
			UID(fmt.Sprintf("%d", b.ID)).
			Icon(icon).
			Var("AUTHOR_ID", authorID)

		it.NewModifier(aw.ModCmd).
			Subtitle(fmt.Sprintf("View author “%s” on Goodreads.com", b.Author.Name)).
			Arg(b.Author.URL).
			Icon(iconAuthor)

		for _, m := range mods {
			keys, value := m.For(b)
			subtitle := m.Name
			if subtitle == "" {
				subtitle = value
			}
			it.NewModifier(keys...).
				Subtitle(subtitle).
				Arg(value).
				Icon(iconLink)
		}

		log.Printf("[%2d/%2d] %q by %s (id=%d)", i+1, len(books), b.Title, b.Author.Name, b.ID)
	}

	if !opts.QueryEmpty() {
		wf.Filter(opts.Query)
	}

	wf.WarnEmpty("No Matching Books", "Try a different query?")

	if icons.HasQueue() {
		if err := icons.Close(); err != nil {
			log.Printf("[ERROR] save icons: %v", err)
		} else if !wf.IsRunning(iconsJob) {
			if err := wf.RunInBackground(iconsJob, exec.Command(os.Args[0], "-icons")); err != nil {
				log.Printf("[ERROR] cache icons: %v", err)
			} else {
				rerun = true
			}
		}
	}

	if rerun || wf.IsRunning(iconsJob) {
		wf.Rerun(0.1)
	}

	wf.SendFeedback()
}

func runConfig() {
	title := "Workflow Is Up To Date"
	subtitle := "↩ or ⇥ to check for update"
	icon := iconUpdateOK
	if wf.UpdateAvailable() {
		title = "Workflow Update Available"
		subtitle = "↩ or ⇥ to install"
		icon = iconUpdateAvailable
	}

	wf.NewItem(title).
		Subtitle(subtitle).
		Valid(false).
		Autocomplete("workflow:update").
		Icon(icon)

	title = "API Key"
	subtitle = "Goodreads API key is set"
	icon = iconOK
	if opts.APIKey == "" {
		title = "API Key Not Set"
		subtitle = "↩ to set up API key"
		icon = iconError
	}

	it := wf.NewItem(title).
		Subtitle(subtitle).
		Valid(true).
		Icon(icon)

	it.Var("command", "apikey")
	if opts.APIKey != "" {
		it.NewModifier(aw.ModOpt).
			Subtitle("Delete API Key").
			Icon(iconDelete).
			Var("command", "delkey")
	}

	it.NewModifier(aw.ModCmd).
		Subtitle("Open Goodreads API page in browser").
		Valid(true).
		Icon(iconURL)

	wf.NewItem("Open Docs").
		Subtitle("Open workflow documentation in your browser").
		Arg(readmeURL).
		Valid(true).
		Icon(iconDocs).
		Var("command", "open")

	wf.NewItem("Get Help").
		Subtitle("Open workflow issue tracker in your browser").
		Arg(issueTrackerURL).
		Valid(true).
		Icon(iconHelp).
		Var("command", "open")

	wf.NewItem("Report Bug").
		Subtitle("Open workflow issue tracker in your browser").
		Arg(issueTrackerURL).
		Valid(true).
		Icon(iconIssue).
		Var("command", "open")

	if !opts.QueryEmpty() {
		_ = wf.Filter(opts.Query)
	}
	wf.WarnEmpty("No Matches", "Try a different query?")

	wf.SendFeedback()
}

// Check for workflow update + clear stale cache files.
func runCheck() {
	wf.Configure(aw.TextErrors(true))

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		log.Println("checking for updates...")
		if err := wf.CheckForUpdate(); err != nil {
			log.Printf("[ERROR] update check: %v", err)
		}
	}()

	clean := func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		// Delete cached queries (.json) and covers (.png).
		x := filepath.Ext(fi.Name())
		if x != ".json" && x != ".png" {
			return nil
		}

		age := time.Now().Sub(fi.ModTime())
		max := opts.MaxCacheAge
		if x == ".png" {
			max = maxIconCacheAge
		}
		if age > max {
			log.Printf("deleting %q (%v) ...", util.PrettyPath(p), age)
			if err := os.Remove(p); err != nil {
				return errors.Wrap(err, util.PrettyPath(p))
			}
		}
		return nil
	}

	go func() {
		defer wg.Done()
		log.Println("cleaning query cache...")

		if err := filepath.Walk(searchCacheDir, clean); err != nil {
			log.Printf("[ERROR] clean query cache: %v", err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		log.Println("cleaning icon cache...")

		if err := filepath.Walk(iconCacheDir, clean); err != nil {
			log.Printf("[ERROR] clean icon cache: %v", err)
			return
		}
	}()

	wg.Wait()
}

func runIcons() {
	wf.Configure(aw.TextErrors(true))
	icons := newIconCache(iconCacheDir)
	if !icons.HasQueue() {
		log.Printf("nothing to download")
		return
	}
	if err := icons.ProcessQueue(); err != nil {
		log.Fatal(err)
	}
}

func hash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func prepareOpts() {
	wf.Args() // call to handle magic actions
	flag.Parse()

	if err := wf.Config.To(opts); err != nil {
		log.Printf("[ERROR] load configuration: %v", err)
	}

	if args := flag.Args(); len(args) > 0 {
		opts.Query = strings.TrimSpace(args[0])
	}

	// Ensure sensible minimums
	if opts.MaxBooks < minBooksPerAuthor {
		opts.MaxBooks = minBooksPerAuthor
	}
	if opts.MaxCacheAge < minCacheAge {
		opts.MaxCacheAge = minCacheAge
	}

	if opts.MinQueryLength == 0 {
		opts.MinQueryLength = 3
	}

	// Try to read API key from Keychain
	if opts.APIKey == "" {
		kc := keychain.New(wf.BundleID())
		if s, err := kc.Get("api_key"); err != nil {
			if err != keychain.ErrNotFound {
				log.Printf("[ERROR] keychain: %v", err)
			}
		} else {
			opts.APIKey = s
		}
	}
	if opts.LastRequest != "" {
		if err := opts.LastRequestParsed.UnmarshalText([]byte(opts.LastRequest)); err != nil {
			log.Printf("[ERROR] invalid LastRequest %q: %v", opts.LastRequest, err)
		}
	}
}

func run() {
	util.MustExist(iconCacheDir)
	util.MustExist(searchCacheDir)

	prepareOpts()

	if opts.FlagAuthor {
		runAuthor()
		return
	}

	if opts.FlagCacheAuthor {
		runCacheBookList()
		return
	}

	if opts.FlagAPIKey {
		runAPIKey()
		return
	}

	if opts.FlagSaveKey {
		runSaveKey()
		return
	}

	if opts.FlagDelKey {
		runDelKey()
		return
	}

	if opts.FlagHelp {
		runHelp()
		return
	}

	if opts.FlagCheck {
		runCheck()
		return
	}

	if opts.FlagConf {
		runConfig()
		return
	}

	if opts.FlagIcons {
		runIcons()
		return
	}

	if opts.LastRequest != "" {
		data, _ := opts.LastRequestParsed.MarshalText()
		wf.Var("LAST_REQUEST", string(data))
	}

	// Search Goodreads
	if wf.UpdateAvailable() && opts.Query == "" {
		wf.NewItem("Update Available!").
			Subtitle("↩ or ⇥ to install update").
			Valid(false).
			Autocomplete("workflow:update").
			Icon(iconUpdateAvailable)
	}

	if wf.UpdateCheckDue() && !wf.IsRunning(updateJob) {
		if err := wf.RunInBackground(updateJob, exec.Command(os.Args[0], "-check")); err != nil {
			log.Printf("[ERROR] check for update: %v", err)
		}
	}

	if opts.APIKey == "" {
		wf.NewItem("API Key Not Set").
			Subtitle("Please set a Goodreads API key in the workflow configuration").
			Valid(true).
			Icon(iconWarning).
			Var("command", "apikey")

		wf.SendFeedback()
		return
	}

	// Search for books
	log.Printf("query=%q, sinceLastRequest=%v", opts.Query, time.Now().Sub(opts.LastRequestParsed))

	if opts.QueryTooShort() {
		return
	}

	icons := newIconCache(iconCacheDir)
	books := cachingSearch(opts.Query)
	mods := LoadModifiers()

	for i, b := range books {
		var (
			title    = b.Title
			subtitle = fmt.Sprintf("%s (%s) ⭑ %0.2f", b.Author.Name, b.PubDate.Format("2006"), b.Rating)
			icon     = icons.BookIcon(b)
			authorID = fmt.Sprintf("%d", b.Author.ID)
		)

		it := wf.NewItem(title).
			Subtitle(subtitle).
			Arg(b.URL).
			Copytext(b.Title).
			Valid(true).
			UID(fmt.Sprintf("%d", b.ID)).
			Icon(icon).
			Var("BOOK_ID", fmt.Sprintf("%d", b.ID)).
			Var("BOOK_TITLE", b.Title).
			Var("AUTHOR_ID", authorID).
			Var("AUTHOR_NAME", b.Author.Name)

		it.NewModifier(aw.ModCmd).
			Subtitle(fmt.Sprintf("View author “%s” on Goodreads.com", b.Author.Name)).
			Arg(b.Author.URL).
			Icon(iconAuthor)

		it.NewModifier(aw.ModOpt).
			Subtitle(fmt.Sprintf("View books by “%s”", b.Author.Name)).
			Arg(authorID).
			Icon(iconAuthor).
			Var("command", "author")

		for _, m := range mods {
			keys, value := m.For(b)
			subtitle := m.Name
			if subtitle == "" {
				subtitle = value
			}
			it.NewModifier(keys...).
				Subtitle(subtitle).
				Arg(value).
				Icon(iconLink)
		}

		log.Printf("[%2d/%2d] %q by %s (id=%d)", i+1, len(books), b.Title, b.Author.Name, b.ID)
	}

	wf.WarnEmpty("No Books Found", "Try a different query?")

	var rerun bool
	if icons.HasQueue() {
		if err := icons.Close(); err != nil {
			log.Printf("[ERROR] save icons: %v", err)
		} else if !wf.IsRunning(iconsJob) {
			if err := wf.RunInBackground(iconsJob, exec.Command(os.Args[0], "-icons")); err != nil {
				log.Printf("[ERROR] cache icons: %v", err)
			} else {
				rerun = true
			}
		}
	}

	if rerun || wf.IsRunning(iconsJob) {
		wf.Rerun(0.2)
	}
	wf.SendFeedback()
}

func cachingSearch(query string) (results []Book) {
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
		return search(query, opts.APIKey)
	}

	util.MustExist(filepath.Dir(filepath.Join(wf.CacheDir(), key)))
	if err := wf.Cache.LoadOrStoreJSON(key, maxCacheAge, reload, &results); err != nil {
		wf.FatalError(err)
	}
	return
}

/*
func cachingAuthor(id string) (results []Book) {
	key := "queries/" + cachefile(hash(id), ".json")
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
		return authorBooks(id, opts.APIKey)
	}

	util.MustExist(filepath.Dir(filepath.Join(wf.CacheDir(), key)))
	if err := wf.Cache.LoadOrStoreJSON(key, maxCacheAge, reload, &results); err != nil {
		wf.FatalError(err)
	}
	return
}
*/
