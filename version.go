package vcardgen

type majorVersion int

const (
	// VersionTwo of vCard
	VersionTwo majorVersion = iota + 1
	// VersionThree of vCard
	VersionThree
	// VersionFour of vCard
	VersionFour
)

func (v majorVersion) getVersionString() string {
	switch v {
	case VersionTwo:
		{
			return "2.1"
		}
	case VersionThree:
		{
			return "3.0"
		}
	case VersionFour:
		{
			return "4.0"
		}
	}
	panic("Invalid Version, only 2-4 supported.")
}
