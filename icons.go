// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
	"github.com/disintegration/imaging"
	"github.com/natefinch/atomic"
	"github.com/pkg/errors"
)

type iconCache struct {
	Dir       string
	Queue     []string
	queueFile string
	seen      map[string]bool
}

func newIconCache(dir string) *iconCache {
	util.MustExist(dir)
	qFile := filepath.Join(dir, "queue.txt")
	icons := &iconCache{
		Dir:       dir,
		Queue:     []string{},
		queueFile: qFile,
		seen:      map[string]bool{},
	}
	if err := icons.loadQueue(); err != nil {
		panic(err)
	}

	return icons
}

// Add URLs to queue.
func (c *iconCache) Add(URL ...string) {
	for _, u := range URL {
		if !c.seen[u] {
			c.Queue = append(c.Queue, u)
			c.seen[u] = true
		}
	}
}

// BookIcon returns icon for a Book.
func (c *iconCache) BookIcon(b Book) *aw.Icon {
	// Assume that any PNG is a placeholder (actual covers are JPEGs)
	if filepath.Ext(b.ImageURL) == ".png" {
		return iconBook
	}
	p := c.cachefile(b.ImageURL)
	if util.PathExists(p) {
		return &aw.Icon{Value: p}
	}
	c.Add(b.ImageURL)
	return iconBook
}

// cachefile returns path of cache file for URL/name s.
func (c *iconCache) cachefile(s string) string {
	return filepath.Join(c.Dir, cachefile(s, ".png"))
}

// HasQueue returns true if there are Queued files.
func (c *iconCache) HasQueue() bool { return len(c.Queue) > 0 }

func (c *iconCache) loadQueue() error {
	var (
		f   *os.File
		scn *bufio.Scanner
		err error
	)
	if f, err = os.Open(c.queueFile); err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "read icon queue")
		}
		return nil
	}
	defer f.Close()

	c.Queue = []string{}
	scn = bufio.NewScanner(f)
	for scn.Scan() {
		s := scn.Text()
		log.Printf("[icons] queue: %q", s)
		if s != "" {
			c.Queue = append(c.Queue, s)
		}
	}
	if err = scn.Err(); err != nil {
		return errors.Wrap(err, "load queue")
	} else {
		if err = f.Close(); err != nil {
			return errors.Wrap(err, "close queue")
		}
		// clear queue
		if err = atomic.WriteFile(c.queueFile, &bytes.Buffer{}); err != nil {
			return errors.Wrap(err, "clear queue")
		}
	}
	return nil
}

// Close atomically writes queue to disk.
func (c *iconCache) Close() error {
	s := strings.Join(c.Queue, "\n")
	if err := ioutil.WriteFile(c.queueFile, []byte(s), 0600); err != nil {
		return errors.Wrap(err, "save queue")
	}
	log.Printf("queued %d icon(s) for download", len(c.Queue))
	c.Queue = []string{}
	return nil
}

// ProcessQueue downloads pending icons.
func (c *iconCache) ProcessQueue() error {
	if !c.HasQueue() {
		return nil
	}

	var wg sync.WaitGroup
	// Allow 3 parallel downloads
	pool := make(chan bool, 3)
	errs := make(chan error)

	wg.Add(len(c.Queue))
	for _, URL := range c.Queue {
		log.Printf("image URL: %s", URL)
		go func(u string) {
			defer wg.Done()

			pool <- true
			defer func() {
				_ = <-pool
			}()

			var (
				img image.Image
				f   *os.File
				err error
			)

			if img, err = remoteImage(u); err != nil {
				log.Printf("[ERROR] download %q: %v", u, err)
				errs <- err
				return
			}
			img = squareImage(img)
			p := c.cachefile(u)
			if err = os.MkdirAll(filepath.Dir(p), 0700); err != nil {
				log.Printf("[ERROR] cache directory %q: %v", filepath.Dir(p), err)
				errs <- err
				return
			}
			if f, err = os.Create(p); err != nil {
				log.Printf("[ERROR] cache directory %q: %v", filepath.Dir(p), err)
				errs <- err
				return
			}
			defer f.Close()
			if err = png.Encode(f, img); err != nil {
				log.Printf("[ERROR] save image %q: %v", u, err)
				errs <- err
				return
			}

			log.Printf("cached %q to %q", u, util.PrettyPath(p))
		}(URL)
	}

	go func() {
		wg.Wait()
		close(pool)
		close(errs)
	}()

	var err error
	for e := range errs {
		err = e
	}
	c.Queue = []string{}
	return err
}

func cachefile(u string, ext ...string) string {
	ext = append([]string{filepath.Ext(u)}, ext...)
	s := hash(u)
	s = s[0:2] + "/" + s[2:4] + "/" + s[4:6] + "/" + s[6:len(s)-1]
	return s + strings.Join(ext, "")
}

/*
func download(URL, path string) (err error) {
	var (
		f *os.File
		r *http.Response
	)

	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return
	}

	if r, err = client.Get(URL); err != nil {
		return errors.Wrap(err, URL)
	}
	defer r.Body.Close()
	if r.StatusCode > 299 {
		return errors.New(r.Status)
	}

	if f, err = os.Create(path); err != nil {
		return
	}
	defer f.Close()

	if _, err = io.Copy(f, r.Body); err != nil {
		return
	}

	return nil
}
*/

func remoteImage(URL string) (image.Image, error) {
	var (
		img image.Image
		r   *http.Response
		err error
	)

	if r, err = client.Get(URL); err != nil {
		return nil, errors.Wrap(err, URL)
	}
	defer r.Body.Close()
	if r.StatusCode > 299 {
		return nil, errors.New(r.Status)
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
	var (
		step    = 15
		max     = (45 / step) - 1
		current = wf.Config.GetInt("RELOAD_PROGRESS", 0)
		next    = current + 1
	)
	if next > max {
		next = 0
	}

	log.Printf("progress: current=%d, next=%d", current, next)

	wf.Var("RELOAD_PROGRESS", fmt.Sprintf("%d", next))

	if current == 0 {
		return iconSpinner
	}

	return &aw.Icon{Value: fmt.Sprintf("icons/spinner-%d.png", current*step)}
}
