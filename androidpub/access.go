// access.go
// Contains common functions for dealing with the Android Publishing API V3.
package androidpub

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/napcatstudio/translate/xlns"

	ap "google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

// GetAPService reads the service credentials from the JSON file and creates
// a new Android Publisher service with them.
func GetAPService(credentialsJson string) (*ap.Service, error) {
	ctx := context.Background()
	option := option.WithCredentialsFile(credentialsJson)
	service, err := ap.NewService(ctx, option)
	if err != nil {
		return nil, fmt.Errorf("creating new service %s got %v", credentialsJson, err)
	}
	return service, nil
}

// EditsInsert gets an edit ID for the given package.
func EditsInsert(service *ap.Service, packageName string) (string, error) {
	appEdit, err := EditsInsertAppEdit(service, packageName)
	if err != nil {
		return "", fmt.Errorf("inserting edit for %s got %v", packageName, err)
	}
	return appEdit.Id, nil
}

// EditsInsertAppEdit returns the full Android Publisher AppEdit
func EditsInsertAppEdit(service *ap.Service, packageName string) (*ap.AppEdit, error) {
	appEdit, err := service.Edits.Insert(packageName, nil).Do()
	if err != nil {
		return nil, fmt.Errorf("inserting edit for %s got %v", packageName, err)
	}
	return appEdit, nil
}

// EditsCommit commits the pending edit for the package.
func EditsCommit(service *ap.Service, packageName string, editId string) error {
	_, err := service.Edits.Commit(packageName, editId).Do()
	if err != nil {
		return fmt.Errorf("commiting edit for %s got %v", packageName, err)
	}
	return nil
}

// PackageInfo write package info to the given io.Writer.
func PackageInfo(w io.Writer, credentialsJson, packageName string) error {
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
	llr, err := service.Edits.Listings.List(packageName, editId).Do()
	if err != nil {
		return fmt.Errorf("getting %s listings got %v", packageName, err)
	}
	for _, listing := range llr.Listings {
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
	credentialsJson, packageName, wordsDir, imagesDir string,
	do_text, do_images bool) error {
	//	fmt.Printf(`credentials: %s
	//words: %s
	//images: %s
	//update text: %v
	//update images: %v
	//`,
	//		credentialsJson, wordsDir, imagesDir, do_text, do_images)

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

	listings, err := listings(service, packageName, editId)
	if err != nil {
		return err
	}

	var translateable []string
	if do_text {
		// Get the BCP-47 codes can check.
		translateable, err = TranslateableGoogleLocales(wordsDir, defBcp47)
		if err != nil {
			return err
		}
		//fmt.Printf("can update:%v\n", translateable)
		// We need to update the default locale also.
		translateable = append(translateable, defBcp47)

		// Show the BCP-47 places we aren't checking.
		//un, err := UntranslateableGoogleLocales(wordsDir, defBcp47)
		//if err == nil {
		//	fmt.Printf("can't update:%v\n", un)
		//}
	}

	// By locale.
	needsCommit := false
	for i, listing := range listings {
		// Output BCP-47.
		fmt.Printf("%s (%d/%d)\n", listing.Language, i+1, len(listings))

		if do_text {
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

			if defBcp47 == listing.Language {
				fmt.Printf("default not changing %s\n", defBcp47)
			} else {
				commit, err := updateDescriptions(
					service, editId,
					packageName, wordsDir,
					defBcp47, listing.Language)
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
		err := EditsCommit(service, packageName, editId)
		if err != nil {
			return err
		}
	}
	return nil
}

// PackageUpdateText updates a Play Store Android package text details using
// the AndroidPublisher API V3.
func PackageUpdateText(
	credentialsJson, packageName, wordsDir string) error {
	return PackageUpdate(
		credentialsJson, packageName, wordsDir, "", true, false)
}

// PackageUpdateText updates a Play Store Android package text details using
// the AndroidPublisher API V3.
func PackageUpdateImages(
	credentialsJson, packageName, imagesDir string) error {
	return PackageUpdate(
		credentialsJson, packageName, "", imagesDir, false, true)
}

// listings returns the listings currently available in the Play Store.
func listings(service *ap.Service, packageName, editId string) ([]*ap.Listing, error) {
	listings, err := service.Edits.Listings.List(
		packageName, editId).Do()
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
func updateDescriptions(
	service *ap.Service, editId,
	packageName, wordsDir,
	defBcp47, bcp47 string) (bool, error) {

	// Get the base language listing.
	baseListing, err := service.Edits.Listings.Get(
		packageName, editId, defBcp47).Do()
	if err != nil {
		return false, fmt.Errorf("getting edit listing for %s got %v", bcp47, err)
	}
	base639 := xlns.Iso639FromBcp47(defBcp47)

	// Create translation map.
	iso639 := xlns.Iso639FromBcp47(bcp47)
	xm, err := xlns.WordsXlnsMap(wordsDir, base639, iso639)
	if err != nil {
		return false, fmt.Errorf("%s %s to %s problem got %v",
			wordsDir, base639, iso639, err)
	}

	// Check if update is needed.
	// Read existing.
	listing, err := service.Edits.Listings.Get(
		packageName, editId, bcp47).Do()
	if err != nil {
		return false, fmt.Errorf("get listing for %s failed %v", bcp47, err)
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

/*
import (
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	ap "google.golang.org/api/androidpublisher/v2"
	gapi "google.golang.org/api/googleapi"
)

const maxTries = 5

var (
	nApCalls int  = 0
	APCalls  *int = &nApCalls
)

// ExpBackoff calls function ebf and if the result is a googleapi 5xx error
// it waits and calls it again (repeat).
func ExpBackoff(ebf func() error) error {
	baseDelay := 1 * time.Second
	offset := 123 * time.Millisecond
	trial := 1
	var err error = nil
	for {
		nApCalls += 1
		err = ebf()
		if err == nil {
			return nil
		}
		gapiErr, ok := err.(*gapi.Error)
		if !ok || trial >= maxTries {
			return err
		}
		if gapiErr.Code < 500 || gapiErr.Code >= 600 {
			return err
		}
		time.Sleep(baseDelay + offset)
		baseDelay *= 2
		trial += 1
	}
	return err
}

// GetAPService returns a Android Publisher Service for the given key.
func GetAPService(serviceKeyFile string) (*ap.Service, error) {
	keyInfo, err := ioutil.ReadFile(serviceKeyFile)
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(keyInfo, ap.AndroidpublisherScope)
	if err != nil {
		return nil, err
	}

	client := conf.Client(oauth2.NoContext)

	var srv *ap.Service
	err = ExpBackoff(func() error {
		var err error
		srv, err = ap.New(client)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve AndroidPublishing Client %v", err)
	}
	return srv, nil
}

// GetAPAppEdit creates a new AppEdit for the package.
func GetAPAppEdit(srv *ap.Service, packageId string) (*ap.AppEdit, error) {
	// Get an AppEdit struct with the edit Id.
	editsInsert := srv.Edits.Insert(packageId, nil)
	var appEdit *ap.AppEdit
	err := ExpBackoff(func() error {
		var err error
		appEdit, err = editsInsert.Do()
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("Edits.Insert failed %v", err)
	}
	return appEdit, nil
}

// APEditInfo combines useful info into one handy package.
type APEditInfo struct {
	Srv         *ap.Service
	Package, Id string
	Details     *ap.AppDetails
}

// GetAPEdit consolidates getting app details and the edit Id.
func GetAPEdit(srv *ap.Service, packageId string) (*APEditInfo, error) {
	appEdit, err := GetAPAppEdit(srv, packageId)
	if err != nil {
		return nil, err
	}
	editId := appEdit.Id

	// App Details
	editDetailsGet := srv.Edits.Details.Get(packageId, editId)
	var appDetails *ap.AppDetails
	err = ExpBackoff(func() error {
		var err error
		appDetails, err = editDetailsGet.Do()
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("Edits.Details.Get failed %v", err)
	}

	apei := &APEditInfo{
		Srv:     srv,
		Package: packageId,
		Id:      editId,
		Details: appDetails}
	return apei, nil
}

// APCommit validates and commits an edit.
func APCommit(apei *APEditInfo) error {
	validate := apei.Srv.Edits.Validate(apei.Package, apei.Id)
	err := ExpBackoff(func() error {
		_, err := validate.Do() // appEdit, err
		return err
	})
	if err != nil {
		return fmt.Errorf("Edits.Validate failed %v", err)
	}
	//fmt.Fprintf(w, "%+v\n", appEdit)
	fmt.Fprintf(w, "validated\n")

	commit := apei.Srv.Edits.Commit(apei.Package, apei.Id)
	err = ExpBackoff(func() error {
		_, err = commit.Do() // appEdit, err
		return err
	})
	if err != nil {
		return fmt.Errorf("Edits.Commit failed %v", err)
	}
	//fmt.Fprintf(w, "%+v\n", appEdit)
	fmt.Fprintf(w, "committed\n")
	return nil
}

*/
