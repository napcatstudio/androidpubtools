// googlepackage displays information about a Google Development Console
// package.
package main

import (
	"fmt"
	"log"
	"os"

	"access"
)

func usage(why string) {
	text := `serviceAccountJson package

    Uses the information in serviceAccountJson to access the Google Play
    Publising API and display information on the given package.

Where:
    serviceAccountJson is a JSON file with Google Service info and keys.
    package is the name of an APK that the service account has access to.`
	log.Fatalf("ERROR: %s.\nUSAGE:\n    %s %s",
		why,
		os.Args[0],
		text)
}

func main() {
	if len(os.Args) != 3 {
		usage("Wrong number of arguments")
	}
	serviceKeyFile := os.Args[1]
	packageId := os.Args[2]

	srv, err := xlns.GetAPService(serviceKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Get an AppEdit struct with the edit Id.
	appEdit, err := xlns.GetAPAppEdit(srv, packageId)
	if err != nil {
		log.Fatal(err)
	}
	editId := appEdit.Id

	// App Details
	editDetailsGet := srv.Edits.Details.Get(packageId, editId)
	appDetails, err := editDetailsGet.Do()
	if err != nil {
		log.Fatalf("Edits.Details.Get failed %v", err)
	}
	//fmt.Printf("%+v\n", appDetails)
	fmt.Printf("%s %s\n%s\n%s\n",
		packageId, appDetails.DefaultLanguage,
		appDetails.ContactEmail,
		appDetails.ContactWebsite)

	defLang := appDetails.DefaultLanguage

	// Tracks
	tracksList := srv.Edits.Tracks.List(packageId, editId)
	tlr, err := tracksList.Do()
	if err != nil {
		log.Fatalf("Edits.Tracks.List failed %v", err)
	}
	var myTracks []string
	for _, track := range tlr.Tracks {
		fmt.Printf("%s %v\n", track.Track, track.VersionCodes)
		myTracks = append(myTracks, track.Track)
	}

	// Images
	for _, imageType := range xlns.GoogleImageTypes {
		fmt.Printf("imageType: %s\n", imageType)
		imagesList := srv.Edits.Images.List(packageId, editId, defLang, imageType)
		ilr, err := imagesList.Do()
		if err != nil {
			log.Fatalf("Edits.Images.List (%s, %s) failed %v",
				defLang, imageType, err)
		}
		if len(ilr.Images) == 0 {
			fmt.Printf("no images\n")
			continue
		}
		for _, image := range ilr.Images {
			fmt.Printf("image.Id: %s\n", image.Id)
		}
	}

	// Listings
	listingsList := srv.Edits.Listings.List(packageId, editId)
	llr, err := listingsList.Do()
	if err != nil {
		log.Fatalf("Edits.Listings.List failed %v", err)
	}
	for _, listing := range llr.Listings {
		fmt.Printf("%s %s\n%s\n%s\n",
			listing.Language, listing.Title,
			listing.ShortDescription,
			listing.FullDescription) // also has Video
	}

	// Testers
	for _, track := range myTracks {
		fmt.Printf("track: %s\n", track)
		testersGet := srv.Edits.Testers.Get(packageId, editId, track)
		testers, err := testersGet.Do()
		if err != nil {
			log.Fatalf("Edits.Testers.Get failed %v", err)
		}
		for _, gg := range testers.GoogleGroups {
			fmt.Printf("%+v\n", gg)
		}
		for _, gpc := range testers.GooglePlusCommunities {
			fmt.Printf("%+v\n", gpc)
		}
	}
}
