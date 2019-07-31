package admin

//Object defines admin accounts
type Object struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	ID        int64  `json:"id"`
}

type Password struct {
	HashPassword string `json:"hashPassword"`
	ID           int64  `json:"id"`
}
