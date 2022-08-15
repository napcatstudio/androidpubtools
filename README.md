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

    androidpkg is a tool for managing Play Store packages.

    It can update the Play Store country text and images.  It uses a meaning
    ordered words system for text.  It uses a directory hierarchy for images.
    If a text translation is too long it used the 'sub' file, if provided, for
    alternative translation text.

    Usage:
        androidpkg [flags..] command packageName

    The commands are:
        info
        Lookup information about packageName.
        update
        Update packageName images and text.
        images
        Update packageName images using the files in images.
        text
        Update packageName text using the files in words.

    -credentials string
            Google Play Developer service credentials. (default "credentials.json")
    -images string
            Images directory. (default "images")
    -sub string
            Default update substitutions. (default "update.sub")
    -words string
            The directory containing the meaning ordered words files. (default "words")


## Reference

### Go modules

### github.com/napcatstudio/translate

Used for managing and updating meaning ordered words files.

### github.com/napcatstudio/androidpubtools

This tool.

### Android Publisher API

[Google Android Publisher API V3](https://github.com/googleapis/google-api-go-client/blob/main/androidpublisher/v3/androidpublisher-api.json)