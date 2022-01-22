// images.go
// Android Publisher image constants and related information.
// see: https://developers.google.com/android-publisher/api-ref/rest/v3/AppImageType
package androidpub

var (
	GooglePlayImageTypes = []string{
		"phoneScreenshots",     // Phone screenshot.
		"sevenInchScreenshots", // Seven inch screenshot.
		"tenInchScreenshots",   // Ten inch screenshot.
		"tvScreenshots",        // TV screenshot.
		"wearScreenshots",      // Wear screenshot.
		"icon",                 // Icon.
		"featureGraphic",       // Feature graphic.
		"tvBanner",             // TV banner.
	}

	imageTypeInfo = map[string]string{
		"phoneScreenshots":     "2-8 phone 24bit PNG (no alpha) ~1080x1920",
		"sevenInchScreenshots": "1-? 7in tablet shot ~1920x1200",
		"tenInchScreenshots":   "1-? 10in tablet shot ~2560x1600",
		"tvScreenshots":        "1-? from Tv",
		"wearScreenshots":      "1-? from Watch",
		"icon":                 "1 32bit PNG (no transparency) 512x512",
		"featureGraphic":       "1 24bit PNG (no alpha) 1024x500",
		"tvBanner":             "1 24bit PNG (no alpha) 1280x720",
	}
)

// GoogleImageTypeInfo returns a string describing the image type.
func GoogleImageTypeInfo(imageType string) string {
	return imageTypeInfo[imageType]
}
