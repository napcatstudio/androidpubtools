// access.go
// Contains common functions for dealing with the Android Publishing API V3.
package androidpub

import (
	"context"
	"fmt"
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
	//fmt.Printf("%+v\n", appEdit)
	fmt.Printf("validated\n")

	commit := apei.Srv.Edits.Commit(apei.Package, apei.Id)
	err = ExpBackoff(func() error {
		_, err = commit.Do() // appEdit, err
		return err
	})
	if err != nil {
		return fmt.Errorf("Edits.Commit failed %v", err)
	}
	//fmt.Printf("%+v\n", appEdit)
	fmt.Printf("committed\n")
	return nil
}

*/
