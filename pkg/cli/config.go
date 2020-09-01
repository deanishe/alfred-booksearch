// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-07-14

package cli

import (
	"fmt"
	"log"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/keychain"
	"go.deanishe.net/alfred-booksearch/pkg/gr"
)

// initiate OAuth workflow
func runAuthorise() {
	wf.Configure(aw.TextErrors(true))
	// delete existing token (if any)
	logIfError(keychain.New(wf.BundleID()).Delete("oauth_tokens"), "delete existing tokens")
	notifyIfError(api.Authorise(), "Authentication Failed", true)
	notifyIfError(getUserInfo(), "Authentication Failed", true)
	checkErr(notify("OAuth Authentication", "Workflow authorised", "search"))
}

// delete OAuth credentials
func runDeauthorise() {
	wf.Configure(aw.TextErrors(true))
	logIfError(keychain.New(wf.BundleID()).Delete(tokensKey), "delete existing tokens")
	err := wf.Config.Set("USER_ID", "", false).
		Set("USER_NAME", "", false).Do()
	notifyIfError(err, "Deauthorisation Failed", true)
	notify("Workflow Deauthorised", "", "config")
}

// Show workflow configuration.
func runConfig() {
	wf.Var("last_action", "config")
	wf.Var("last_query", opts.Query)

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

	title = "Workflow Authorised"
	subtitle = "Workflow has an OAuth token"
	icon = iconOK
	var arg string
	valid := false
	action := "config"
	hide := ""

	if !opts.Authorised() {
		title = "Workflow Not Authorised"
		subtitle = "↩ to authorise workflow via OAuth"
		icon = iconLocked
		arg = "-authorise"
		action = ""
		valid = true
		hide = "true"
	}

	it := wf.NewItem(title).
		Subtitle(subtitle).
		Icon(icon).
		Arg(arg).
		Valid(valid).
		Var("action", action).
		Var("hide_alfred", hide)

	if opts.Authorised() {
		it.NewModifier(aw.ModCmd).
			Subtitle("Delete OAuth token").
			Icon(iconDelete).
			Arg("-deauthorise").
			Var("action", "config").
			Var("hide_alfred", "")
	}

	wf.NewItem("Open Scripts Folder").
		Subtitle("Open custom scripts folder").
		Arg("-open", userScriptsDir).
		Copytext(userScriptsDir).
		Valid(true).
		Icon(iconScript)

	wf.NewItem("Open Docs").
		Subtitle("Open workflow documentation in your browser").
		Arg("-open", helpURL).
		Copytext(helpURL).
		Valid(true).
		Icon(iconDocs)

	wf.NewItem("Get Help").
		Subtitle("Open workflow issue tracker in your browser").
		Arg("-open", issueTrackerURL).
		Copytext(issueTrackerURL).
		Valid(true).
		Icon(iconHelp)

	wf.NewItem("Report Bug").
		Subtitle("Open workflow issue tracker in your browser").
		Arg("-open", issueTrackerURL).
		Copytext(issueTrackerURL).
		Valid(true).
		Icon(iconIssue)

	addNavActions("config")

	if !opts.QueryEmpty() {
		_ = wf.Filter(opts.Query)
	}
	wf.WarnEmpty("No Matches", "Try a different query?")

	wf.SendFeedback()
}

// fetch user info from API.
func runUserInfo() {
	wf.Configure(aw.TextErrors(true))
	if !opts.Authorised() {
		return
	}

	checkErr(getUserInfo())
}

func getUserInfo() error {
	var (
		user gr.User
		err  error
	)
	if user, err = api.UserInfo(); err != nil {
		return err
	}

	log.Printf("[user] id=%d, name=%s", user.ID, user.Name)

	return wf.Config.
		Set("USER_ID", fmt.Sprintf("%d", user.ID), false).
		Set("USER_NAME", user.Name, false).
		Do()
}
