package constants

// github user name and slack user id
type Users struct {
	PaulWaltersDev     string `json:"PaulWaltersDev"`
	MeganSitoyPractera string `json:"MeganSitoy-Practera"`
	Rodentskie         string `json:"rodentskie"`
}

func SlackUsers() *Users {
	return &Users{
		PaulWaltersDev:     "U06Q5GKADME",
		MeganSitoyPractera: "U06Q5GKADME",
		Rodentskie:         "U06Q5GKADME",
	}
}
