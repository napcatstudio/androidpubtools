package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	ap "google.golang.org/api/androidpublisher/v3"

	apta "github.com/napcatstudio/androidpubtools/access"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		usage("incorrect number of arguments")
	}
	print(args)
	apservice, err := apta.X(args[0])
	if err != nil {
		log.Fatalf("error %v", err)
	}
	//_, err := apta.details(apservice, args[1])
	packageName := args[1]
	//apservice.Edits.Details.Get(packageName)
	//editsService := ap.NewEditsService(apservice)
	//appEdit := ap.AppEdit{}
	//editsInsertCall := editsService.Insert(packageName, &appEdit)
	var appEdit *ap.AppEdit
	appEdit, err = apservice.Edits.Insert(packageName, appEdit).Do()
	//_, err = editsInsertCall.Do()
	if err != nil {
		log.Fatalf("error %v", err)
	}
	fmt.Printf("\n%v\n", appEdit)

	var appDetails *ap.AppDetails
	appDetails, err = apservice.Edits.Details.Get(packageName, appEdit.Id).Do()
	if err != nil {
		log.Fatalf("error %v", err)
	}
	fmt.Printf("\n%v\n", appDetails)
	//detailsService := ap.NewDetailsService(apservice)
	//detailsGetCall := detailsService.Get(packageName, appEdit.Id)
	//var appDetails ap.appDetails = nil
	//appDetails, err = detailsGetCall.Do()
	//if err != nil {
	//	log.Fatalf("error %v", err)
	//}

	fmt.Println()
	//fmt.Printf("%v\n", appDetails)
	fmt.Println("it worked?")
}

func usage(why string) {
	text := `clientSecretJson serviceAccountJson packageName

	TODO: client_secret as flag?
	"client_secret.json"

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
