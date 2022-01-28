// packageupdate.go
// Updates a Google Play Development package.
package main

import (
	//"crypto/sha1"
	"flag"
	"fmt"
	"log"

	//"io/ioutil"
	//"log"
	"os"

	//"path/filepath"
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
	info := editInfo{
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
	info.service, err = apt.GetAPService(*credentialsJson)
	if err != nil {
		usage(fmt.Errorf("connecting to %s got %v", *credentialsJson, err))
	}
	info.editId, err = apt.EditsInsert(info.service, info.packageName)
	if err != nil {
		usage(fmt.Errorf("getting edits insert got %v", err))
	}

	// Details
	appDetails, err := info.service.Edits.Details.Get(
		info.packageName, info.editId).Do()
	if err != nil {
		fatal(fmt.Errorf("getting %s details got %v", info.packageName, err))
	}
	// TODO: limit this.
	fmt.Printf("%s %s\n%s\n%s\n",
		info.packageName, appDetails.DefaultLanguage,
		appDetails.ContactEmail,
		appDetails.ContactWebsite)
	// Finish setting up info.
	info.defBcp47 = appDetails.DefaultLanguage

	listings, err := listings(&info)
	if err != nil {
		fatal(err)
	}

	// Get the BCP-47 codes we need to check.
	translateable, err := xlns.TranslateableGoogleLocales(info.wordsDir, info.defBcp47)
	if err != nil {
		fatal(err)
	}
	// We need to update the default locale also.
	translateable = append(translateable, info.defBcp47)

	// Set update limit.
	//if *updateLimit == -1 {
	//	*updateLimit = len(translateable)
	//}

	// Show the BCP-47 places we aren't checking.
	un, err := xlns.UntranslateableGoogleLocales(info.wordsDir, info.defBcp47)
	if err == nil {
		fmt.Printf("Can't update:%v\n", un)
	}

	/*
		// How many calls to Google?
		//log.Printf("%d calls", *xlns.APCalls)
		lastCalls := *xlns.APCalls
	*/

	// By locale.
	//needingUpdate := 0
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

		//has, err := hasCountry(&info, bcp47)
		//if err != nil {
		//	fatal(err)
		//}
		//if !has {
		//	fmt.Printf("missing %s", bcp47)
		//	continue
		//}

		//localeNeedsUpdate := false
		if do_text {
			if info.defBcp47 != listing.Language {
				// The description for the "base" locale is by definition
				// correct (and not in need of translating!).
				//localeNeedsUpdate = updateDescriptions(&info)
				err := updateDescriptions(&info, listing.Language)
				if err != nil {
					fatal(err)
				}
			}
		}
		//if do_images {
		//	localeNeedsUpdate =
		//		localeNeedsUpdate
		//			updateImages(*shotsDir, apei, defBcp47, bcp47) || localeNeedsUpdate
		//}

		//needsCommit = needsCommit || localeNeedsUpdate

		//log.Printf("%d calls +%d", *xlns.APCalls, *xlns.APCalls-lastCalls)
		//lastCalls = *xlns.APCalls

		//if localeNeedsUpdate {
		//	needingUpdate++
		//	//if needingUpdate >= *updateLimit {
		//	//	break
		//	//}
		//}
	}

	if info.needsCommit {
		err := apt.EditsCommit(info.service, info.packageName, info.editId)
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

//func hasCountry(ei *editInfo, bcp47) (bool, error) {
//	avail, err := ei.service.Edits.Countryavailability.
//}

var showBase = true

// updateDescription updates the description information for a package for
// each BCP-47 location it has information for.
func updateDescriptions(ei *editInfo, bcp47 string) error {
	// Get the base language listing.
	baseListing, err := ei.service.Edits.Listings.Get(
		ei.packageName, ei.editId, ei.defBcp47).Do()
	//err := xlns.ExpBackoff(func() error {
	//	var err error
	//	baseListing, err = baseLangGet.Do()
	//	return err
	//})
	if err != nil {
		return fmt.Errorf("getting edit listing for %s got %v", bcp47, err)
	}
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
	// TODO: on failure add...
	listing, err := ei.service.Edits.Listings.Get(
		ei.packageName, ei.editId, bcp47).Do()
	//err = xlns.ExpBackoff(func() error {
	//	var err error
	//	listing, err = langGet.Do()
	//	return err
	//})
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
		fmt.Printf("no changes for %s", bcp47)
		return nil
	}

	_, err = ei.service.Edits.Listings.Update(
		ei.packageName, ei.editId, bcp47, &translated).Do()
	//err = xlns.ExpBackoff(func() error {
	//	_, err := update.Do() // *Listing, error
	//	return err
	//})
	if err != nil {
		return fmt.Errorf("listing update for %s got %v", bcp47, err)
	}
	ei.needsCommit = true

	return nil
}

/*
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
