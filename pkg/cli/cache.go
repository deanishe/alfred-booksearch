// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-07-18

package cli

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
	"github.com/pkg/errors"

	"go.deanishe.net/alfred-booksearch/pkg/gr"
)

const feedsKey = "FeedsLastUpdate.json"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func runIcons() {
	wf.Configure(aw.TextErrors(true))
	icons := newIconCache(iconCacheDir)
	if icons.HasQueue() {
		checkErr(icons.ProcessQueue())
	}
}

// Fetch RSS feeds and cache icons.
func runFeeds() {
	wf.Configure(aw.TextErrors(true))
	// fetch RSS feeds
	if !wf.Cache.Exists(shelvesKey) {
		log.Printf("[feeds] no shelves")
		return
	}

	if opts.UserID == 0 {
		log.Printf("[feeds] user ID not set; not fetching RSS feeds")
		return
	}

	var (
		icons   = newIconCache(iconCacheDir)
		shelves []gr.Shelf
		err     error
	)

	checkErr(wf.Cache.LoadJSON(shelvesKey, &shelves))
	checkErr(wf.Cache.StoreJSON(feedsKey, time.Now()))

	log.Println("[feeds] fetching RSS feeds...")
	for _, s := range shelves {
		var feed gr.Feed
		if feed, err = api.FetchFeed(opts.UserID, s.Name); err == nil {
			log.Printf("[feeds] %d book(s) in feed %q", len(feed.Books), s.Name)
			icons.Add(feed.Books...)
		}
		logIfError(err, "fetch feed %q: %v", s.Name)
	}

	if icons.HasQueue() {
		log.Printf("[feeds] %d icon(s) queued for download", len(icons.Queue))
		if err = icons.Close(); err == nil {
			err = runJob(iconsJob, "-icons")
		}
		checkErr(err)
	}
}

// Check for workflow update + clear stale cache files.
func runHousekeeping() {
	wf.Configure(aw.TextErrors(true))

	// wait a bit for current search to complete before clearing caches
	// time.Sleep(15)

	wg := sync.WaitGroup{}
	wg.Add(3)

	// check for update
	go func() {
		defer wg.Done()
		log.Println("[housekeeping] checking for updates...")
		logIfError(wf.CheckForUpdate(), "[housekeeping] update check: %v")
	}()

	// clean covers cache
	go func() {
		defer wg.Done()
		log.Println("[housekeeping] cleaning cover cache...")
		dc := &dirCleaner{
			root: iconCacheDir,
			maxAge: func() time.Duration {
				// fuzzy age of max cache age +/- 72 hours
				delta := time.Hour * time.Duration(rand.Int31n(72))
				return opts.MaxCache.Icons - (time.Hour * 72) + delta
			},
		}
		logIfError(dc.Clean(), "[housekeeping] clean icon cache: %v")
	}()

	// clean other caches
	go func() {
		defer wg.Done()
		dirs := []string{authorsCacheDir, booksCacheDir, searchCacheDir}
		ages := []time.Duration{opts.MaxCache.Default, opts.MaxCache.Default, opts.MaxCache.Search}
		for i, dir := range dirs {
			i := i
			log.Printf("[housekeeping] cleaning %s cache...", filepath.Base(dir))
			dc := &dirCleaner{
				root:   dir,
				maxAge: func() time.Duration { return ages[i] },
			}
			logIfError(dc.Clean(), "[housekeeping] clean query cache: %v")
		}
	}()

	wg.Wait()
}

type cacheDir struct {
	info os.FileInfo
	path string
}

// cacheDirs sorts directories by name.
type cacheDirs []cacheDir

// Implement sort.Interface
func (s cacheDirs) Len() int           { return len(s) }
func (s cacheDirs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s cacheDirs) Less(i, j int) bool { return s[i].path < s[j].path }

// removes old files and empty directories from a cache directory.
type dirCleaner struct {
	root   string
	maxAge func() time.Duration
	dirs   []cacheDir
}

func (dc *dirCleaner) addDir(fi os.FileInfo, path string) {
	dc.dirs = append(dc.dirs, cacheDir{fi, path})
}

func (dc *dirCleaner) Clean() error {
	if err := dc.cleanFiles(); err != nil {
		return err
	}
	return dc.cleanDirs()
}

func (dc *dirCleaner) cleanDirs() error {
	sort.Sort(sort.Reverse(cacheDirs(dc.dirs)))
	for _, dir := range dc.dirs {
		if dir.path == dc.root {
			continue
		}

		if time.Since(dir.info.ModTime()) < time.Hour*72 {
			continue
		}

		infos, err := ioutil.ReadDir(dir.path)
		if err != nil {
			return err
		}
		if len(infos) == 0 {
			log.Printf("[housekeeping] deleting empty directory %q ...", util.PrettyPath(dir.path))
			if err := os.Remove(dir.path); err != nil {
				return errors.Wrap(err, util.PrettyPath(dir.path))
			}
		}
	}
	return nil
}

func (dc *dirCleaner) cleanFiles() error {
	clean := func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			dc.addDir(fi, p)
			return nil
		}
		// delete cached queries (.json) and covers (.png).
		x := filepath.Ext(fi.Name())
		if x != ".json" && x != ".png" {
			return nil
		}

		age := time.Since(fi.ModTime())
		if age > dc.maxAge() {
			log.Printf("[housekeeping] deleting %q (%v) ...", util.PrettyPath(p), age)
			if err := os.Remove(p); err != nil {
				return errors.Wrap(err, util.PrettyPath(p))
			}
		}
		return nil
	}

	return filepath.Walk(dc.root, clean)
}
