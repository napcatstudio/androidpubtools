# images

The directory `images` contains directories for each `ImageType` the *Play 
Store* supports.  In order for updating to happen correctly images must 
follow a standard naming convention (see below).

The images have different *Play Store* requirements.

## Naming convention

Either *BCP47_#.TYPE* or *ISO639_#.TYPE* where:

* *BCP47* is a BCP-47 code.
* *ISO639* is a ISO-639 code.
* *#* is the image number 0-n.  Where n varies by `ImageType`.  This can be 
  omitted if there is only one image (like `en-US.png`).
* *TYPE* is the image format.  Different `ImageType`s have differing format 
  requirements.

### Examples

* `it_0.png`
  The first italian image for a type.  Using ISO-639 and a number.
* `en-US_2.png`
  The third english United States image for a type.  Using BCP-47 and a number.
* `en-AU.jpg`
  A single english Australian image for a type.  Using BCP-47.

## Directories

### icon

A transparent PNG or JPEG, up to 1 MB, 512 px by 512 px.

### featureGraphic

Your feature graphic must be a PNG or JPEG, up to 1MB, and 1,024 px by 500 px.

### phoneScreenshots

Screenshots must be PNG or JPEG, up to 8 MB each, 16:9 or 9:16 aspect ratio, with each side between 320 px and 3,840 px.

### sevenInchScreenshots

Screenshots must be PNG or JPEG, up to 8 MB each, 16:9 or 9:16 aspect ratio, with each side between 320 px and 3,840 px.

### tenInchScreenshots

Screenshots must be PNG or JPEG, up to 8 MB each, 16:9 or 9:16 aspect ratio, with each side between 320 px and 3,840 px.

### tvBanner

1 24bit PNG (no alpha) 1280x720.

# ##wearScreenshots

No info.

### tvScreenshots

No info.
