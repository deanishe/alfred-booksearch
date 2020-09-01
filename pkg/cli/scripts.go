// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-07-29

package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
	"github.com/pkg/errors"
	"go.deanishe.net/alfred-booksearch/pkg/gr"
)

var (
	runner    util.Runner
	imageExts = map[string]struct{}{
		".png":  {},
		".icns": {},
		".gif":  {},
		".jpg":  {},
		".jpeg": {},
	}
)

func init() {
	runner = util.Runners{util.Executable, util.Script}
}

// Show scripts in Alfred.
func runScripts() {
	updateStatus()
	if !authorisedStatus() {
		return
	}

	var scripts []Script
	for _, s := range LoadScripts() {
		scripts = append(scripts, s)
	}
	sort.Sort(Scripts(scripts))

	for _, s := range scripts {
		wf.NewItem(s.Name).
			Arg("-script", s.Name).
			UID(s.Name).
			Copytext(s.Name).
			Valid(true).
			Icon(s.Icon).
			Var("hide_alfred", "true").
			Var("action", "")
	}

	addNavActions()

	if !opts.QueryEmpty() {
		wf.Filter(opts.Query)
	}

	wf.WarnEmpty("No Matching Scripts", "Try a different query?")
	wf.SendFeedback()
}

// Execute specified script.
func runScript() {
	wf.Configure(aw.TextErrors(true))

	scripts := LoadScripts()
	if s, ok := scripts[opts.Query]; ok {
		if opts.ExportDetails { // add all book variables to environment
			b, err := bookDetails(opts.BookID)
			if err != nil {
				notifyError("Fetch Book Details", err)
				log.Fatalf("[ERROR] book details: %v", err)
			}
			for k, v := range bookVariables(b) {
				os.Setenv(k, v)
			}
		}

		log.Printf("[actions] running script %q ...", util.PrettyPath(s.Path))
		data, err := util.RunCmd(runner.Cmd(s.Path))
		if err != nil {
			notifyError(fmt.Sprintf("Run Script %q", s.Name), err)
			log.Fatalf("[ERROR] run script %q: %v", s.Name, err)
		}

		fmt.Print(string(data))
		return
	}
	notifyError("Unknown Script", errors.New(opts.Query))
}

// Script is a built-in or user script.
type Script struct {
	Name string
	Path string
	Icon *aw.Icon
}

// Scripts sorts Scripts by name.
type Scripts []Script

// Implement sort.Interface
func (s Scripts) Len() int           { return len(s) }
func (s Scripts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Scripts) Less(i, j int) bool { return s[i].Name < s[j].Name }

// LoadScripts reads built-in and user scripts.
func LoadScripts() map[string]Script {
	var (
		files = map[string]string{}
		icons = map[string]*aw.Icon{}
	)

	for _, dir := range []string{scriptsDir, userScriptsDir} {
		infos, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Printf("[scripts] [ERROR] read script directory %q: %v", dir, err)
			continue
		}

		for _, fi := range infos {
			var (
				filename = fi.Name()
				path     = filepath.Join(dir, filename)
				ext      = strings.ToLower(filepath.Ext(filename))
				name     = filename[0 : len(filename)-len(ext)]
			)
			if _, ok := imageExts[ext]; ok {
				icons[name] = &aw.Icon{Value: path}
			} else if runner.CanRun(path) {
				files[name] = path
			}
		}
	}

	var (
		scripts = map[string]Script{}
		icon    *aw.Icon
		ok      bool
	)
	for name, path := range files {
		if icon, ok = icons[name]; !ok {
			icon = iconScript
		}
		scripts[name] = Script{name, path, icon}
	}

	for _, s := range scripts {
		log.Printf("[scripts] name=%q, path=%q", s.Name, util.PrettyPath(s.Path))
	}

	return scripts
}

// returns Book populated with all details.
func bookDetails(id int64) (gr.Book, error) {
	key := "books/" + cachefileID(id)
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
		return api.BookDetails(id)
	}

	var b gr.Book
	util.MustExist(filepath.Dir(filepath.Join(wf.CacheDir(), key)))
	if err := wf.Cache.LoadOrStoreJSON(key, opts.MaxCache.Default, reload, &b); err != nil {
		return gr.Book{}, errors.Wrap(err, "book details")
	}

	return b, nil
}

func bookVariables(b gr.Book) map[string]string {
	data := map[string]string{}
	for k, v := range b.Data() {
		data[k] = v
		data[k+"_QUOTED"] = url.PathEscape(v)
		data[k+"_QUOTED_PLUS"] = url.QueryEscape(v)
	}
	return data
}
