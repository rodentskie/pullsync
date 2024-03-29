package constants

// github user name and slack user id
type Users struct {
	PaulWaltersDev     string `json:"PaulWaltersDev"`
	MeganSitoyPractera string `json:"MeganSitoy-Practera"`
	Rodentskie         string `json:"rodentskie"`
	Trtshen            string `json:"trtshen"`
	TerenceCoder       string `json:"TerenceCoder"`
	Jazzmind           string `json:"jazzmind"`
	Shawnm0705         string `json:"shawnm0705"`
	Sasangachathumal   string `json:"sasangachathumal"`
	Sunilsbcloud       string `json:"Sunilsbcloud"`
}

func SlackUsers() *Users {
	return &Users{
		PaulWaltersDev:     "U020E8T5PC5",
		MeganSitoyPractera: "U04MFRK6350",
		Rodentskie:         "U046JKYN3BQ",
		Trtshen:            "U1GBY5XKJ",
		TerenceCoder:       "U1GJ48N5V",
		Jazzmind:           "U02TB2WV7",
		Shawnm0705:         "U0DFSQLJC",
		Sasangachathumal:   "U9JD7Q9GF",
		Sunilsbcloud:       "U016R04BE81",
	}
}
