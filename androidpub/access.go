// access.go
// Contains common functions for dealing with the Android Publishing API V3
// access and edits.
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
