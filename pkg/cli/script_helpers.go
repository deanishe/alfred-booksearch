// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-08-02

package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	aw "github.com/deanishe/awgo"
	"github.com/keegancsmith/shell"
)

// Play "morse" sound
func runBeep() {
	wf.Configure(aw.TextErrors(true))
	log.Print("[beep] playing sound")
	wf.Alfred.RunTrigger("beep", "")
}

func runVariables() {
	wf.Configure(aw.TextErrors(true))
	v := aw.NewArgVars()

	fs.Visit(func(f *flag.Flag) {
		log.Printf("[variables] %s=%s", f.Name, f.Value)
		switch f.Name {
		case "notify":
			v.Var("notification_title", f.Value.String())

		case "message":
			v.Var("notification_text", f.Value.String())

		case "action":
			v.Var("action", f.Value.String())

		case "hide":
			if !opts.FlagHide {
				v.Var("hide_alfred", "")
			}

		case "passvars":
			var s string
			if opts.FlagPassvars {
				s = "true"
			}
			v.Var("passvars", s)

		case "query":
			v.Var("query", f.Value.String())
		}
	})

	// if hidden, clear all settings
	if opts.FlagHide {
		v.Var("hide_alfred", "true")
		v.Var("action", "")
		v.Var("passvars", "")
		v.Var("query", "")
	}

	checkErr(v.Send())
}

// Export book details as shell variables
func runExport(asJSON bool) {
	wf.Configure(aw.TextErrors(true))

	b, err := bookDetails(opts.BookID)
	if err != nil {
		notifyError("Fetch Book Details", err)
		log.Fatalf("fetch book details: %v", err)
	}

	if asJSON {
		data, err := json.MarshalIndent(b, "", "  ")
		checkErr(err)
		fmt.Print(string(data))
		return
	}

	for k, v := range bookVariables(b) {
		fmt.Println(shell.Sprintf("export %s=%s", k, v))
		// fmt.Printf("export %s=%s\n", k, shell.EscapeArg(v))
	}
}
