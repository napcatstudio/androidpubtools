// package.go
// Contains common functions for dealing with the Android Publishing API V3
// packages.
package androidpub

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	xlns "github.com/napcatstudio/translate"

	ap "google.golang.org/api/androidpublisher/v3"
)

// PackageInfo write package info to the given io.Writer.
func PackageInfo(
	w io.Writer, credentialsJson, packageName string,
	langs []string) error {

	service, err := GetAPService(credentialsJson)
	if err != nil {
		return fmt.Errorf("error %v", err)
	}
	editId, err := EditsInsert(service, packageName)
	if err != nil {
		return fmt.Errorf("error %v", err)
	}

	// Details
	appDetails, err := service.Edits.Details.Get(packageName, editId).Do()
	if err != nil {
		return fmt.Errorf("getting %s details got %v", packageName, err)
	}
	fmt.Fprintf(w, "%s %s\n%s\n%s\n",
		packageName, appDetails.DefaultLanguage,
		appDetails.ContactEmail,
		appDetails.ContactWebsite)
	defLang := appDetails.DefaultLanguage

	// Tracks
	tlr, err := service.Edits.Tracks.List(packageName, editId).Do()
	if err != nil {
		return fmt.Errorf("getting %s tracks got %v", packageName, err)
	}
	var myTracks []string
	fmt.Println("tracks:")
	for _, track := range tlr.Tracks {
		fmt.Fprintf(w, "\t%s\n", track.Track)
		myTracks = append(myTracks, track.Track)
	}

	// Images
	for _, imageType := range GooglePlayImageTypes {
		fmt.Fprintf(w, "imageType: %s\n", imageType)
		ilr, err := service.Edits.Images.List(
			packageName, editId, defLang, imageType).Do()
		if err != nil {
			return fmt.Errorf("getting %s %s images got %v", packageName, defLang, err)
		}
		if len(ilr.Images) == 0 {
			fmt.Fprintf(w, "\tno images\n")
			continue
		}
		for _, image := range ilr.Images {
			fmt.Fprintf(w, "\timage.Id: %s\n", image.Id)
		}
	}

	// Listings
	listings, err := listings(service, packageName, editId, langs)
	if err != nil {
		return fmt.Errorf("getting %s listings got %v", packageName, err)
	}
	for _, listing := range listings {
		if !useListing(langs, listing) {
			continue
		}
		fmt.Fprintf(w, "%s %s\n%s\n%s\n",
			listing.Language, listing.Title,
			listing.ShortDescription,
			listing.FullDescription)
	}

	return nil
}

// PackageUpdate updates a Play Store Android package using the
// AndroidPublisher API V3.
func PackageUpdate(
	credentialsJson, packageName, subFile, wordsDir, imagesDir string,
	langs []string,
	do_text, do_images bool) error {
	//	fmt.Printf(`credentials: %s
	//words: %s
	//images: %s
	//update text: %v
	//update images: %v
	//`,
	//		credentialsJson, wordsDir, imagesDir, do_text, do_images)

	alternates, err := readSubstitutions(subFile)
	if err != nil {
		return err
	}

	service, err := GetAPService(credentialsJson)
	if err != nil {
		return fmt.Errorf("connecting to %s got %v", credentialsJson, err)
	}
	editId, err := EditsInsert(service, packageName)
	if err != nil {
		return fmt.Errorf("getting edits insert got %v", err)
	}

	// Details
	appDetails, err := service.Edits.Details.Get(packageName, editId).Do()
	if err != nil {
		return fmt.Errorf("getting %s details got %v", packageName, err)
	}
	fmt.Printf("%s default lang:%s\ne-mail:%s\nwebsite:%s\n",
		packageName, appDetails.DefaultLanguage,
		appDetails.ContactEmail,
		appDetails.ContactWebsite)
	// Finish setting up info.
	defBcp47 := appDetails.DefaultLanguage

	listings, err := listings(service, packageName, editId, langs)
	if err != nil {
		return err
	}
	if len(listings) == 0 {
		return fmt.Errorf("no listings")
	}
	if len(langs) != 0 && len(listings) != len(langs) {
		return fmt.Errorf("bad language in %v", langs)
	}

	// By locale.
	needsCommit := false
	for i, listing := range listings {
		// Output BCP-47.
		fmt.Printf("%s (%d/%d)\n", listing.Language, i+1, len(listings))

		if do_text {
			if defBcp47 == listing.Language {
				fmt.Printf("default not changing %s\n", defBcp47)
			} else {
				commit, err := updateDescriptions(
					service, editId,
					packageName, wordsDir,
					defBcp47, listing.Language, alternates)
				if err != nil {
					return err
				}
				if commit {
					needsCommit = true
				}
			}
		}

		if do_images {
			commit, err := updateImages(
				service, editId, packageName, imagesDir, listing.Language)
			if err != nil {
				return err
			}
			if commit {
				needsCommit = true
			}
		}
	}

	if needsCommit {
		commit := service.Edits.Commit(packageName, editId)
		//commit.ChangesNotSentForReview(true)
		//err := EditsCommit(service, packageName, editId)
		//editsCommitCall := ap.EditsCommitCall(packageName, editId)
		//editsCommitCall.ChangesNotSentForReview(true)
		_, err := commit.Do()
		if err != nil {
			return err
		}
	}
	return nil
}

// readSubstitutions reads a translation substitution file into a map.  If
// there is not a file it prints a warning and returns an empty map.
func readSubstitutions(subFile string) (map[string]string, error) {
	subs := make(map[string]string)
	f, err := os.Open(subFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: no substitution file '%s'", subFile)
		return subs, nil
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		toks := strings.Split(line, ":")
		if len(toks) == 0 {
			// Blank line.
			continue
		}
		if len(toks) != 2 {
			return nil, fmt.Errorf("bad substitution '%s' in %s", line, subFile)
		}
		// Original to alternate.
		subs[toks[0]] = toks[1]
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading %s got %v", subFile, err)
	}
	return subs, nil
}

// PackageUpdateText updates a Play Store Android package text details using
// the AndroidPublisher API V3.
func PackageUpdateText(
	credentialsJson, packageName, subFile, wordsDir string,
	langs []string) error {

	return PackageUpdate(
		credentialsJson, packageName, subFile, wordsDir, "", langs, true, false)
}

// PackageUpdateText updates a Play Store Android package text details using
// the AndroidPublisher API V3.
func PackageUpdateImages(
	credentialsJson, packageName, imagesDir string,
	langs []string) error {

	return PackageUpdate(
		credentialsJson, packageName, "", "", imagesDir, langs, false, true)
}

// listings returns the listings currently available in the Play Store.
func listings(
	service *ap.Service,
	packageName, editId string,
	langs []string) ([]*ap.Listing, error) {

	listings, err := service.Edits.Listings.List(
		packageName, editId).Do()
	if err != nil {
		return nil, fmt.Errorf("getting listings got %v", err)
	}
	var ls []*ap.Listing
	for _, listing := range listings.Listings {
		if !useListing(langs, listing) {
			continue
		}
		ls = append(ls, listing)
	}
	return ls, nil
}

func useListing(langs []string, listing *ap.Listing) bool {
	if len(langs) == 0 {
		// No restrictions use them all.
		return true
	}
	for _, lang := range langs {
		if lang == listing.Language {
			// This is one is requested.
			return true
		}
	}
	// Not in langs list.
	return false
}

// updateDescription updates the description information for a package for
// each BCP-47 location it has information for.
func updateDescriptions(
	service *ap.Service, editId,
	packageName, wordsDir,
	defBcp47, bcp47 string,
	alternates map[string]string) (bool, error) {

	// Get the base language listing.
	baseListing, err := service.Edits.Listings.Get(
		packageName, editId, defBcp47).Do()
	if err != nil {
		return false, fmt.Errorf("getting edit listing for %s got %v", bcp47, err)
	}

	baseLang, err := langToUse(wordsDir, defBcp47)
	if err != nil {
		return false, err
	}
	lang, err := langToUse(wordsDir, bcp47)
	if err != nil {
		return false, err
	}

	// Create translation map.
	xm, err := xlns.WordsXlnsMap(wordsDir, baseLang, lang)
	if err != nil {
		return false, fmt.Errorf("%s %s to %s problem got %v",
			wordsDir, baseLang, lang, err)
	}

	// Check if update is needed.
	// Read existing.
	listing, err := service.Edits.Listings.Get(
		packageName, editId, bcp47).Do()
	if err != nil {
		return false, fmt.Errorf("get listing for %s failed %v", bcp47, err)
	}

	altTitle := alternates[baseListing.Title]
	translated := ap.Listing{
		Language:         bcp47,
		Title:            xm.TranslateByLineWithAlternate(baseListing.Title, altTitle, 30),
		ShortDescription: xm.TranslateByLine(baseListing.ShortDescription),
		FullDescription:  xm.TranslateByLine(baseListing.FullDescription),
	}

	// Compare.
	isTheSame := listing.Title == translated.Title &&
		listing.ShortDescription == translated.ShortDescription &&
		listing.FullDescription == translated.FullDescription
	if isTheSame {
		fmt.Printf("no listing changes for %s\n", bcp47)
		return false, nil
	}

	_, err = service.Edits.Listings.Update(
		packageName, editId, bcp47, &translated).Do()
	if err != nil {
		return false, fmt.Errorf("listing update for %s got %v", bcp47, err)
	}
	return true, nil
}

func langToUse(wordsDir, bcp47 string) (string, error) {
	// Do we have this as a full BCP-47 language?
	lang, err := xlns.WordsGetLang(wordsDir, bcp47)
	if err != nil {
		return "", err
	}
	return lang, nil
}

type shotInfo struct {
	file, sha1 string
	image      *ap.Image
}

// updateImages checks for image updates.
func updateImages(
	service *ap.Service, editId,
	packageName, imagesDir,
	bcp47 string) (bool, error) {

	iso639 := xlns.Iso639FromBcp47(bcp47)

	needsCommit := false
	// Go through shots.
	for _, imageType := range GooglePlayImageTypes {
		locImageDir := filepath.Join(imagesDir, imageType)
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
		sis, err := getLocalImagesInfo(matches)
		if err != nil {
			return false, err
		}
		ilr, err := service.Edits.Images.List(
			packageName, editId, bcp47, imageType).Do()
		if err != nil {
			return false, fmt.Errorf("image list for %s %s got %v",
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
			err := service.Edits.Images.Delete(
				packageName, editId, bcp47, imageType, doomed.Id).Do()
			if err != nil {
				return false, err
			}
			needsCommit = true
		}
		// Upload new images.
		for _, si := range sis {
			if si.image == nil {
				// Update.
				fmt.Printf("upload %s\n", si.file)
				fPng, err := os.Open(si.file)
				if err != nil {
					return false, fmt.Errorf("can't open %s got %v", si.file, err)
				}
				defer fPng.Close()
				_, err = service.Edits.Images.Upload(
					packageName, editId, bcp47, imageType).Media(fPng).Do()
				if err != nil {
					return false, fmt.Errorf("uploading %s got %v", si.file, err)
				}
				needsCommit = true
			}
		}
	}
	if !needsCommit {
		fmt.Printf("no images changes for %s\n", bcp47)
	}
	return needsCommit, nil
}

func getLocalImagesInfo(files []string) ([]shotInfo, error) {
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
