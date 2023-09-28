package entities

type Yandexoid struct {
	Login   string
	Name    string
	Surname string
	Status  *Status
}

func NewYandexoid(login, name, surname string, status *Status) *Yandexoid {
	return &Yandexoid{
		Login:   login,
		Name:    name,
		Surname: surname,
		Status:  status,
	}
}
