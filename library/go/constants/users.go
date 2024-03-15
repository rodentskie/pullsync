package constants

// github user name and slack user id
type Users struct {
	PaulWaltersDev     string `json:"PaulWaltersDev"`
	MeganSitoyPractera string `json:"MeganSitoy-Practera"`
	Rodentskie         string `json:"rodentskie"`
}

func SlackUsers() *Users {
	return &Users{
		PaulWaltersDev:     "U020E8T5PC5",
		MeganSitoyPractera: "U04MFRK6350",
		Rodentskie:         "U046JKYN3BQ",
	}
}
