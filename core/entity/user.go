package entity

type User struct {
	ID            string
	Username      string
	Avatar        string
	Discriminator string
	PublicFlags   int
	Flags         int
	Banner        string
	AccentColor   string
	GlobalName    string
	MFAEnabled    bool
	Locale        string
	PremiumType   int
	Email         string
	Verified      bool
}
