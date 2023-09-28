package parser

import (
	"errors"
	"strconv"

	"github.com/meowshi/pasco/internal/entities"
	"google.golang.org/api/sheets/v4"
)

type TableParser interface {
	Parse(service *sheets.Service, spreadsheetId string) *entities.Day
}

func GetSheetTitleById(service *sheets.Service, spreadsheetId string, sheetId string) (string, error) {
	spreadsheet, err := service.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		return "", err
	}
	sheetIdInt, _ := strconv.Atoi(sheetId)
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.SheetId == int64(sheetIdInt) {
			return sheet.Properties.Title, nil
		}
	}

	return "", errors.New("Листа с указанным id нет")
}

func HaveEmpty(slice []interface{}) bool {
	for _, el := range slice {
		if len(el.(string)) == 0 {
			return true
		}
	}
	return false
}
