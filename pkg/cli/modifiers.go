// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	aw "github.com/deanishe/awgo"
)

var validMods = map[string]struct{}{
	"cmd":   {},
	"alt":   {},
	"opt":   {},
	"ctrl":  {},
	"shift": {},
	"fn":    {},
}

// Modifier is a user-specified template.
type Modifier struct {
	Keys   []aw.ModKey
	Script Script
}

// String implements Stringer.
func (m Modifier) String() string {
	// keys := make([]string, len(m.Keys))
	// for i, k := range m.Keys {
	// 	keys[i] = string(k)
	// }
	return fmt.Sprintf("Modifier{Keys: %v, Script: %q}", m.Keys, m.Script.Path)
}

/*
func lookup(data map[string]string) func(string) string {
	return func(key string) string { return data[key] }
}

// For applies Book to template.
func (m Modifier) For(b gr.Book) ([]aw.ModKey, string) {
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
*/

func parseEnv() map[string]string {
	env := map[string]string{}
	for _, s := range os.Environ() {
		i := strings.Index(s, "=")
		if i < 0 || i == len(s)-1 {
			continue
		}
		env[s[0:i]] = s[i+1:]
	}
	return env
}

// LoadModifiers loads user's custom hotkeys.
func LoadModifiers() []Modifier {
	var (
		scripts = LoadScripts()
		mods    []Modifier
		script  Script
		ok      bool
	)
	for k, name := range parseEnv() {
		if k == "ACTION_DEFAULT" || !strings.HasPrefix(k, "ACTION_") {
			continue
		}
		var keys []aw.ModKey
		for _, s := range strings.Split(strings.ToLower(k[7:]), "_") {
			if _, ok := validMods[s]; ok {
				keys = append(keys, aw.ModKey(s))
			}
		}
		if len(keys) == 0 {
			log.Printf("[actions] invalid modifiers: %s", k[7:])
			continue
		}
		if script, ok = scripts[name]; !ok {
			log.Printf("[actions] unknown script: %s", name)
			continue
		}

		log.Printf("[modifiers] %v -> %q", keys, name)
		mods = append(mods, Modifier{Keys: keys, Script: script})
	}
	return mods
}

/*
func LoadModifiers() []Modifier {
	var mods []Modifier
	for _, s := range os.Environ() {
		i := strings.Index(s, "=")
		if i < 0 || i == len(s)-1 {
			continue
		}
		key, value := s[0:i], s[i+1:]
		if !strings.HasPrefix(key, "URL_") {
			continue
		}
		key = key[4:]
		if m, err := newModifier(key, value); err != nil {
			log.Printf("[ERROR] %v", err)
		} else {
			log.Printf("mod=%v", m)
			mods = append(mods, m)
		}
	}
	return mods
}
*/

/*
func newModifier(key, value string) (Modifier, error) {
	m := Modifier{
		Name:  os.Getenv("NAME_" + key),
		Value: value,
	}
	for _, k := range strings.Split(strings.ToLower(key), "_") {
		if _, ok := validMods[k]; ok {
			m.Keys = append(m.Keys, aw.ModKey(k))
		}
	}
	if len(m.Keys) == 0 {
		return Modifier{}, fmt.Errorf("invalid modifiers: %s", key)
	}
	return m, nil
}
*/
