// androidpkg
// Manages Play Store packages using the Play Store API.  Handles text and
// images.
package main

import (
	"flag"
	"fmt"
	"os"
	//"google.golang.org/grpc/credentials"
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
	androidpkg [flags..] command packageName

The commands are:
	info
	  Lookup information about packageName.
	update [flags..]
	  Update packageName.

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
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	textOnly := updateCmd.Bool(
		"textonly", false,
		"Only update text.",
	)
	imagesOnly := updateCmd.Bool(
		"imagesonly", false,
		"Only update images.",
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, USAGE)
		fmt.Fprintf(os.Stderr, "androidpkg flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "update flags:\n")
		updateCmd.PrintDefaults()
	}
	flag.Parse()
	if err := isFile(*credentialsJson); err != nil {
		fatal_usage(fmt.Errorf("credentialsJson got %v", err))
	}

	if flag.NArg() != 2 {
		fatal_usage(fmt.Errorf("wrong number of arguments"))
	}
	do_text := !*imagesOnly
	do_images := !*textOnly

	// Run command.
	var err error
	switch flag.Arg(0) {
	case "info":
		err = info(*credentialsJson, flag.Arg(1))
	case "update":
		if err = isDir(*wordsDir); do_text && err != nil {
			fatal_usage(err)
		}
		if err = isDir(*imagesDir); do_images && err != nil {
			fatal_usage(err)
		}
		err = update(
			*wordsDir, *imagesDir,
			*credentialsJson,
			flag.Arg(1),
			do_text, do_images)
	}
	if err != nil {
		fatal(err)
	}
}

func update(
	wordsDir, imagesDir, credentialsJson, packageName string,
	do_text, do_images bool) error {
	return fmt.Errorf("not implemented")
}

func info(credentialsJson, packageName string) error {
	return fmt.Errorf("not implemented")
}

func fatal_usage(err error) {
	fmt.Fprintf(os.Stderr, "error: %v", err)
	flag.Usage()
	os.Exit(2)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v", err)
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
