// packageinfo prints information about a Android Publishing package (app).
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	//ap "google.golang.org/api/androidpublisher/v3"
	//"google.golang.org/grpc/credentials"

	apta "github.com/napcatstudio/androidpubtools/androidpub"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		usage("incorrect number of arguments")
	}
	credentialsJson := args[0]
	packageName := args[1]
	service, err := apta.GetAPService(credentialsJson)
	if err != nil {
		log.Fatalf("error %v", err)
	}
	var editId string
	editId, err = apta.EditsInsert(service, packageName)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	// Details
	appDetails, err := service.Edits.Details.Get(packageName, editId).Do()
	if err != nil {
		log.Fatalf("getting %s details got %v", packageName, err)
	}
	fmt.Printf("%s %s\n%s\n%s\n",
		packageName, appDetails.DefaultLanguage,
		appDetails.ContactEmail,
		appDetails.ContactWebsite)
	defLang := appDetails.DefaultLanguage

	// Tracks
	tlr, err := service.Edits.Tracks.List(packageName, editId).Do()
	if err != nil {
		log.Fatalf("getting %s tracks got %v", packageName, err)
	}
	var myTracks []string
	fmt.Println("tracks:")
	for _, track := range tlr.Tracks {
		fmt.Printf("\t%s\n", track.Track)
		myTracks = append(myTracks, track.Track)
	}

	// Images
	for _, imageType := range apta.GooglePlayImageTypes {
		fmt.Printf("imageType: %s\n", imageType)
		ilr, err := service.Edits.Images.List(
			packageName, editId, defLang, imageType).Do()
		if err != nil {
			log.Fatalf("getting %s %s images got %v", packageName, defLang, err)
		}
		if len(ilr.Images) == 0 {
			fmt.Printf("\tno images\n")
			continue
		}
		for _, image := range ilr.Images {
			fmt.Printf("\timage.Id: %s\n", image.Id)
		}
	}

	// Listings
	llr, err := service.Edits.Listings.List(packageName, editId).Do()
	if err != nil {
		log.Fatalf("getting %s listings got %v", packageName, err)
	}
	for _, listing := range llr.Listings {
		fmt.Printf("%s %s\n%s\n%s\n",
			listing.Language, listing.Title,
			listing.ShortDescription,
			listing.FullDescription)
	}
}

func usage(why string) {
	text := `credentialsJson packageName

    Uses the service information in credentialsJson to access the Google Play
    Publising API and display information on the given package (app).

Where:
    credentialsJson is a JSON file with Google Service info and keys.
    packageName is the name of an APK that the service account has access to.
	
Example:
    packageinfo yourServiceKey.json com.yoursite.yourapp
`
	log.Fatalf("ERROR: %s.\nUSAGE:\n    %s %s",
		why,
		os.Args[0],
		text)
}
