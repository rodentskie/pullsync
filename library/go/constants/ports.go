package constants

type Ports struct {
	MainApi int
}

func Port() *Ports {
	return &Ports{
		MainApi: 8080,
	}
}
