package vcardgen

import (
	"fmt"
	"strings"
	"time"
)

// Vcard is a Vcard with pretty much all-optional fields.
type Vcard struct {
	// Name stuff
	FirstName, MiddleName, LastName string
	NamePrefix, NameSuffix          string
	FormattedName, NickName         string
	// Work or Vocation stuff
	Title, Role                               string
	Organization, WorkURL                     string
	WorkAddress                               *VcardAddress
	WorkPhone, WorkFax, PagerPhone, WorkEmail string
	Logo                                      *CardImage
	// Other contact stuff
	HomeAddress                          *VcardAddress
	CellPhone, HomePhone, HomeFax, Email string
	URL                                  string
	Birthday, Anniversary                time.Time
	// SocialURLs adds extension entries for anything given. This may be trash-in-trash-out, beware.
	SocialURLs map[string]string
	Photo      *CardImage
	// Gender, Source and Note are data about the card information itself or person.
	Gender, Source, Note string
	// Tag all the things!
	Categories []string
	// The vCard version to export with.
	Version majorVersion
}

// New returns a new, zero-initialised, vCard struct.
// This is the preferred way to start using vCards to avoid random non-zero fields.
func New() *Vcard {
	v := new(Vcard)
	v.SocialURLs = make(map[string]string)
	return v
}

// Return a formatted name, whether given as card.FormattedName or constructed
// from first + middle + last.
func (card *Vcard) getFormattedName() string {
	if card.FormattedName != "" {
		return card.FormattedName
	}
	formattedName := ""
	if card.FirstName != "" {
		formattedName = formattedName + card.FirstName + " "
	}
	if card.MiddleName != "" {
		formattedName = formattedName + card.MiddleName + " "
	}
	if card.LastName != "" {
		formattedName = formattedName + card.LastName + " "
	}
	return encodeString(strings.TrimSpace(formattedName))
}

func (card *Vcard) getStructuredName() string {
	return fmt.Sprintf("%s;%s;%s;%s;%s",
		encodeString(card.LastName),
		encodeString(card.FirstName),
		encodeString(card.MiddleName),
		encodeString(card.NamePrefix),
		encodeString(card.NameSuffix))
}

func (card *Vcard) makeEmailField(kind, email string) string {
	localEmail := normaliseEmail(email)
	if localEmail == "" {
		// Just go with what we were given..
		localEmail = email
	}
	email = encodeString(localEmail)
	switch card.Version {
	case VersionFour:
		{
			return fmt.Sprintf("EMAIL;type=HOME:%s\r\n", email)
		}
	case VersionThree:
		{
			return fmt.Sprintf("EMAIL;CHARSET=UTF-8;type=HOME,INTERNET:%s\r\n", email)
		}
	case VersionTwo:
		{
			return fmt.Sprintf("EMAIL;CHARSET=UTF-8;HOME;INTERNET:%s\r\n", email)
		}
	}
	panic("Bad vcard version, only 2-4 supported.")
}

func (card *Vcard) getFormattedContactNumber(semanticType, formatType, number string) string {
	if card.Version == VersionFour {
		return fmt.Sprintf("TEL;VALUE=uri;TYPE=\"%s,%s\":tel:%s\r\n", formatType, semanticType, encodeString(number))
	}
	return fmt.Sprintf("TEL;TYPE=%s,%s:%s\r\n", strings.ToUpper(semanticType), strings.ToUpper(formatType), encodeString(number))
}

// GetFormattedString returns the vCard as a string value. It only includes non-zero-values
// so most fields of the struct are considered optional.
// Format will depend on the Value field of the vCard, default is VersionFour.
// This does not currently wrap lines. Consider it TODO.
func (card *Vcard) GetFormattedString() string {
	var (
		encodingPrefix string
		fmtd           string
	)
	if card.Version == 0 {
		card.Version = VersionFour
	}
	fmtd += "BEGIN:VCARD\r\n"
	fmtd += "VERSION:" + card.Version.getVersionString() + "\r\n"
	if card.Version == VersionFour {
		encodingPrefix = ""
	} else {
		encodingPrefix = ";CHARSET=UTF-8"
	}
	// Both of the below pre-escape name fields.
	fmtd += fmt.Sprintf("FN%s:%s\r\n", encodingPrefix, card.getFormattedName())
	fmtd += fmt.Sprintf("N%s:%s\r\n", encodingPrefix, card.getStructuredName())
	if card.NickName != "" && card.Version >= VersionThree {
		fmtd += fmt.Sprintf("NICKNAME%s:%s\r\n", encodingPrefix, encodeString(card.NickName))
	}
	if card.Gender != "" {
		// Deliberately ignoring "Male" or "Female" because fuck binary gender
		fmtd += fmt.Sprintf("X-GENDER%s:%s\r\n", encodingPrefix, encodeString(card.Gender))
	}
	if !card.Birthday.IsZero() {
		fmtd += fmt.Sprintf("BDAY:%s\r\n", card.Birthday.Format("20060102"))
	}
	if !card.Anniversary.IsZero() {
		fmtd += fmt.Sprintf("ANNIVERSARY:%s\r\n", card.Anniversary.Format("20060102"))
	}
	if card.Email != "" {
		fmtd += card.makeEmailField("HOME", card.Email)
	}
	if card.WorkEmail != "" {
		fmtd += card.makeEmailField("WORK", card.WorkEmail)
	}
	if card.Logo != nil {
		card.Logo.Version = card.Version
		fmtd += card.Logo.getFormattedPhoto()
	}
	if card.Photo != nil {
		card.Photo.Version = card.Version
		fmtd += card.Photo.getFormattedPhoto()
	}
	if card.CellPhone != "" {
		fmtd += card.getFormattedContactNumber("cell", "voice", card.CellPhone)
	}
	if card.PagerPhone != "" {
		fmtd += card.getFormattedContactNumber("cell", "pager", card.CellPhone)
	}
	if card.HomePhone != "" {
		fmtd += card.getFormattedContactNumber("home", "voice", card.CellPhone)
	}
	if card.WorkPhone != "" {
		fmtd += card.getFormattedContactNumber("work", "voice", card.WorkPhone)
	}
	if card.HomeFax != "" {
		fmtd += card.getFormattedContactNumber("home", "fax", card.HomeFax)
	}
	if card.WorkFax != "" {
		fmtd += card.getFormattedContactNumber("work", "fax", card.WorkFax)
	}
	if card.HomeAddress != nil {
		card.HomeAddress.Version = card.Version
		fmtd += card.HomeAddress.getFormattedAddress(encodingPrefix)
	}
	if card.WorkAddress != nil {
		card.WorkAddress.Version = card.Version
		fmtd += card.HomeAddress.getFormattedAddress(encodingPrefix)
	}
	if card.Title != "" {
		fmtd += fmt.Sprintf("TITLE%s:%s\r\n", encodingPrefix, encodeString(card.Title))
	}
	if card.Role != "" {
		fmtd += fmt.Sprintf("ROLE%s:%s\r\n", encodingPrefix, encodeString(card.Role))
	}
	if card.Organization != "" {
		fmtd += fmt.Sprintf("ORG%s:%s\r\n", encodingPrefix, encodeString(card.Organization))
	}
	if card.URL != "" {
		fmtd += fmt.Sprintf("URL%s:%s\r\n", encodingPrefix, encodeString(card.URL))
	}
	if card.WorkURL != "" {
		fmtd += fmt.Sprintf("URL%s;type=WORK:%s\r\n", encodingPrefix, encodeString(card.WorkURL))
	}
	if card.Note != "" {
		fmtd += fmt.Sprintf("NOTE%s:%s\r\n", encodingPrefix, encodeString(card.Note))
	}
	for key, URL := range card.SocialURLs {
		fmtd += fmt.Sprintf("X-SOCIALPROFILE%s;TYPE=%s:%s\r\n", encodingPrefix, key, encodeString(URL))
	}
	cats := ""
	for _, tag := range card.Categories {
		cats += "," + encodeString(tag)
	}
	if cats != "" {
		fmtd += fmt.Sprintf("CATEGORIES:%s\r\n", strings.Trim(cats, ","))
	}
	if card.Source != "" {
		fmtd += fmt.Sprintf("SOURCE%s:%s\r\n", encodingPrefix, encodeString(card.Source))
	}
	fmtd += fmt.Sprintf("REV:%s\r\n", time.Now().Format("2006-01-02T22:04:05Z00:00"))
	fmtd += fmt.Sprintf("END:VCARD\r\n")
	return fmtd
}
