package types

const (
	Tier0 string = "Superuser"
	Tier1 string = "Admin"
	Tier2 string = "Moderator"
	Tier3 string = "Reader"
)

type Account struct {
	ID string `json:"acc_id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Rights string `json:"rights"`
}

type SignIn struct {
	Email string `json:"email"`
	Password string `json:"password"`
}