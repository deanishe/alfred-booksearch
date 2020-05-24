// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package cli

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"image"
	"image/png"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
	"github.com/disintegration/imaging"
	"github.com/natefinch/atomic"
	"github.com/pkg/errors"
	"go.deanishe.net/alfred-booksearch/pkg/gr"
)

// Workflow icons
var (
	iconBook            = &aw.Icon{Value: "icons/book.png"}
	iconConfig          = &aw.Icon{Value: "icons/config.png"}
	iconDelete          = &aw.Icon{Value: "icons/delete.png"}
	iconDocs            = &aw.Icon{Value: "icons/docs.png"}
	iconError           = &aw.Icon{Value: "icons/error.png"}
	iconHelp            = &aw.Icon{Value: "icons/help.png"}
	iconIssue           = &aw.Icon{Value: "icons/issue.png"}
	iconLocked          = &aw.Icon{Value: "icons/locked.png"}
	iconMore            = &aw.Icon{Value: "icons/more.png"}
	iconOK              = &aw.Icon{Value: "icons/ok.png"}
	iconReload          = &aw.Icon{Value: "icons/reload.png"}
	iconSave            = &aw.Icon{Value: "icons/save.png"}
	iconScript          = &aw.Icon{Value: "icons/script.png"}
	iconShelf           = &aw.Icon{Value: "icons/shelf.png"}
	iconShelfSelected   = &aw.Icon{Value: "icons/shelf-selected.png"}
	iconUpdateAvailable = &aw.Icon{Value: "icons/update-available.png"}
	iconUpdateOK        = &aw.Icon{Value: "icons/update-ok.png"}
	iconWarning         = &aw.Icon{Value: "icons/warning.png"}
	// iconAuthor          = &aw.Icon{Value: "icons/author.png"}
	// iconLink            = &aw.Icon{Value: "icons/link.png"}
	// iconURL             = &aw.Icon{Value: "icons/url.png"}
	// iconDefault         = &aw.Icon{Value: "icon.png"}

	spinnerIcons = []*aw.Icon{
		{Value: "icons/spinner-0.png"},
		{Value: "icons/spinner-1.png"},
		{Value: "icons/spinner-2.png"},
		// {Value: "icons/spinner-3.png"},
	}
)

var (
	userAgent  string
	httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   30 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
		},
	}
)

func init() {
	userAgent = "Alfred Booksearch Workflow " + version + " (+https://github.com/deanishe/alfred-booksearch)"
}

type cacheIcon struct {
	ID   int64
	URL  string
	Path string
}

type iconCache struct {
	Dir       string
	Queue     []cacheIcon
	queueFile string
	seen      map[int64]bool
}

func newIconCache(dir string) *iconCache {
	util.MustExist(dir)
	icons := &iconCache{
		Dir:       dir,
		Queue:     []cacheIcon{},
		queueFile: filepath.Join(dir, "queue.txt"),
		seen:      map[int64]bool{},
	}
	if err := icons.loadQueue(); err != nil {
		panic(err)
	}

	return icons
}

// Add URLs to queue.
func (c *iconCache) Add(books ...gr.Book) {
	for _, b := range books {
		// ignore PNGs, as they're placeholders (real covers are JPG)
		if filepath.Ext(b.ImageURL) == ".png" {
			continue
		}
		if !c.seen[b.ID] {
			if !c.Exists(b) {
				// log.Printf("[icons] queuing for cover retrieval: %s", b)
				c.Queue = append(c.Queue, cacheIcon{
					ID:  b.ID,
					URL: b.ImageURL,
				})
			}
			c.seen[b.ID] = true
		}
	}
}

// BookIcon returns icon for a Book.
func (c *iconCache) BookIcon(b gr.Book) *aw.Icon {
	p := c.cachefile(b.ID)
	if util.PathExists(p) {
		return &aw.Icon{Value: p}
	}
	// Assume that any PNG is a placeholder (actual covers are JPEGs)
	if filepath.Ext(b.ImageURL) == ".png" {
		return iconBook
	}
	// Queue icon for caching
	c.Add(b)
	return iconBook
}

// Exists returns true if book's icon is already cached.
func (c *iconCache) Exists(b gr.Book) bool {
	return util.PathExists(c.cachefile(b.ID))
}

// cachefile returns path of cache file for URL/name s.
func (c *iconCache) cachefile(id int64) string {
	return filepath.Join(c.Dir, cachefileID(id, "png"))
}

// HasQueue returns true if there are Queued files.
func (c *iconCache) HasQueue() bool { return len(c.Queue) > 0 }

func (c *iconCache) loadQueue() error {
	var (
		seen    = map[int64]bool{}
		f       *os.File
		r       *csv.Reader
		records [][]string
		err     error
	)
	if f, err = os.Open(c.queueFile); err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "read icon queue")
		}
		return nil
	}
	defer f.Close()

	c.Queue = []cacheIcon{}
	r = csv.NewReader(f)
	r.Comma = '\t'
	r.FieldsPerRecord = 2
	if records, err = r.ReadAll(); err != nil {
		return errors.Wrap(err, "load queue")
	}
	for _, row := range records {
		id, _ := strconv.ParseInt(row[0], 10, 64)
		if seen[id] {
			continue
		}
		c.Queue = append(c.Queue, cacheIcon{ID: id, URL: row[1]})
		seen[id] = true
	}

	if err = f.Close(); err != nil {
		return errors.Wrap(err, "close queue")
	}
	// clear queue
	if err = atomic.WriteFile(c.queueFile, &bytes.Buffer{}); err != nil {
		return errors.Wrap(err, "clear queue")
	}

	return nil
}

// Close atomically writes queue to disk.
func (c *iconCache) Close() error {
	var (
		buf = &bytes.Buffer{}
		w   = csv.NewWriter(buf)
	)
	w.Comma = '\t'

	for _, icon := range c.Queue {
		if err := w.Write([]string{fmt.Sprintf("%d", icon.ID), icon.URL}); err != nil {
			return errors.Wrapf(err, "write icon %#v", icon)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return errors.Wrap(err, "write TSV")
	}

	if err := atomic.WriteFile(c.queueFile, buf); err != nil {
		return errors.Wrap(err, "write queue file")
	}

	log.Printf("[icons] %d icon(s) queued for download", len(c.Queue))
	c.Queue = []cacheIcon{}
	return nil
}

// ProcessQueue downloads pending icons.
func (c *iconCache) ProcessQueue() error {
	if !c.HasQueue() {
		return nil
	}

	type status struct {
		icon cacheIcon
		err  error
	}

	var (
		wg   sync.WaitGroup
		pool = make(chan struct{}, 5) // Allow 5 parallel downloads
		ch   = make(chan status)
	)
	wg.Add(len(c.Queue))

	for _, icon := range c.Queue {
		icon.Path = c.cachefile(icon.ID)
		go func(icon cacheIcon) {
			defer wg.Done()

			pool <- struct{}{}
			defer func() { <-pool }()

			var (
				img image.Image
				buf = &bytes.Buffer{}
				err error
			)

			if util.PathExists(icon.Path) {
				return
			}

			if img, err = remoteImage(icon.URL); err != nil {
				ch <- status{err: errors.Wrapf(err, "download %q", icon.URL)}
				return
			}
			img = squareImage(img)
			if err = os.MkdirAll(filepath.Dir(icon.Path), 0700); err != nil {
				ch <- status{err: errors.Wrapf(err, "cache directory %q", filepath.Dir(icon.Path))}
				return
			}

			if err = png.Encode(buf, img); err != nil {
				ch <- status{err: errors.Wrapf(err, "convert image %q", icon.URL)}
				return
			}

			if err = atomic.WriteFile(icon.Path, buf); err != nil {
				ch <- status{err: errors.Wrapf(err, "save image %q", icon.URL)}
				return
			}

			ch <- status{icon: icon}
		}(icon)
	}

	go func() {
		wg.Wait()
		close(pool)
		close(ch)
	}()

	var (
		n   int
		err error
	)
	for st := range ch {
		n++
		if st.err != nil {
			logIfError(st.err, "cache icon: %v")
			err = st.err
		} else {
			log.Printf("[icons] [%3d/%d] cached %q to %q", n, len(c.Queue), st.icon.URL, st.icon.Path)
		}
	}
	c.Queue = []cacheIcon{}
	return err
}

func cachefile(key string, ext ...string) string {
	ext = append([]string{filepath.Ext(key)}, ext...)
	s := hash(key)
	path := s[0:2] + "/" + s[2:4] + "/" + s
	return path + strings.Join(ext, "")
}

func cachefileID(id int64, ext ...string) string {
	x := "json"
	if len(ext) > 0 {
		x = ext[0]
	}
	s := fmt.Sprintf("%d", id)
	for len(s) < 4 {
		s = "0" + s
	}
	return fmt.Sprintf("%s/%s/%d.%s", s[0:2], s[2:4], id, x)
}

func remoteImage(URL string) (image.Image, error) {
	var (
		img image.Image
		req *http.Request
		r   *http.Response
		err error
	)

	if req, err = http.NewRequest("GET", URL, nil); err != nil {
		return nil, errors.Wrap(err, "build HTTP request")
	}
	req.Header.Set("User-Agent", userAgent)

	if r, err = httpClient.Do(req); err != nil {
		return nil, errors.Wrap(err, "retrieve URL")
	}
	defer r.Body.Close()
	log.Printf("[%d] %s", r.StatusCode, URL)

	if r.StatusCode > 299 {
		return nil, errors.Wrap(fmt.Errorf("%s: %s", URL, r.Status), "retrieve URL")
	}

	if img, _, err = image.Decode(r.Body); err != nil {
		return nil, errors.Wrap(err, "decode image")
	}
	return img, nil
}

func squareImage(img image.Image) image.Image {
	max := img.Bounds().Max
	n := max.X
	if max.Y > n {
		n = max.Y
	}
	bg := image.NewRGBA(image.Rect(0, 0, n, n))
	return imaging.OverlayCenter(bg, img, 1.0)
}

// spinnerIcon returns a spinner icon. It rotates by 15 deg on every
// subsequent call. Use with wf.Reload(0.1) to implement an animated
// spinner.
func spinnerIcon() *aw.Icon {
	n := wf.Config.GetInt("RELOAD_PROGRESS", 0)
	wf.Var("RELOAD_PROGRESS", fmt.Sprintf("%d", n+1))
	return spinnerIcons[n%3]
}
