// androidpkg
// Manages Play Store packages using the Play Store API.  Handles text and
// images.
package main

import (
	"flag"
	"fmt"
	"os"

	apt "github.com/napcatstudio/androidpubtools/androidpub"
)

const (
	defaultCredentials = "credentials.json"
	defaultWordsDir    = "words"
	defaultImagesDir   = "images"
	USAGE              = `androidpkg is a tool for managing Play Store packages.

It can update the Play Store country text and images.  It uses a meaning
ordered words system for text.  It uses a directory hierarchy for images.
It uses the Google Translate API V3 for translating.

Usage:
	androidpkg [flags..] command packageName [lang..]

The commands are:
	info
	  Lookup information about packageName.
	update
	  Update packageName images and text.
	images
	  Update packageName images using the files in images.
	text
	  Update packageName text using the files in words.
	  
  If one or more lang arguments are provided only check those.

`
)

func main() {
	credentialsJson := flag.String(
		"credentials", defaultCredentials,
		"Google Play Developer service credentials.",
	)
	wordsDir := flag.String(
		"words", defaultWordsDir,
		"The directory containing the meaning ordered words files.",
	)
	imagesDir := flag.String(
		"images", defaultImagesDir,
		"Images directory.",
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, USAGE)
		flag.PrintDefaults()
	}
	flag.Parse()
	if err := isFile(*credentialsJson); err != nil {
		fatal_usage(fmt.Errorf("credentialsJson got %v", err))
	}
	if flag.NArg() < 2 {
		fatal_usage(fmt.Errorf("missing arguments"))
	}
	packageName := flag.Arg(1)
	langs := flag.Args()[2:]

	// Run command.
	var err error
	switch flag.Arg(0) {
	case "info":
		err = apt.PackageInfo(os.Stdout, *credentialsJson, packageName, langs)
	case "images":
		if err = isDir(*imagesDir); err != nil {
			fatal_usage(err)
		}
		err = apt.PackageUpdateImages(
			*credentialsJson, packageName, *imagesDir, langs)
	case "text":
		if err = isDir(*wordsDir); err != nil {
			fatal_usage(err)
		}
		err = apt.PackageUpdateText(*credentialsJson, packageName, *wordsDir, langs)
	case "update":
		if err = isDir(*wordsDir); err != nil {
			fatal_usage(err)
		}
		if err = isDir(*imagesDir); err != nil {
			fatal_usage(err)
		}
		err = apt.PackageUpdate(
			*credentialsJson, packageName, *wordsDir, *imagesDir, langs, true, true)
	}
	if err != nil {
		fatal(err)
	}
}

func fatal_usage(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	flag.Usage()
	os.Exit(2)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(2)
}

func isDir(dir string) error {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("bad path %s", dir)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("%s not directory", dir)
	}
	return nil
}

func isFile(file string) error {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("bad path %s", file)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s not file", file)
	}
	return nil
}
