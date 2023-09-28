package entities

type Status struct {
	Value string
	Cell  Cell
}

func NewStatus(value string, cell Cell) *Status {
	return &Status{
		Value: value,
		Cell:  cell,
	}
}
