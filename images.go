package vcardgen

import "fmt"

type cardImageType string

const (
	// PhotoImage is the type of images that represent the human face of a vCard
	PhotoImage cardImageType = "PHOTO"
	// LogoImage is the type of images that represent the faceless Corporatehood of a vCard
	LogoImage cardImageType = "LOGO"
)

// CardImage represents a URI pointing to an image for a Logo or Photo.
type CardImage struct {
	// Photo Type, i.e LOGO or PHOTO
	PhotoType           cardImageType
	Contents, MediaType string
	Base64              bool
	Version             majorVersion
}

// NewLogo returns a logo image from a URL.
func NewLogo(URL, MediaType string) *CardImage {
	return &CardImage{
		PhotoType: LogoImage,
		Contents:  URL,
		MediaType: MediaType,
		Base64:    false,
		Version:   VersionFour,
	}
}

// NewPhoto returns a logo image from a URL.
func NewPhoto(URL, MediaType string) *CardImage {
	return &CardImage{
		PhotoType: PhotoImage,
		Contents:  URL,
		MediaType: MediaType,
		Base64:    false,
		Version:   VersionFour,
	}
}

// Returns a formatted URL for a photo, or the photo contents.
func (img *CardImage) getFormattedPhoto() string {
	var params string
	switch img.Version {
	case VersionFour:
		{
			if img.Base64 {
				params = ";ENCODING=b;MEDIATYPE=image/"
			} else {
				params = ";MEDIATYPE=image/"
			}
		}
	case VersionThree:
		{
			if img.Base64 {
				params = ";ENCODING=b;TYPE="
			} else {
				params = ";TYPE="
			}
		}
	case VersionTwo:
		{
			if img.Base64 {
				params = ";ENCODING=BASE64;"
			} else {
				params = ";"
			}
		}
	}
	return fmt.Sprintf("%s%s%s:%s\r\n", string(img.PhotoType), params, img.MediaType, encodeString(img.Contents))
}
