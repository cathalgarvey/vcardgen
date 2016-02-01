package vcardgen

// VcardAddress represents a structured address in vCard.
type VcardAddress struct {
	Label         string
	Street        string
	City          string
	StateProvince string
	PostalCode    string
	CountryRegion string
	Version       majorVersion
	Type          string
}

// Get formatted address
func (address *VcardAddress) getFormattedAddress(encodingPrefix string) string {
	var formattedAddress string

	if address.Label != "" ||
		address.Street != "" ||
		address.City != "" ||
		address.StateProvince != "" ||
		address.PostalCode != "" ||
		address.CountryRegion != "" {

		if address.Version == VersionFour {
			formattedAddress = "ADR" + encodingPrefix + ";TYPE=" + address.Type
			if address.Label != "" {
				formattedAddress = formattedAddress + ";LABEL=" + encodeString(address.Label) + `"`
			} else {
				formattedAddress = formattedAddress + ""
			}
			formattedAddress = formattedAddress + ":;;" +
				encodeString(address.Street) + ";" +
				encodeString(address.City) + ";" +
				encodeString(address.StateProvince) + ";" +
				encodeString(address.PostalCode) + ";" +
				encodeString(address.CountryRegion) + "\r\n"
		} else {
			if address.Label != "" {
				formattedAddress = "LABEL" + encodingPrefix + ";TYPE=" + address.Type + ":" + encodeString(address.Label) + "\r\n"
			}
			formattedAddress = formattedAddress + "ADR" + encodingPrefix + ";TYPE=" + address.Type + ":;;" +
				encodeString(address.Street) + ";" +
				encodeString(address.City) + ";" +
				encodeString(address.StateProvince) + ";" +
				encodeString(address.PostalCode) + ";" +
				encodeString(address.CountryRegion) + "\r\n"
		}
	}

	return formattedAddress
}
