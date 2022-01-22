// distribution.go
// Contains constant Google Play distribution information.
// TODO: Make this either JSON or CSV.  Get it from an API?
package androidpub

import "strings"

type GooglePlayDistribution struct {
	country    string
	free, paid bool
	ccyPrice   string
}

// Google Play Supported Locations from
// date ?
// https://support.google.com/googleplay/android-developer/table/3541286?hl=en
// ★ - Users who purchase apps in these locations may see prices in their
//	local currency on Google Play, but the transaction will take place using
//	the Developer's Default Price and Currency.
// ☆ - Users in these locations may not download Paid apps from Google Play.
var distribution = []GooglePlayDistribution{
	{"Albania", true, true, "★"},
	{"Algeria", true, true, "DZD 109.00 - 45,000.00"},
	{"Angola", true, true, "★"},
	{"Antigua and Barbuda", true, true, "★"},
	{"Argentina", true, true, "★"},
	{"Armenia", true, true, "★"},
	{"Aruba", true, true, "★"},
	{"Australia", true, true, "AUD .99 - 550.00"},
	{"Austria", true, true, "EUR .50 - 350.00"},
	{"Azerbaijan", true, true, "★"},
	{"Bahamas", true, true, "★"},
	{"Bahrain", true, true, "USD .99 - 400.00"},
	{"Bangladesh", true, true, "BDT 80.00 - 33,000.00"},
	{"Belarus", true, true, "★"},
	{"Belgium", true, true, "EUR .50 - 350.00"},
	{"Belize", true, true, "★"},
	{"Benin", true, true, "★"},
	{"Bermuda", true, true, "USD .99 - 400.00"},
	{"Bolivia", true, true, "BOB 7.00 - 2,800.00"},
	{"Bosnia and Herzegovina", true, true, "★"},
	{"Botswana", true, true, "★"},
	{"Brazil", true, true, "BRL 0.99 - 1,500.00"},
	{"British Virgin Islands", true, true, "USD .99 - 400.00"},
	{"Bulgaria", true, true, "BGN 1.50 - 700.00"},
	{"Burkina Faso", true, true, "★"},
	{"Cambodia", true, true, "USD .99 - 400.00"},
	{"Cameroon", true, true, "★"},
	{"Canada", true, true, "CAD .99 - 500.00"},
	{"Cape Verde", true, true, "★"},
	{"Cayman Islands", true, true, "USD .99 - 400.00"},
	{"Chile", true, true, "CLP 200.00 - 270,000.00"},
	{"China", true, false, "☆"},
	{"Colombia", true, true, "COP 800.00 - 1,337,000.00"},
	{"Costa Rica", true, true, "CRC 500.00 - 270,000.00"},
	{"Cote d'Ivoire", true, true, "★"},
	{"Croatia", true, true, "HRK 6.60 - 2,700.00"},
	{"Cuba", true, false, "☆"},
	{"Cyprus", true, true, "EUR .50 - 350.00"},
	{"Czech Republic", true, true, "CZK 19.50 - 10,000.00"},
	{"Denmark", true, true, "DKK 6.00 - 2,600.00"},
	{"Dominican Republic", true, true, "★"},
	{"Ecuador", true, true, "★"},
	{"Egypt", true, true, "EGP 2.00 - 7,500.00"},
	{"El Salvador", true, true, "★"},
	{"Estonia", true, true, "EUR .50 - 350.00"},
	{"Fiji", true, true, "★"},
	{"Finland", true, true, "EUR .50 - 350.00"},
	{"France", true, true, "EUR .50 - 350.00"},
	{"Gabon", true, true, "★"},
	{"Georgia", true, true, "GEL 2.00 GEL to 1,100.00"},
	{"Germany", true, true, "EUR .50 - 350.00"},
	{"Ghana", true, true, "GHS 4.00 - 1,700.00"},
	{"Greece", true, true, "EUR .50 - 350.00"},
	{"Guatemala", true, true, "★"},
	{"Guinea-Bissau", true, true, "★"},
	{"Haiti", true, true, "★"},
	{"Honduras", true, true, "★"},
	{"Hong Kong", true, true, "HKD 7.00 - 3,100.00"},
	{"Hungary", true, true, "HUF 125.00 - 133,700.00"},
	{"Iceland", true, true, "★"},
	{"India", true, true, "INR 10.00 - 26,000.00"},
	{"Indonesia", true, true, "IDR 3,000.00 - 5,500,000.00"},
	{"Iran", true, false, "☆"},
	{"Iraq", true, true, "IQD 1,190.00 - 476,000"},
	{"Ireland", true, true, "EUR .50 - 350.00"},
	{"Israel", true, true, "ILS 3.00 - 1,337.00"},
	{"Italy", true, true, "EUR .50 - 350.00"},
	{"Jamaica", true, true, "★"},
	{"Japan", true, true, "JPY 99.00 - 48,000.00"},
	{"Jordan", true, true, "JOD 0.70 - 285"},
	{"Kazakhstan", true, true, "KZT 300.00 - 120,000.00"},
	{"Kenya", true, true, "KES 102.00 - 42,000.00"},
	{"Kuwait", true, true, "USD .99 - 400.00"},
	{"Kyrgyzstan", true, true, "★"},
	{"Laos", true, true, "★"},
	{"Latvia", true, true, "EUR .50 - 350.00"},
	{"Lebanon", true, true, "LBP 1,500.00 - 600,000.00"},
	{"Liechtenstein", true, true, "CHF .99 - 350.00"},
	{"Lithuania", true, true, "EUR .50 - 350.00"},
	{"Luxembourg", true, true, "EUR .50 - 350.00"},
	{"Macau", true, true, "MOP 7.50 to 3,250.00"},
	{"Macedonia", true, true, "★"},
	{"Malaysia", true, true, "MYR 1.00 - 1,337.00"},
	{"Mali", true, true, "★"},
	{"Malta", true, true, "★"},
	{"Mauritius", true, true, "★"},
	{"Mexico", true, true, "MXN 5.00 - 7,000.00"},
	{"Moldova", true, true, "★"},
	{"Morocco", true, true, "MAD 8.50 - 4,000.00"},
	{"Mozambique", true, true, "★"},
	{"Myanmar", true, true, "MMK 1,500 to 620,000"},
	{"Namibia", true, true, "★"},
	{"Nepal", true, true, "★"},
	{"Netherlands", true, true, "EUR .50 - 350.00"},
	{"Netherlands Antilles", true, true, "★"},
	{"New Zealand", true, true, "NZD .99 - 600.00"},
	{"Nicaragua", true, true, "★"},
	{"Niger", true, true, "★"},
	{"Nigeria", true, true, "NGN 40.00 - 80,000.00"},
	{"Norway", true, true, "NOK 6.00 - 3,370.00"},
	{"Oman", true, true, "USD .99 - 400.00"},
	{"Pakistan", true, true, "PKR 105.00 - 42,000.00"},
	{"Panama", true, true, "★"},
	{"Papua New Guinea", true, true, "★"},
	{"Paraguay", true, true, "PYG 5,700 to 2,400,000"},
	{"Peru", true, true, "PEN 0.99 - 1,337.00"},
	{"Philippines", true, true, "PHP 15.00 - 18,000.00"},
	{"Poland", true, true, "PLN 1.79 - 1,600.00"},
	{"Portugal", true, true, "EUR .50 - 350.00"},
	{"Qatar", true, true, "QAR 3.50 - 1,500.00"},
	{"Rest of the World", true, false, "☆"}, // note 1
	{"Romania", true, true, "RON 3.50 - 1,600.00"},
	{"Russia", true, true, "RUB 15.00 - 42,000.00"}, // note 2
	{"Rwanda", true, true, "★"},
	{"Saudi Arabia", true, true, "SAR 0.99 - 1,337.00"},
	{"Senegal", true, true, "★"},
	{"Serbia", true, true, "RSD 99 to 41,000"},
	{"Singapore", true, true, "SGD .99 - 550.00"},
	{"Slovakia", true, true, "EUR .50 - 350.00"},
	{"Slovenia", true, true, "EUR .50 - 350.00"},
	{"South Africa", true, true, "ZAR 3.99 - 5,500.00"},
	{"South Korea", true, true, "KRW 999 - 450,000.00"},
	{"Spain", true, true, "EUR .50 - 350.00"},
	{"Sri Lanka", true, true, "LKR 151.00 - 61,005.00"},
	{"Sudan", true, false, "☆"},
	{"Sweden", true, true, "SEK 7.00 - 3,000.00"},
	{"Switzerland", true, true, "CHF .99 - 350.00"},
	{"Taiwan", true, true, "TWD 30.00 - 13,370.00"},
	{"Tajikistan", true, true, "★"},
	{"Tanzania", true, true, "TZS 2,200.00 - 894,000.00"},
	{"Thailand", true, true, "THB 10.00 - 13,370.00"},
	{"Togo", true, true, "★"},
	{"Trinidad and Tobago", true, true, "★"},
	{"Tunisia", true, true, "★"},
	{"Turkey", true, true, "TRY 0.59 - 1,337.00"},
	{"Turkmenistan", true, true, "★"},
	{"Turks and Caicos Islands", true, true, "USD .99 - 400.00"},
	{"Uganda", true, true, "★"},
	{"Ukraine", true, true, "UAH 5.00 - 9,000.00"}, // note 2
	{"United Arab Emirates", true, true, "AED 3.50 - 1,337.00"},
	{"United Kingdom", true, true, "GBP .50 - 300.00"},
	{"United States", true, true, "USD .99 - 400.00"}, // note 3
	{"Uruguay", true, true, "★"},
	{"Uzbekistan", true, true, "★"},
	{"Venezuela", true, true, "★"},
	{"Vietnam", true, true, "VND 6,000.00 - 9,000,000.00"},
	{"Yemen", true, true, "★"},
	{"Zambia", true, true, "★"},
	{"Zimbabwe", true, true, "★"},
}

// note 1 - Locations that are not explicitly listed in this table fall under
//	the Rest of the world category. By enabling distribution to Rest of the
//	world, you automatically include all locations grouped in this category.
// note 2 - Due to recently enacted international sanctions against the
//	Crimea region, the availability of products listed above for Russia and
//	Ukraine does not apply to the Crimea region
// note 3 - Includes Puerto Rico

// getGooglePlayDistribution finds the requested struct or nil.
func getGooglePlayDistribution(country string) *GooglePlayDistribution {
	lower := strings.ToLower(country)
	for _, gd := range distribution {
		if lower == strings.ToLower(gd.country) {
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

// GooglePlayHasFreeDistribution returns whether or not Google has free app
// distribution in the given country.
func GooglePlayHasFreeDistribution(country string) bool {
	gd := getGooglePlayDistribution(country)
	if gd == nil {
		return false
	}
	return gd.free
}

// GooglePlayHasPaidDistribution returns whether or not Google has paid app
// distribution in the given country.
func GooglePlayHasPaidDistribution(country string) bool {
	gd := getGooglePlayDistribution(country)
	if gd == nil {
		return false
	}
	return gd.paid
}

// GooglePlayCountries returns a slice of countries that Google distributes to.
func GooglePlayCountries() []string {
	countries := make([]string, len(distribution), len(distribution))
	for i, dist := range distribution {
		countries[i] = dist.country
	}
	return countries
}

/*
These are languages devices support.

https://support.google.com/googleplay/android-developer/table/4419860?hl=en

Nov 16, 2019

Afrikaans	af
Amharic	am
Bulgarian	bg
Catalan	ca
Chinese (Hong Kong)	zh-HK
Chinese (PRC)	zh-CN
Chinese (Taiwan)	zh-TW
Croatian	hr
Czech	cs
Danish	da
Dutch	nl
English (UK)	en-GB
English (US)	en-US
Estonian	et
Filipino	fil
Finnish	fi
French (Canada)	fr-CA
French (France)	fr-FR
German	de
Greek	el
Hebrew	he
Hindi	hi
Hungarian	hu
Icelandic	is
Indonesian	id / in
Italian	it
Japanese	ja
Korean	ko
Latvian	lv
Lithuanian	lt
Malay	ms
Norwegian	no
Polish	pl
Portuguese (Brazil)	pt-BR
Portuguese (Portugal)	pt-PT
Romanian	ro
Russian	ru
Serbian	sr
Slovak	sk
Slovenian	sl
Spanish (Latin America)	es-419
Spanish (Spain)	es-ES
Swahili	sw
Swedish	sv
Thai	th
Turkish	tr
Ukrainian	uk
Vietnamese	vi
Zulu	zu
*/
