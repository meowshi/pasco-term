package entities

type Event struct {
	Name          string
	Cell          Cell
	Yandexoids    map[string]*Yandexoid
	LockerEventId string
	IsPlusOne     bool
}

func NewEvent(name string, cell Cell) *Event {
	return &Event{
		Name:       name,
		Cell:       cell,
		Yandexoids: make(map[string]*Yandexoid, 0),
		IsPlusOne:  true,
	}
}

func (e *Event) IsRegistered(login string) bool {
	_, ok := e.Yandexoids[login]
	return ok
}
