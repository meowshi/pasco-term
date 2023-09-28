package parser

import (
	"github.com/meowshi/pasco/internal/entities"
	"google.golang.org/api/sheets/v4"
)

type AutoTableParser struct{}

func (atp *AutoTableParser) Parse(service *sheets.Service, spreadsheetId string) *entities.Day {
	panic("not impl")
}
