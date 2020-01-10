// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	aw "github.com/deanishe/awgo"
)

// Modifier is a user-specified template.
type Modifier struct {
	Name  string // optional
	Keys  []aw.ModKey
	Value string
}

func lookup(data map[string]string) func(string) string {
	return func(key string) string { return data[key] }
}

// For applies Book to template.
func (m Modifier) For(b Book) ([]aw.ModKey, string) {
	data := map[string]string{}
	for k, v := range b.Data() {
		data[k] = url.QueryEscape(v)
		data[k+"Alt"] = url.PathEscape(v)
		data[k+"Raw"] = v
	}
	return m.Keys, os.Expand(m.Value, lookup(data))
}

// String formats Modifier for printing.
func (m Modifier) String() string {
	return fmt.Sprintf("Modifier{Keys: %+v, Value: %q}", m.Keys, m.Value)
}

func LoadModifiers() []Modifier {
	var mods []Modifier
	for _, s := range os.Environ() {
		i := strings.Index(s, "=")
		if i < 0 || i == len(s)-1 {
			continue
		}
		key, value := s[0:i], s[i+1:len(s)]
		if !strings.HasPrefix(key, "URL_") {
			continue
		}
		key = key[4:len(key)]
		if m, err := newModifier(key, value); err != nil {
			log.Printf("[ERROR] %v", err)
		} else {
			log.Printf("mod=%v", m)
			mods = append(mods, m)
		}
	}
	return mods
}

var validMods = map[string]bool{
	"cmd":   true,
	"alt":   true,
	"opt":   true,
	"ctrl":  true,
	"shift": true,
	"fn":    true,
}

func newModifier(key, value string) (Modifier, error) {
	m := Modifier{
		Name:  os.Getenv("NAME_" + key),
		Value: value,
	}
	for _, k := range strings.Split(strings.ToLower(key), "_") {
		if validMods[k] {
			m.Keys = append(m.Keys, aw.ModKey(k))
		}
	}
	if len(m.Keys) == 0 {
		return Modifier{}, fmt.Errorf("invalid modifiers: %s", key)
	}
	return m, nil
}
