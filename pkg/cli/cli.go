// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

// Package cli implements the Book Search workflow for Alfred.
package cli

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/keychain"
	"github.com/deanishe/awgo/update"
	"github.com/deanishe/awgo/util"
	"github.com/pkg/errors"

	"go.deanishe.net/alfred-booksearch/pkg/gr"
)

// Configuration values set via LD_FLAGS.
var (
	// Goodreads API key & secret.
	apiKey    = ""
	apiSecret = ""

	// Workflow version.
	version = ""

	navActions []navAction
)

const (
	// workflow links
	repo            = "deanishe/alfred-booksearch"
	helpURL         = "https://github.com/deanishe/alfred-booksearch/tree/master/doc"
	issueTrackerURL = "https://github.com/deanishe/alfred-booksearch/issues"

	// background job names
	booksJob   = "booklist"
	cacheJob   = "housekeeping"
	iconsJob   = "icons"
	shelfJob   = "shelf"
	shelvesJob = "shelves"
	feedsJob   = "feeds"
	userJob    = "user"
	seriesJob  = "series"
	bookJob    = "book"

	tokensKey = "oauth_tokens"

	// how often to re-run workflow if a background job is running
	rerunInterval = 0.2
)

var (
	// cache directories
	authorsCacheDir string
	booksCacheDir   string
	iconCacheDir    string
	searchCacheDir  string
	seriesCacheDir  string
	shelvesCacheDir string

	scriptsDir     string
	userScriptsDir string

	api   *gr.Client
	store *keychainStore
	wf    *aw.Workflow
)

func init() {
	aw.IconError = iconError
	aw.IconWarning = iconWarning
	wf = aw.New(
		aw.HelpURL(issueTrackerURL),
		update.GitHub(repo),
	)
	authorsCacheDir = filepath.Join(wf.CacheDir(), "authors")
	booksCacheDir = filepath.Join(wf.CacheDir(), "books")
	iconCacheDir = filepath.Join(wf.CacheDir(), "covers")
	searchCacheDir = filepath.Join(wf.CacheDir(), "queries")
	shelvesCacheDir = filepath.Join(wf.CacheDir(), "shelves")
	seriesCacheDir = filepath.Join(wf.CacheDir(), "series")

	scriptsDir = "scripts"
	userScriptsDir = filepath.Join(wf.DataDir(), "scripts")

	navActions = []navAction{
		{"Search", "Search for books", "search", iconBook},
		{"Shelves", "List bookshelves", "shelves", iconShelf},
		{"Configuration", "Workflow configuration", "config", iconConfig},
	}
}

// Logger for goodreads library.
type logger struct{}

func (l logger) Printf(format string, args ...interface{}) {
	log.Output(3, fmt.Sprintf(format, args...))
}
func (l logger) Print(args ...interface{}) {
	log.Output(3, fmt.Sprint(args...))
}

var _ gr.Logger = logger{}

type navAction struct {
	title    string
	subtitle string
	action   string
	icon     *aw.Icon
}

// keychainStore implements gr.TokenStore.
type keychainStore struct {
	name          string
	token, secret string
}

// Save saves token & secret to Keychain.
func (s *keychainStore) Save(token, secret string) error {
	kc := keychain.New(s.name)
	if err := kc.Set(tokensKey, token+" "+secret); err != nil {
		return errors.Wrap(err, "save token to Keychain")
	}
	s.token, s.secret = token, secret
	return nil
}

// Load returns OAuth token and secret.
func (s *keychainStore) Load() (token, secret string, err error) {
	return s.token, s.secret, nil
}

var _ gr.TokenStore = (*keychainStore)(nil)

func runHelp() error {
	wf.Configure(aw.TextErrors(true))
	fs.Usage()
	return nil
}

// Run executes the workflow
func Run() {
	wf.Run(run)
}

func run() {
	checkErr(opts.Prepare(wf.Args()))

	if opts.FlagNoop {
		return
	}

	if opts.FlagHelp {
		runHelp()
		return
	}

	if opts.FlagOpen {
		_, err := util.RunCmd(exec.Command("/usr/bin/open", opts.Query))
		notifyIfError(err, "open failed", true)
		return
	}

	checkErr(bootstrap())

	if opts.FlagAuthor {
		runAuthor()
		return
	}

	if opts.FlagCacheAuthor {
		runCacheAuthorList()
		return
	}

	if opts.FlagAuthorise {
		runAuthorise()
		return
	}

	if opts.FlagDeauthorise {
		runDeauthorise()
		return
	}

	if opts.FlagUserInfo {
		runUserInfo()
		return
	}

	if opts.FlagHousekeeping {
		runHousekeeping()
		return
	}

	if opts.FlagFeeds {
		runFeeds()
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

	if opts.FlagShelves {
		runShelves()
		return
	}

	if opts.FlagShelf {
		runShelf()
		return
	}

	if opts.FlagSelectShelves {
		runSelectShelves()
		return
	}

	if opts.FlagCacheShelf {
		runCacheShelf()
		return
	}

	if opts.FlagCacheShelves {
		runCacheShelves()
		return
	}

	if opts.FlagReloadShelf {
		runReloadShelf()
		return
	}

	if opts.FlagReloadShelves {
		runReloadShelves()
		return
	}

	if opts.FlagAddToShelves {
		runAddToShelves()
		return
	}

	if opts.FlagRemoveFromShelf {
		runRemoveFromShelf()
		return
	}

	if opts.FlagSeries {
		runSeries()
		return
	}

	if opts.FlagCacheSeries {
		runCacheSeries()
		return
	}

	if opts.FlagCacheBook {
		runCacheBook()
		return
	}

	if opts.FlagScript {
		runScript()
		return
	}

	if opts.FlagScripts {
		runScripts()
		return
	}

	if opts.FlagExport {
		runExport(opts.FlagJSON)
		return
	}

	if opts.FlagBeep {
		runBeep()
		return
	}

	if opts.FlagSearch {
		runSearch()
		return
	}

	var runVars bool
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "action" || f.Name == "hide" || f.Name == "passvars" || f.Name == "notify" || f.Name == "query" {
			runVars = true
		}
	})
	if runVars {
		runVariables()
		return
	}
}

// Create cache directories & fetch essential data.
func bootstrap() error {
	util.MustExist(authorsCacheDir)
	util.MustExist(booksCacheDir)
	util.MustExist(iconCacheDir)
	util.MustExist(searchCacheDir)
	util.MustExist(seriesCacheDir)
	util.MustExist(shelvesCacheDir)
	util.MustExist(userScriptsDir)

	store = &keychainStore{
		name:   wf.BundleID(),
		token:  opts.AccessToken,
		secret: opts.AccessSecret,
	}

	var err error
	if api, err = gr.New(apiKey, apiSecret, store); err != nil {
		return errors.Wrap(err, "create API client")
	}
	api.Log = logger{}

	if !opts.Authorised() {
		return nil
	}

	// fetch user ID & name if not already set
	if opts.UserID == 0 {
		if err := runJob(userJob, "-userinfo"); err != nil {
			return err
		}
	} else if wf.Cache.Exists(shelvesKey) {
		if !wf.IsRunning(feedsJob) {
			var t time.Time
			if wf.Cache.Exists(feedsKey) {
				if err := wf.Cache.LoadJSON(feedsKey, &t); err != nil {
					return err
				}
			}

			if time.Since(t) > opts.MaxCache.Feeds {
				if err := runJob(feedsJob, "-feeds"); err != nil {
					return err
				}
			}
		}
	} else if err := runJob(shelvesJob, "-saveshelves"); err != nil {
		return err
	}

	return nil
}

// Show "update available" message and check for update if due.
func updateStatus() {
	if wf.UpdateAvailable() && opts.Query == "" {
		wf.NewItem("Update Available!").
			Subtitle("↩ or ⇥ to install update").
			Valid(false).
			Autocomplete("workflow:update").
			Icon(iconUpdateAvailable)
	}

	if wf.UpdateCheckDue() && !wf.IsRunning(cacheJob) {
		logIfError(runJob(cacheJob, "-housekeeping"), "check for update: %v")
	}
}

// Show "authorise workflow" action if workflow has no OAuth token.
func authorisedStatus() bool {
	if !opts.Authorised() {
		wf.NewItem("Authorise Workflow").
			Subtitle("Action this item to authorise workflow to access your Goodreads account").
			Arg("-authorise").
			Valid(true).
			Icon(iconLocked).
			Var("hide_alfred", "true")

		wf.SendFeedback()
		return false
	}
	_, err := api.AuthedClient()
	checkErr(err)
	return true
}

func addNavActions(ignore ...string) {
	if len(opts.Query) < 3 {
		return
	}

	ig := make(map[string]bool, len(ignore))
	for _, s := range ignore {
		ig[s] = true
	}

	for _, a := range navActions {
		if ig[a.action] || !strings.HasPrefix(a.action, strings.ToLower(opts.Query)) {
			continue
		}
		wf.NewItem(a.title).
			Subtitle(a.subtitle).
			Arg("-noop").
			UID(a.action).
			Valid(true).
			Icon(a.icon).
			Var("action", a.action).
			Var("last_action", "").
			Var("last_query", "").
			Var("query", "").
			Var("hide_alfred", "")
	}
}

func hash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func notify(title, msg string, action ...string) error {
	v := aw.NewArgVars().
		Var("notification_title", title).
		Var("notification_text", msg).
		Var("hide_alfred", "true")

	if len(action) > 0 {
		v.Var("action", action[0])
	}

	return v.Send()
}

func notifyError(title string, err error, command ...string) error {
	return notify("💀 "+title+" 💀", err.Error(), command...)
}

func notifyIfError(err error, title string, fatal bool) {
	if err == nil {
		return
	}
	_ = notify("💀 "+title+" 💀", err.Error())
	log.Fatalf("[ERROR] %s: %v", title, err)
}

func logIfError(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	args = append(args, err)
	log.Printf("[ERROR] "+format, args...)
}

func checkErr(err error) {
	if err == nil {
		return
	}
	panic(err)
}

// start a named background job, passing the given arguments to this executable.
func runJob(name string, args ...string) error {
	if wf.IsRunning(name) {
		return nil
	}
	return wf.RunInBackground(name, exec.Command(os.Args[0], args...))
}
