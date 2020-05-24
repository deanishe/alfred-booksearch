// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/deanishe/awgo/util"
	"github.com/deanishe/awgo/util/build"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

var (
	info     *build.Info
	env      map[string]string
	ldflags  string
	workDir  string
	buildDir = "./build"
	distDir  = "./dist"
	iconsDir = "./icons"
)

func init() {
	var err error
	if info, err = build.NewInfo(); err != nil {
		panic(err)
	}
	if workDir, err = os.Getwd(); err != nil {
		panic(err)
	}
	env = info.Env()
	env["API_KEY"] = os.Getenv("GOODREADS_API_KEY")
	env["API_SECRET"] = os.Getenv("GOODREADS_API_SECRET")
	env["VERSION"] = info.Version
	env["PKG_CLI"] = "go.deanishe.net/alfred-booksearch/pkg/cli"
	env["PKG_GR"] = "go.deanishe.net/alfred-booksearch/pkg/gr"
	ldflags = `-X "$PKG_GR.version=$VERSION" -X "$PKG_CLI.version=$VERSION" -X "$PKG_CLI.apiKey=$API_KEY" -X "$PKG_CLI.apiSecret=$API_SECRET"`
}

func mod(args ...string) error {
	argv := append([]string{"mod"}, args...)
	return sh.RunWith(env, "go", argv...)
}

// Aliases are mage command aliases.
var Aliases = map[string]interface{}{
	"b": Build,
	"c": Clean,
	"d": Dist,
	"l": Link,
}

// Build builds workflow in ./build
func Build() error {
	mg.Deps(cleanBuild)
	fmt.Println("building ...")

	err := sh.RunWith(env,
		"go", "build", "-ldflags", ldflags,
		"-o", "./build/alfred-booksearch", ".",
	)
	if err != nil {
		return err
	}

	globs := build.Globs(
		"*.png",
		"info.plist",
		"*.html",
		"README.md",
		"LICENCE.txt",
		"icons/*.png",
		"scripts/*",
	)

	if err := build.SymlinkGlobs(buildDir, globs...); err != nil {
		return err
	}

	scriptIcons := []struct {
		script, icon string
	}{
		{"Add to Currently Reading", "shelf"},
		{"Add to Shelves", "shelf"},
		{"Add to Want to Read", "shelf"},
		{"Copy Goodreads Link", "link"},
		{"Mark as Read", "shelf"},
		{"Open Author Page", "author"},
		{"Open Book Page", "link"},
		{"View Author’s Books", "author"},
		{"View Series", "series"},
		{"View Series Online", "link"},
		{"View Similar Books", "link"},
	}

	for _, st := range scriptIcons {
		target := filepath.Join(buildDir, "icons", st.icon+".png")
		link := filepath.Join(buildDir, "scripts", st.script+".png")
		if err := build.Symlink(link, target, true); err != nil {
			return err
		}
	}
	icons := []struct {
		src, dst string
	}{
		{"icons/config.png", "6B3CB906-52D2-4266-8E5F-2F3C1155A05C.png"},
		{"icons/shelf.png", "303CAB58-86FE-497C-995C-11F659969015.png"},
	}

	for _, i := range icons {
		src, dst := filepath.Join(buildDir, i.src), filepath.Join(buildDir, i.dst)
		if err := build.Symlink(dst, src, true); err != nil {
			return err
		}
	}
	return nil
}

// Run run workflow
func Run() error {
	mg.Deps(Build)
	fmt.Println("running ...")
	return sh.RunWith(env, buildDir+"/alfred-booksearch", "-h")
}

// Dist build an .alfredworkflow file in ./dist
func Dist() error {
	mg.SerialDeps(Clean, Build)
	fmt.Printf("exporting %q to %q ...\n", buildDir, distDir)
	p, err := build.Export(buildDir, distDir)
	if err != nil {
		return err
	}

	fmt.Printf("built workflow file %q\n", p)
	return nil
}

// Config display configuration
func Config() {
	fmt.Println("     Workflow name:", info.Name)
	fmt.Println("         Bundle ID:", info.BundleID)
	fmt.Println("  Workflow version:", info.Version)
	fmt.Println("  Preferences file:", info.AlfredPrefsBundle)
	fmt.Println("       Sync folder:", info.AlfredSyncDir)
	fmt.Println("Workflow directory:", info.AlfredWorkflowDir)
	fmt.Println("    Data directory:", info.DataDir)
	fmt.Println("   Cache directory:", info.CacheDir)
}

// Link symlinks ./build directory to Alfred's workflow directory.
func Link() error {
	mg.Deps(Build)

	fmt.Println("linking ./build to workflow directory ...")
	target := filepath.Join(info.AlfredWorkflowDir, info.BundleID)
	// fmt.Printf("target: %s\n", target)

	if util.PathExists(target) {
		fmt.Println("removing existing workflow ...")
	}
	// try to remove it anyway, as dangling symlinks register as existing
	if err := os.RemoveAll(target); err != nil && !os.IsNotExist(err) {
		return err
	}

	src, err := filepath.Abs(buildDir)
	if err != nil {
		return err
	}
	return build.Symlink(target, src, true)
}

// Deps ensure dependencies
func Deps() error {
	mg.Deps(cleanDeps)
	fmt.Println("downloading deps ...")
	return mod("download")
}

// Vendor copy dependencies to ./vendor
func Vendor() error {
	mg.Deps(Deps)
	fmt.Println("vendoring deps ...")
	return mod("vendor")
}

// Clean remove build files
func Clean() {
	fmt.Println("cleaning ...")
	mg.Deps(cleanBuild, cleanMage, cleanDeps)
}

func cleanDeps() error {
	return mod("tidy", "-v")
}

// remove & recreate directory
func cleanDir(name string) error {
	if err := sh.Rm(name); err != nil {
		return err
	}
	return os.MkdirAll(name, 0755)
}

func cleanBuild() error {
	return cleanDir(buildDir)
}

func cleanMage() error {
	return sh.Run("mage", "-clean")
}

// CleanIcons delete all generated icons from ./icons
func CleanIcons() error {
	return cleanDir(iconsDir)
}
