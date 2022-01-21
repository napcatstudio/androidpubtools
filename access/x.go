package access

import (
	"context"
	"fmt"
	//"golang.org/x/oauth2"
	//"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
	//"google.golang.org/api/gensupport"
	//"io/ioutil"
)

// explicit reads credentials from the specified path.
//func explicit(jsonPath, projectID string) {
//	ctx := context.Background()
//	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonPath))
//	if err != nil {
//			log.Fatal(err)
//	}
//	defer client.Close()
//	fmt.Println("Buckets:")
//	it := client.Buckets(ctx, projectID)
//	for {
//			battrs, err := it.Next()
//			if err == iterator.Done {
//					break
//			}
//			if err != nil {
//					log.Fatal(err)
//			}
//			fmt.Println(battrs.Name)
//	}
//}

func X(clientSecretJson string) (*androidpublisher.Service, error) {
	ctx := context.Background()
	option := option.WithCredentialsFile(clientSecretJson)
	//androidpublisherService, err := androidpublisher.NewService(ctx, option)
	service, err := androidpublisher.NewService(ctx, option)
	if err != nil {
		return nil, fmt.Errorf("got %v creating new service from %s ", err, clientSecretJson)
	}
	return service, nil

	//bs, err := ioutil.ReadFile(clientSecretJson)
	//if err != nil {
	//	return fmt.Errorf("can't open client secret JSON file %s", clientSecretJson)
	//}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	//config, err := google.ConfigFromJSON(b, androidpublisher.)
	//if err != nil {
	//	return fmt.Errorf("format string", a ...interface{})
	//	log.Fatalf("Unable to parse client secret file to config: %v", err)
	//}
	//client := getClient(ctx, config)
	//service, err := youtube.New(client)
}
