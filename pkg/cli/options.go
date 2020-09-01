// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-07-14

package cli

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/deanishe/awgo/keychain"
	"github.com/pkg/errors"
)

const (
	// defaults
	minCacheAge       = 3 * time.Minute
	maxBooksPerAuthor = 100
	minBooksPerAuthor = 30
)

var (
	fs   *flag.FlagSet
	opts *options
)

func init() {
	opts = &options{
		MaxBooks:          maxBooksPerAuthor,
		LastRequestParsed: time.Time{},
	}
	// default cache values
	opts.MaxCache.Default = 24 * time.Hour
	opts.MaxCache.Search = 12 * time.Hour
	opts.MaxCache.Icons = 336 * time.Hour // 14 days
	opts.MaxCache.Shelf = 5 * time.Minute
	opts.MaxCache.Feeds = 90 * time.Minute
}

type options struct {
	MaxBooks int      // How many books to load per author
	MaxCache struct { // How long data are cached
		Default time.Duration
		Search  time.Duration
		Shelf   time.Duration
		Icons   time.Duration
		Feeds   time.Duration
	}
	AccessToken    string // OAuth token
	AccessSecret   string // OAuth secret
	MinQueryLength int    // Minimum length of search query

	// Whether to always export book details to scripts
	ExportDetails bool

	// RSS feed/shelves data
	UserID   int64  // User's Goodreads ID
	UserName string // User's Goodreads username (may not be set)

	// Time of last request to Goodreads API.
	// Requests are throttled to 1/sec.
	LastRequest       string
	LastRequestParsed time.Time

	// Scripts
	DefaultScript string `env:"ACTION_DEFAULT"`

	// Workflow data
	BookID     int64
	BookTitle  string `env:"TITLE"`
	AuthorID   int64
	AuthorName string
	ShelfID    int64
	ShelfName  string
	ShelfTitle string
	SeriesID   int64
	SeriesName string `env:"SERIES"`

	// Alternate actions
	FlagAuthor          bool `env:"-"`
	FlagCacheAuthor     bool `env:"-"`
	FlagShelf           bool `env:"-"`
	FlagCacheShelf      bool `env:"-"`
	FlagShelves         bool `env:"-"`
	FlagAddToShelves    bool `env:"-"`
	FlagRemoveFromShelf bool `env:"-"`
	FlagSelectShelf     bool `env:"-"`
	FlagSelectShelves   bool `env:"-"`
	FlagCacheShelves    bool `env:"-"`
	FlagReloadShelf     bool `env:"-"`
	FlagReloadShelves   bool `env:"-"`
	FlagFeeds           bool `env:"-"`
	FlagConf            bool `env:"-"`
	FlagAuthorise       bool `env:"-"`
	FlagDeauthorise     bool `env:"-"`
	FlagOpen            bool `env:"-"`
	FlagScript          bool `env:"-"`
	FlagScripts         bool `env:"-"`
	FlagSearch          bool `env:"-"`
	FlagSeries          bool `env:"-"`
	FlagCacheSeries     bool `env:"-"`
	FlagCacheBook       bool `env:"-"`
	FlagUserInfo        bool `env:"-"`
	FlagHousekeeping    bool `env:"-"`
	FlagIcons           bool `env:"-"`
	FlagHelp            bool `env:"-"`
	FlagNoop            bool `env:"-"`

	// script helper functions
	FlagExport        bool   `env:"-"`
	FlagJSON          bool   `env:"-"`
	FlagBeep          bool   `env:"-"`
	FlagNotify        string `env:"-"`
	FlagNotifyMessage string `env:"-"`
	FlagAction        string `env:"-"`
	FlagHide          bool   `env:"-"`
	FlagPassvars      bool   `env:"-"`
	FlagQuery         string `env:"-"`

	// Search query. Populated from first argument.
	Query string `env:"-"`
	// All args
	Args []string `env:"-"`
}

// QueryEmpty returns true if trimmed query is empty.
func (opts *options) QueryEmpty() bool { return strings.TrimSpace(opts.Query) == "" }

// QueryTooShort returns true if query is empty.
func (opts *options) QueryTooShort() bool {
	return len(strings.TrimSpace(opts.Query)) < opts.MinQueryLength
}

// Authorised returns true if workflow has an OAuth token.
func (opts *options) Authorised() bool { return opts.AccessToken != "" }

func (opts *options) Prepare(args []string) error {
	log.Printf("argv=%#v", args)

	fs = flag.NewFlagSet("alfred-booksearch", flag.ExitOnError)

	fs.BoolVar(&opts.FlagSearch, "search", false, "search for books")
	fs.BoolVar(&opts.FlagConf, "conf", false, "show workflow configuration")

	fs.BoolVar(&opts.FlagAuthor, "author", false, "list books for author")
	fs.BoolVar(&opts.FlagCacheAuthor, "savebooks", false, "cache all books by author")

	fs.BoolVar(&opts.FlagSeries, "series", false, "list books in a series")
	fs.BoolVar(&opts.FlagCacheSeries, "saveseries", false, "cache all books in a series")

	fs.BoolVar(&opts.FlagCacheBook, "savebook", false, "cache book details")

	fs.BoolVar(&opts.FlagShelf, "shelf", false, "list books on shelf")
	fs.BoolVar(&opts.FlagShelves, "shelves", false, "list user shelves")
	fs.BoolVar(&opts.FlagCacheShelf, "saveshelf", false, "cache user shelf")
	fs.BoolVar(&opts.FlagCacheShelves, "saveshelves", false, "cache all user shelves")
	fs.BoolVar(&opts.FlagAddToShelves, "add", false, "add book to shelves")
	fs.BoolVar(&opts.FlagRemoveFromShelf, "remove", false, "remove book from shelf")
	fs.BoolVar(&opts.FlagSelectShelves, "selection", false, "select shelves to add a book to")
	fs.BoolVar(&opts.FlagSelectShelf, "select", false, "toggle shelf selected")
	fs.BoolVar(&opts.FlagReloadShelf, "reload", false, "reload shelf")
	fs.BoolVar(&opts.FlagReloadShelves, "reloadshelves", false, "reload shelves")

	fs.BoolVar(&opts.FlagFeeds, "feeds", false, "fetch RSS feeds")
	fs.BoolVar(&opts.FlagHousekeeping, "housekeeping", false, "check for a new version & clear stale caches")
	fs.BoolVar(&opts.FlagIcons, "icons", false, "download queued icons")
	fs.BoolVar(&opts.FlagAuthorise, "authorise", false, "intiate OAuth authorisation flow")
	fs.BoolVar(&opts.FlagDeauthorise, "deauthorise", false, "delete OAuth credentials")
	fs.BoolVar(&opts.FlagUserInfo, "userinfo", false, "retrieve user info from API")
	fs.BoolVar(&opts.FlagHelp, "h", false, "show this message and exit")
	fs.BoolVar(&opts.FlagOpen, "open", false, "open URL/file")

	fs.BoolVar(&opts.FlagScript, "script", false, "run named script")
	fs.BoolVar(&opts.FlagScripts, "scripts", false, "show scripts")

	fs.BoolVar(&opts.FlagExport, "export", false, "export book details as shell variables")
	fs.BoolVar(&opts.FlagJSON, "json", false, "export book data as JSON")

	fs.BoolVar(&opts.FlagBeep, "beep", false, `play "morse" sound`)
	fs.BoolVar(&opts.FlagNoop, "noop", false, "do nothing")

	fs.StringVar(&opts.FlagNotify, "notify", "", "show notification")
	fs.StringVar(&opts.FlagNotifyMessage, "message", "", "show notification")
	fs.StringVar(&opts.FlagAction, "action", "", "next action")
	fs.BoolVar(&opts.FlagHide, "hide", false, "hide Alfred")
	fs.BoolVar(&opts.FlagPassvars, "passvars", false, "pass variables to next action")
	fs.StringVar(&opts.FlagQuery, "query", "", "search query")

	if err := fs.Parse(args); err != nil {
		return errors.Wrap(err, "parse CLI args")
	}

	logIfError(wf.Config.To(opts), "load configuration: %v")

	opts.Args = fs.Args()
	if len(opts.Args) > 0 {
		opts.Query = strings.TrimSpace(opts.Args[0])
	}
	log.Printf("query=%q", opts.Query)

	// Ensure sensible minimums
	if opts.MaxBooks < minBooksPerAuthor {
		opts.MaxBooks = minBooksPerAuthor
	}
	if opts.MaxCache.Search < minCacheAge {
		opts.MaxCache.Search = minCacheAge
	}
	if opts.MaxCache.Default < minCacheAge {
		opts.MaxCache.Default = minCacheAge
	}
	if opts.MaxCache.Icons < minCacheAge {
		opts.MaxCache.Icons = minCacheAge
	}
	if opts.MaxCache.Shelf < minCacheAge {
		opts.MaxCache.Shelf = minCacheAge
	}
	if opts.MaxCache.Feeds < minCacheAge {
		opts.MaxCache.Feeds = minCacheAge
	}

	if opts.MinQueryLength == 0 {
		opts.MinQueryLength = 2
	}

	if opts.DefaultScript == "" {
		opts.DefaultScript = "View Book Online"
	}

	// Try to read API key from Keychain
	if opts.AccessToken == "" {
		kc := keychain.New(wf.BundleID())
		if s, err := kc.Get(tokensKey); err == nil {
			parts := strings.Split(s, " ")
			opts.AccessToken, opts.AccessSecret = parts[0], parts[1]
		} else if err != keychain.ErrNotFound {
			return errors.Wrap(err, "Keychain")
		}

		wf.Var("ACCESS_TOKEN", opts.AccessToken)
		wf.Var("ACCESS_SECRET", opts.AccessSecret)
	}

	if opts.LastRequest != "" {
		if err := opts.LastRequestParsed.UnmarshalText([]byte(opts.LastRequest)); err != nil {
			return errors.Wrap(err, "parse LastRequest")
		}
	}
	// log.Println("opts=" + spew.Sdump(opts))
	return nil
}
