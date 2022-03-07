// images.go
// Android Publisher image constants and related information.
// updated: January 22, 2022
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
		"phoneScreenshots":     "2-8 PNG or JPEG, up to 8 MB each, 16:9 or 9:16, each side 320-3,840px.",
		"sevenInchScreenshots": "1-? PNG or JPEG, up to 8 MB each, 16:9 or 9:16, each side 320-3,840px.",
		"tenInchScreenshots":   "1-? PNG or JPEG, up to 8 MB each, 16:9 or 9:16, each side 320-3,840px.",
		"tvScreenshots":        "1-? from Tv",
		"wearScreenshots":      "1-? from Watch",
		"icon":                 "1 A transparent PNG or JPEG, up to 1 MB, 512px by 512px.",
		"featureGraphic":       "1 PNG or JPEG, up to 1MB, and 1,024px by 500px.",
		"tvBanner":             "1 24bit PNG (no alpha) 1280x720",
	}
)

// GoogleImageTypeInfo returns a string describing the image type.
func GoogleImageTypeInfo(imageType string) string {
	return imageTypeInfo[imageType]
}
