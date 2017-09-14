package vcardgen

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mkSimpleCard() (*Vcard, error) {
	var err error
	var addr VcardAddress
	v := New()
	v.FirstName = "Cathal"
	v.LastName = "Garvey"
	v.MiddleName = "Joseph"
	v.CellPhone = "+353876638976"    // Not real.
	v.Email = "cathal@foobarbaz.xyz" // Not real.
	v.Birthday, err = time.Parse(time.RFC822, "12 Nov 85 01:10 GMT")
	v.Note = "This is meee"
	v.Source = "Source code of vcardgen."

	addr = VcardAddress{
		Label:         "Cathal Works Inc",
		Street:        "First Street",
		City:          "Foobarcity",
		StateProvince: "FB",
		CountryRegion: "US",
		PostalCode:    "12345-678",
	}
	v.WorkAddress = &addr

	if err != nil {
		return nil, err
	}
	return v, nil
}

func TestSimpleCard(t *testing.T) {
	v, err := mkSimpleCard()
	assert.Nil(t, err)
	// Test V4
	v.Version = VersionFour
	cardStrv4 := v.GetFormattedString()
	assert.NotEmpty(t, cardStrv4)
	lines := strings.Split(strings.TrimSpace(cardStrv4), "\r\n")
	assert.Equal(t, lines[0], "BEGIN:VCARD")
	assert.Equal(t, lines[1], "VERSION:4.0")
	assert.Equal(t, lines[len(lines)-1], "END:VCARD")
	for _, line := range lines[2 : len(lines)-1] {
		switch line[:2] {
		case "FN":
			assert.Equal(t, line, "FN:Cathal Joseph Garvey")
		case "N:":
			assert.Equal(t, line, "N:Garvey;Cathal;Joseph;;")
		case "BD":
			assert.Equal(t, line, "BDAY:19851112")
		case "EM":
			assert.Equal(t, line, "EMAIL;type=HOME:cathal@foobarbaz.xyz")
		case "TE":
			assert.Equal(t, line, "TEL;VALUE=uri;TYPE=\"voice,cell\":tel:+353876638976")
		case "NO":
			assert.Equal(t, line, "NOTE:This is meee")
		case "SO":
			assert.Equal(t, line, "SOURCE:Source code of vcardgen.")
		case "AD":
			assert.Equal(t, line, "ADR;TYPE=;LABEL=Cathal Works Inc\":;;First Street;Foobarcity;FB;12345-678;US")
		}
	}
	// TODO Test V3
	// TODO Test v2
}

// Ensure that the REV element changes each second in a way that ensures difference
// between different serialisations of the same card.
func TestCardRevisioning(t *testing.T) {
	v, err := mkSimpleCard()
	assert.Nil(t, err)
	vstr := v.GetFormattedString()
	<-time.After(time.Millisecond * 1100)
	assert.NotEqual(t, vstr, v.GetFormattedString())
}
