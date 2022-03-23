// distribution.go
// Contains constant Google Play distribution information.
// I would love to get this list from an API but have not figured out a way
// yet.
package androidpub

import (
	"fmt"
	"strings"

	xlns "github.com/napcatstudio/translate"
)

type GooglePlayDistribution struct {
	Country string
	Bcp47   string
}

// Google Play Supported Locations
// Copied from the Play Store "manage translations" dialog.
// January 26, 2022
var distribution = []GooglePlayDistribution{
	{"Afrikaans", "af"},
	{"Albanian", "sq"},
	{"Amharic", "am"},
	{"Arabic", "ar"},
	{"Armenian", "hy-AM"},
	{"Azerbaijani", "az-AZ"},
	{"Bangla", "bn-BD"},
	{"Basque", "eu-ES"},
	{"Belarusian", "be"},
	{"Bulgarian", "bg"},
	{"Burmese", "my-MM"},
	{"Catalan", "ca"},
	{"Chinese (Hong Kong)", "zh-HK"},
	{"Chinese (Simplified)", "zh-CN"},
	{"Chinese (Traditional)", "zh-TW"},
	{"Croatian", "hr"},
	{"Czech", "cs-CZ"},
	{"Danish", "da-DK"},
	{"Dutch", "nl-NL"},
	{"English (Australia)", "en-AU"},
	{"English (Canada)", "en-CA"},
	{"English (United Kingdom)", "en-GB"},
	{"English", "en-IN"},
	{"English", "en-SG"},
	{"English", "en-ZA"},
	{"Estonian", "et"},
	//{"Filipino", "fil"}, //TODO: not ISO-639?
	{"Finnish", "fi-FI"},
	{"French (Canada)", "fr-CA"},
	{"French (France)", "fr-FR"},
	{"Galician", "gl-ES"},
	{"Georgian", "ka-GE"},
	{"German", "de-DE"},
	{"Greek", "el-GR"},
	{"Gujarati", "gu"},
	//{"Hebrew", "iw-IL"}, //TODO: not ISO-639?
	{"Hindi", "hi-IN"},
	{"Hungarian", "hu-HU"},
	{"Icelandic", "is-IS"},
	{"Indonesian", "id"},
	{"Italian", "it-IT"},
	{"Japanese", "ja-JP"},
	{"Kannada", "kn-IN"},
	{"Kazakh", "kk"},
	{"Khmer", "km-KH"},
	{"Korean", "ko-KR"},
	{"Kyrgyz", "ky-KG"},
	{"Lao", "lo-LA"},
	{"Latvian", "lv"},
	{"Lithuanian", "lt"},
	{"Macedonian", "mk-MK"},
	{"Malay (Malaysia)", "ms-MY"},
	{"Malay", "ms"},
	{"Malayalam", "ml-IN"},
	{"Marathi", "mr-IN"},
	{"Mongolian", "mn-MN"},
	{"Nepali", "ne-NP"},
	{"Norwegian", "no-NO"},
	{"Persian", "fa"},
	{"Persian", "fa-AE"},
	{"Persian", "fa-AF"},
	{"Persian", "fa-IR"},
	{"Polish", "pl-PL"},
	{"Portuguese (Brazil)", "pt-BR"},
	{"Portuguese (Portugal)", "pt-PT"},
	{"Punjabi", "pa"},
	{"Romanian", "ro"},
	{"Romansh", "rm"},
	{"Russian", "ru-RU"},
	{"Serbian", "sr"},
	{"Sinhala", "si-LK"},
	{"Slovak", "sk"},
	{"Slovenian", "sl"},
	{"Spanish (Latin America)", "es-419"},
	{"Spanish (Spain)", "es-ES"},
	{"Spanish (United States)", "es-US"},
	{"Swahili", "sw"},
	{"Swedish", "sv-SE"},
	{"Tamil", "ta-IN"},
	{"Telugu", "te-IN"},
	{"Thai", "th"},
	{"Turkish", "tr-TR"},
	{"Ukrainian", "uk"},
	{"Urdu", "ur"},
	{"Vietnamese", "vi"},
	{"Zulu", "zu"}}

// getGooglePlayDistribution finds the requested struct or nil.
func getGooglePlayDistribution(country string) *GooglePlayDistribution {
	lower := strings.ToLower(country)
	for _, gd := range distribution {
		if lower == strings.ToLower(gd.Country) {
			return &gd
		}
	}
	return nil
}

// GooglePlayHasCountry returns whether or not Google distributes to the given
// country.
func GooglePlayHasCountry(country string) bool {
	gd := getGooglePlayDistribution(country)
	return gd != nil
}

// GooglePlayCountries returns a slice of countries that Google distributes to.
func GooglePlayCountries() []string {
	countries := make([]string, len(distribution), len(distribution))
	for i, dist := range distribution {
		countries[i] = dist.Country
	}
	return countries
}

// TranslateableGoogleLocales returns a list of BCP-47 locales we have
// languages for.  It excludes the defLang locale.
func TranslateableGoogleLocales(wordsDir, defLang string) ([]string, error) {
	var translateable []string
	for _, info := range distribution {
		if info.Bcp47 == defLang {
			// We can't translate X to X.
			continue
		}
		_, err := xlns.WordsGetLang(wordsDir, info.Bcp47)
		if err != nil {
			return nil, err
		}
		translateable = append(translateable, info.Bcp47)
	}
	return translateable, nil
}

// UntranslateableGoogleLocales returns a list of BCP-47 locales we don't
// have languages for.  It excludes the defLang locale.
//func UntranslateableGoogleLocales(wordsDir, defLang string) ([]string, error) {
//	var un []string
//	for _, info := range distribution {
//		if info.Bcp47 == defLang {
//			continue
//		}
// TODO: This is too strict.
//		iso639 := xlns.Iso639FromBcp47(info.Bcp47)
//		has, err := xlns.WordsHasLanguage(wordsDir, iso639)
//		if err != nil {
//			return nil, fmt.Errorf("bad words directory %s got %v",
//				wordsDir, err)
//		}
//		if !has {
//			un = append(un, info.Bcp47)
//		}
//	}
//	return un, nil
//}

// GoogleLocaleForLang tries to find a Google supported locale for the given
// language.
func GoogleLocaleForLang(lang string) (string, error) {
	trialBcp47 := fmt.Sprintf("%s-%s", lang, strings.ToUpper(lang))
	for _, info := range distribution {
		if trialBcp47 == info.Bcp47 {
			return info.Bcp47, nil
		}
	}
	for _, info := range distribution {
		if lang == xlns.Iso639FromBcp47(info.Bcp47) {
			return info.Bcp47, nil
		}
	}
	return "", fmt.Errorf("%s is not a language in a Google locale", lang)
}
