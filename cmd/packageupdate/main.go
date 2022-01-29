// packageupdate.go
// Updates a Google Play Development package.
package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	ap "google.golang.org/api/androidpublisher/v3"

	apt "github.com/napcatstudio/androidpubtools/androidpub"
	"github.com/napcatstudio/translate/xlns"
)

const (
	defaultCredentials = "credentials.json"
	defaultWordsDir    = "words"
	defaultShotsDir    = "images"
	defaultUpdateLimit = -1
	USAGE              = `usage:
packageupdate [flags...] packageName

    Uses the information in credentialsJson to access the Google Play
    Publishing API and update the given package.

where:
    packageName  The name of an APK that the service account has access to.
`
)

type editInfo struct {
	wordsDir, shotsDir, packageName, editId, defBcp47 string
	service                                           *ap.Service
	needsCommit                                       bool
}

func (ei editInfo) String() string {
	return fmt.Sprintf("%s-%s", ei.packageName, ei.defBcp47)
}

func main() {
	credentialsJson := flag.String(
		"credentials", defaultCredentials,
		"Google Play Developer service credentials.",
	)
	wordsDir := flag.String(
		"words", defaultWordsDir,
		"The directory containing the meaning ordered words files.",
	)
	shotsDir := flag.String(
		"images", defaultShotsDir,
		"Images directory.",
	)
	textOnly := flag.Bool(
		"textonly", false,
		"Only update text.",
	)
	imagesOnly := flag.Bool(
		"imagesonly", false,
		"Only update images.",
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, USAGE)
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		usage(fmt.Errorf("wrong number of arguments"))
	}
	ei := editInfo{
		wordsDir:    *wordsDir,
		shotsDir:    *shotsDir,
		packageName: flag.Arg(0),
		needsCommit: false,
	}
	do_text := !*imagesOnly
	do_images := !*textOnly
	fmt.Printf(`credentials: %s
words: %s
images: %s
update text: %v
update images: %v
`,
		*credentialsJson, *wordsDir, *shotsDir, do_text, do_images)

	var err error
	ei.service, err = apt.GetAPService(*credentialsJson)
	if err != nil {
		usage(fmt.Errorf("connecting to %s got %v", *credentialsJson, err))
	}
	ei.editId, err = apt.EditsInsert(ei.service, ei.packageName)
	if err != nil {
		usage(fmt.Errorf("getting edits insert got %v", err))
	}

	// Details
	appDetails, err := ei.service.Edits.Details.Get(
		ei.packageName, ei.editId).Do()
	if err != nil {
		fatal(fmt.Errorf("getting %s details got %v", ei.packageName, err))
	}
	fmt.Printf("%s default lang:%s\ne-mail:%s\nwebsite:%s\n",
		ei.packageName, appDetails.DefaultLanguage,
		appDetails.ContactEmail,
		appDetails.ContactWebsite)
	// Finish setting up info.
	ei.defBcp47 = appDetails.DefaultLanguage

	listings, err := listings(&ei)
	if err != nil {
		fatal(err)
	}

	// Get the BCP-47 codes we need to check.
	translateable, err := xlns.TranslateableGoogleLocales(ei.wordsDir, ei.defBcp47)
	if err != nil {
		fatal(err)
	}
	// We need to update the default locale also.
	translateable = append(translateable, ei.defBcp47)

	// Show the BCP-47 places we aren't checking.
	un, err := xlns.UntranslateableGoogleLocales(ei.wordsDir, ei.defBcp47)
	if err == nil {
		fmt.Printf("Can't update:%v\n", un)
	}

	// By locale.
	for i, listing := range listings {
		// Output BCP-47.
		fmt.Printf("%s (%d/%d)\n", listing.Language, i+1, len(listings))

		can_translate := false
		for _, bcp47 := range translateable {
			if bcp47 == listing.Language {
				can_translate = true
				break
			}
		}
		if !can_translate {
			fmt.Printf("no words for %s\n", listing.Language)
			continue
		}

		if do_text {
			if ei.defBcp47 != listing.Language {
				// The description for the "base" locale is by definition
				// correct (and not in need of translating!).
				err := updateDescriptions(&ei, listing.Language)
				if err != nil {
					fatal(err)
				}
			}
		}
		if do_images {
			err := updateImages(&ei, listing.Language)
			if err != nil {
				fatal(err)
			}
		}
	}

	if ei.needsCommit {
		err := apt.EditsCommit(ei.service, ei.packageName, ei.editId)
		if err != nil {
			fatal(err)
		}
	}
}

func usage(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	flag.Usage()
	os.Exit(2)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(2)
}

// listings returns the listings currently available in the Play Store.
func listings(ei *editInfo) ([]*ap.Listing, error) {
	listings, err := ei.service.Edits.Listings.List(
		ei.packageName, ei.editId).Do()
	if err != nil {
		return nil, fmt.Errorf("getting listings got %v", err)
	}
	var ls []*ap.Listing
	for _, listing := range listings.Listings {
		ls = append(ls, listing)
	}
	return ls, nil
}

// updateDescription updates the description information for a package for
// each BCP-47 location it has information for.
func updateDescriptions(ei *editInfo, bcp47 string) error {
	// Get the base language listing.
	baseListing, err := ei.service.Edits.Listings.Get(
		ei.packageName, ei.editId, ei.defBcp47).Do()
	if err != nil {
		return fmt.Errorf("getting edit listing for %s got %v", bcp47, err)
	}
	base639 := xlns.Iso639FromBcp47(ei.defBcp47)

	// Create translation map.
	iso639 := xlns.Iso639FromBcp47(bcp47)
	xm, err := xlns.WordsXlnsMap(ei.wordsDir, base639, iso639)
	if err != nil {
		return fmt.Errorf("%s %s to %s problem got %v",
			ei.wordsDir, base639, iso639, err)
	}

	// Check if update is needed.
	// Read existing.
	listing, err := ei.service.Edits.Listings.Get(
		ei.packageName, ei.editId, bcp47).Do()
	if err != nil {
		return fmt.Errorf("get listing for %s failed %v", bcp47, err)
	}

	translated := ap.Listing{
		Language:         bcp47,
		Title:            xm.TranslateByLine(baseListing.Title),
		ShortDescription: xm.TranslateByLine(baseListing.ShortDescription),
		FullDescription:  xm.TranslateByLine(baseListing.FullDescription),
	}

	// Compare.
	isTheSame := listing.Title == translated.Title &&
		listing.ShortDescription == translated.ShortDescription &&
		listing.FullDescription == translated.FullDescription
	if isTheSame {
		fmt.Printf("no changes for %s\n", bcp47)
		return nil
	}

	_, err = ei.service.Edits.Listings.Update(
		ei.packageName, ei.editId, bcp47, &translated).Do()
	if err != nil {
		return fmt.Errorf("listing update for %s got %v", bcp47, err)
	}
	ei.needsCommit = true

	return nil
}

type shotInfo struct {
	file, sha1 string
	image      *ap.Image
}

// updateImages checks for image updates.
func updateImages(ei *editInfo, bcp47 string) error {
	iso639 := xlns.Iso639FromBcp47(bcp47)

	// Go through shots.
	for _, imageType := range apt.GooglePlayImageTypes {
		locImageDir := filepath.Join(ei.shotsDir, imageType)
		// Look for locale specific images first.
		pattern := filepath.Join(locImageDir, bcp47+"*.png")
		// Glob only has errors for bad patterns.
		matches, _ := filepath.Glob(pattern)
		if len(matches) == 0 {
			// Look for language specific images.
			pattern = filepath.Join(locImageDir, iso639+"*.png")
			matches, _ = filepath.Glob(pattern)
		}
		// Get info from the directory and from Google.
		sis, err := getShotInfo(matches)
		if err != nil {
			return err
		}
		ilr, err := ei.service.Edits.Images.List(ei.packageName, ei.editId, bcp47, imageType).Do()
		if err != nil {
			return fmt.Errorf("image list for %s %s got %v",
				bcp47, imageType, err)
		}

		// Match up the info.
		var toDelete []*ap.Image
		for _, image := range ilr.Images {
			found := false
			for sii, si := range sis {
				if si.sha1 == image.Sha1 {
					sis[sii].image = image
					found = true
					break
				}
			}
			if !found {
				toDelete = append(toDelete, image)
			}
		}
		// Delete unwanted images.
		for _, doomed := range toDelete {
			fmt.Printf("delete %s %s %s\n", bcp47, imageType, doomed.Id)
			err := ei.service.Edits.Images.Delete(
				ei.packageName, ei.editId, bcp47, imageType, doomed.Id).Do()
			if err != nil {
				return err
			}
			ei.needsCommit = true
		}
		// Upload new images.
		for _, si := range sis {
			if si.image == nil {
				// Update.
				fmt.Printf("upload %s\n", si.file)
				fPng, err := os.Open(si.file)
				if err != nil {
					return fmt.Errorf("can't open %s got %v", si.file, err)
				}
				defer fPng.Close()
				_, err = ei.service.Edits.Images.Upload(
					ei.packageName, ei.editId, bcp47, imageType).Media(fPng).Do()
				if err != nil {
					return fmt.Errorf("uploading %s got %v", si.file, err)
				}
				ei.needsCommit = true
			}
		}
	}
	return nil
}

func getImageInfo(ei *editInfo, bcp47, imageType string) ([]*ap.Image, error) {
	list, err := ei.service.Edits.Images.List(ei.packageName, ei.editId, bcp47, imageType).Do()
	if err != nil {
		return nil, fmt.Errorf("image list for %s %s got %v",
			bcp47, imageType, err)
	}
	return list.Images, nil
}

func getShotInfo(files []string) ([]shotInfo, error) {
	sis := make([]shotInfo, len(files))
	for i, file := range files {
		sha1, err := fileSha1(file)
		if err != nil {
			return nil, fmt.Errorf("SHA1 for %s got %v", file, err)
		}
		sis[i] = shotInfo{file: file, sha1: sha1, image: nil}
	}
	return sis, nil
}

func fileSha1(file string) (string, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha1.Sum(bytes)), nil
}
