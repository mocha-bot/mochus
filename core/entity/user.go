package entity

type User struct {
	ID            string
	Username      string
	Discriminator string
	Avatar        string
	Email         string
}
