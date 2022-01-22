# androidpubtools

Google tools for Android Publishing using golang and Google APIs.

## Tools

### packageinfo

#### packageinfo credentialsJson packageName

    Uses the service information in credentialsJson to access the Google Play
    Publising API and display information on the given package (app).

#### Where:
    credentialsJson is a JSON file with Google Service info and keys.
    packageName is the name of an APK that the service account has access to.
	
#### Example:
    packageinfo yourServiceKey.json com.yoursite.yourapp

## Reference

go mod: github.com/napcatstudio/androidpubtools

[Google Android Publisher API V3](https://github.com/googleapis/google-api-go-client/blob/main/androidpublisher/v3/androidpublisher-api.json)