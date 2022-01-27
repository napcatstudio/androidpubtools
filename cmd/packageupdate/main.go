// packageupdate.go
// Updates a Google Play Development package.
package main

import (
	//"crypto/sha1"
	"flag"
	"fmt"

	//"io/ioutil"
	//"log"
	"os"

	//"path/filepath"

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
	//updateLimit := flag.Int(
	//	"limit", defaultUpdateLimit,
	//	"Limit update to this number of locales (-1 = all).",
	//)
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
	packageName := flag.Arg(0)
	do_text := !*imagesOnly
	do_images := !*textOnly
	fmt.Printf(`credentials: %s
words: %s
images: %s
update text: %v
update images: %v
`,
		*credentialsJson, *wordsDir, *shotsDir, do_text, do_images)

	service, err := apt.GetAPService(*credentialsJson)
	if err != nil {
		usage(fmt.Errorf("connecting to %s got %v", *credentialsJson, err))
	}
	editId, err := apt.EditsInsert(service, packageName)
	if err != nil {
		usage(fmt.Errorf("getting edits insert got %v", err))
	}
	needsCommit := false

	// Details
	appDetails, err := service.Edits.Details.Get(packageName, editId).Do()
	if err != nil {
		fatal(fmt.Errorf("getting %s details got %v", packageName, err))
	}
	// TODO: limit this.
	fmt.Printf("%s %s\n%s\n%s\n",
		packageName, appDetails.DefaultLanguage,
		appDetails.ContactEmail,
		appDetails.ContactWebsite)
	var defBcp47 = appDetails.DefaultLanguage

	// Get the BCP-47 codes we need to check.
	translateable, err := xlns.TranslateableGoogleLocales(*wordsDir, defBcp47)
	if err != nil {
		fatal(err)
	}
	// We need to update the default locale also.
	translateable = append(translateable, defBcp47)

	// Set update limit.
	//if *updateLimit == -1 {
	//	*updateLimit = len(translateable)
	//}

	// Show the BCP-47 places we aren't checking.
	un, err := xlns.UntranslateableGoogleLocales(*wordsDir, defBcp47)
	if err == nil {
		fmt.Printf("Not updating:%v\n", un)
	}

	/*


		// How many calls to Google?
		//log.Printf("%d calls", *xlns.APCalls)
		lastCalls := *xlns.APCalls

		// By locale.
		needsCommit := false
		needingUpdate := 0
		for i, bcp47 := range translateable {
			// Output BCP-47.
			log.Printf("%s (%d/%d)", bcp47, i+1, len(translateable))

			localeNeedsUpdate := false
			if text {
				if defBcp47 != bcp47 {
					// The description for the "base" locale is by definition
					// correct (and not in need of translating!).
					localeNeedsUpdate = updateDescriptions(*wordsDir, apei, defBcp47, bcp47)
				}
			}
			if images {
				localeNeedsUpdate =
					localeNeedsUpdate ||
						updateImages(*shotsDir, apei, defBcp47, bcp47) || localeNeedsUpdate
			}

			needsCommit = needsCommit || localeNeedsUpdate

			//log.Printf("%d calls +%d", *xlns.APCalls, *xlns.APCalls-lastCalls)
			lastCalls = *xlns.APCalls

			if localeNeedsUpdate {
				needingUpdate++
				if needingUpdate >= *updateLimit {
					break
				}
			}
		}

	*/
	if needsCommit {
		err := apt.EditsCommit(service, packageName, editId)
		if err != nil {
			fatal(err)
		}
	}

	//log.Printf("%d calls +%d", *xlns.APCalls, *xlns.APCalls-lastCalls)
	fmt.Println("it worked?")
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

//var showBase = true

/*
// updateDescription updates the description information for a package for
// each BCP-47 location it has information for.
func updateDescriptions(
	wordsDir string,
	apei *xlns.APEditInfo,
	defBcp47, bcp47 string) bool {
	// Get the base language listing.
	baseLangGet := apei.Srv.Edits.Listings.Get(apei.Package, apei.Id, defBcp47)
	var baseListing *ap.Listing
	err := xlns.ExpBackoff(func() error {
		var err error
		baseListing, err = baseLangGet.Do()
		return err
	})
	if err != nil {
		log.Fatalf("Edits.Listings.Get failed %v", err)
	}
	if showBase {
		log.Printf("\n%s\n%s\n%s\n",
			baseListing.Title,
			baseListing.ShortDescription,
			baseListing.FullDescription)
		showBase = false
	}
	base639 := xlns.Iso639FromBcp47(defBcp47)

	// Create translation map.
	iso639 := xlns.Iso639FromBcp47(bcp47)
	xm, err := xlns.WordsXlnsMap(wordsDir, base639, iso639)
	if err != nil {
		log.Fatalf("%s %s to %s problem got %v",
			wordsDir, base639, iso639, err)
	}

	// Check if update is needed.
	// Read existing.
	// TODO: on failure add...
	langGet := apei.Srv.Edits.Listings.Get(apei.Package, apei.Id, bcp47)
	var listing *ap.Listing
	err = xlns.ExpBackoff(func() error {
		var err error
		listing, err = langGet.Do()
		return err
	})
	if err != nil {
		log.Fatalf("Edits.Listings.Get %s failed %v", bcp47, err)
	}

	translated := ap.Listing{
		Language:         bcp47,
		Title:            xm.TranslateByLine(baseListing.Title),
		ShortDescription: xm.TranslateByLine(baseListing.ShortDescription),
		FullDescription:  xm.TranslateByLine(baseListing.FullDescription),
	}
	// TODO: Check for bad translation.

	// Compare.
	isTheSame := listing.Title == translated.Title &&
		listing.ShortDescription == translated.ShortDescription &&
		listing.FullDescription == translated.FullDescription
	if isTheSame {
		log.Printf("no changes for %s", bcp47)
		return false
	}

	update := apei.Srv.Edits.Listings.Update(
		apei.Package, apei.Id, bcp47, &translated)
	err = xlns.ExpBackoff(func() error {
		_, err := update.Do() // *Listing, error
		return err
	})
	if err != nil {
		log.Fatalf("Edits.Listings.Update failed %v", err)
	}

	return true
}

type shotInfo struct {
	file, sha1 string
	image      *ap.Image
}

func fileSha1(file string) (string, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha1.Sum(bytes)), nil
}

func getShotInfo(files []string) []shotInfo {
	sis := make([]shotInfo, len(files))
	for i, file := range files {
		sha1, err := fileSha1(file)
		if err != nil {
			log.Fatalf("SHA1 for %s failed %v", file, err)
		}
		sis[i] = shotInfo{file: file, sha1: sha1, image: nil}
	}
	return sis
}

func getImageInfo(
	apei *xlns.APEditInfo,
	bcp47, imageType string) []*ap.Image {
	list := apei.Srv.Edits.Images.List(apei.Package, apei.Id, bcp47, imageType)
	var ilr *ap.ImagesListResponse
	err := xlns.ExpBackoff(func() error {
		var err error
		ilr, err = list.Do()
		return err
	})
	if err != nil {
		log.Fatalf("Edits.Images.List for %s %s failed %v",
			bcp47, imageType, err)
	}
	return ilr.Images
}

// updateImages checks for image updates.
func updateImages(
	shotsDir string,
	apei *xlns.APEditInfo,
	defBcp47, bcp47 string) bool {
	iso639 := xlns.Iso639FromBcp47(bcp47)
	needsCommit := false

	// Go through shots.
	for _, imageType := range xlns.GoogleImageTypes {
		locImageDir := filepath.Join(shotsDir, imageType)
		// Look for locale specific images first.
		pattern := filepath.Join(locImageDir, bcp47+"*.png")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			log.Fatalf("Glob for %s failed %v", pattern, err)
		}
		if len(matches) == 0 {
			// Look for language specific images.
			pattern = filepath.Join(locImageDir, iso639+"*.png")
			matches, err = filepath.Glob(pattern)
			if err != nil {
				log.Fatalf("Glob for %s failed %v", pattern, err)
			}
		}
		// Get info from the directory and from Google.
		sis := getShotInfo(matches)
		images := getImageInfo(apei, bcp47, imageType)
		// Match up the info.
		var toDelete []*ap.Image
		for _, image := range images {
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
			log.Printf("delete %s %s %s\n", bcp47, imageType, doomed.Id)
			del := apei.Srv.Edits.Images.Delete(
				apei.Package, apei.Id, bcp47, imageType, doomed.Id)
			err := xlns.ExpBackoff(func() error {
				return del.Do()
			})
			if err != nil {
				log.Fatalf("Edits.Images.Delete %s %s %s failed %v",
					bcp47, imageType, doomed.Id, err)
			}
			needsCommit = true
		}
		// Upload new images.
		for _, si := range sis {
			if si.image == nil {
				// Update.
				log.Printf("upload %s\n", si.file)
				upload := apei.Srv.Edits.Images.Upload(
					apei.Package, apei.Id, bcp47, imageType)
				fPng, err := os.Open(si.file)
				if err != nil {
					log.Fatalf("Open %s %v", si.file, err)
				}
				defer fPng.Close()
				_ = upload.Media(fPng)
				err = xlns.ExpBackoff(func() error {
					_, err := upload.Do()
					return err
				})
				if err != nil {
					log.Fatalf("Edits.Images.Upload failed %v", err)
				}
				needsCommit = true
			}
		}
	}
	return needsCommit
}
*/
