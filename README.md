# androidpubtools

Google tools for Android Publishing using golang and Google APIs.

The library has tools for updating the Play Store images and text for the 
diffent regions supported by the Play Store.

It uses meaning ordered words files for updating the text (see below).

It uses a directory hierarchy and naming convention for the images (see the
`README.md` in the images directory).

## Meaning ordered words files

A meaning ordered words file is a file which has words, in one language, based on another file in a different language. The file name specifies the language. The filename must be of the form XX.words, where XX is a ISO-639 two letter language code.

For instance:

* `en.words`
    Easy to use.
    Easy.

* `de.words`
    Einfach zu gebrauchen.
    Einfach.

* `fi.words`
    Helppokäyttöinen.
    Helppo.

Meaning ordered words files are used so that translations can be controlled
separately from updating.

## Tools

### androidpub

## Reference

### Go modules

### github.com/napcatstudio/translate

Used for managing and updating meaning ordered words files.

### github.com/napcatstudio/androidpubtools

This tool.

### Android Publisher API

[Google Android Publisher API V3](https://github.com/googleapis/google-api-go-client/blob/main/androidpublisher/v3/androidpublisher-api.json)