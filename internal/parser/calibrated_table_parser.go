package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/meowshi/pasco/internal/entities"
	"github.com/meowshi/pasco/internal/env"
	"github.com/meowshi/pasco/internal/utils"
	"github.com/rickb777/date"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"google.golang.org/api/sheets/v4"
)

type CalibratedTableParser struct{}

func (ctb *CalibratedTableParser) Parse(service *sheets.Service, spreadsheetId string) *entities.Day {
	sheetTitle, err := GetSheetTitleById(service, env.SpreadsheetId, env.SheetId)
	if err != nil {
		panic(err)
	}
	logrus.Info("Sheet title получен.")

	var day entities.Day
	day.Date = date.Today()

loop:
	for {
		fmt.Print("Перечисли ячейки с названиями мероприятий через пробел: ")
		input := utils.ReadlnTrimmed()
		cells := strings.Split(input, " ")

		var wg sync.WaitGroup
		var mx sync.Mutex
		ch := make(chan bool, len(cells))

		for _, rawCell := range cells {
			wg.Add(1)

			go func(rawCell string) {
				defer wg.Done()

				cell := entities.NewCellFromString(rawCell)
				logrus.Infof("Ячейка с названием мероприятия: %s.", cell.ToString())

				rangee := fmt.Sprintf("%s!%s:%s", sheetTitle, rawCell, cell.AddColumn(3).Column)
				logrus.Infof("Range, включающий название мероприятия и сотрудников: %s.", rangee)

				res, err := service.Spreadsheets.Values.Get(env.SpreadsheetId, rangee).Do()
				if err != nil {
					fmt.Println("Скорее всего ты неправильно указал ячейку.")
					logrus.Infof("Скорее всего неправильно указана ячейка. Ошибка: %s", err)

					ch <- false
					return
				}
				if len(res.Values) == 0 {
					fmt.Printf("Видимо указана лишняя ячейка %s. Она не будет учитана.\n", rawCell)
					ch <- true
					return
				}

				startCell := entities.NewCellFromString(rawCell)

				event := entities.NewEvent(res.Values[0][0].(string), startCell)
				logrus.Infof("Получили инфу о мероприятии: %s.", event.Name)

				startCell.AddRow(1).AddColumn(3)
				people := res.Values[1:]
				for _, val := range people {
					if len(val) == 0 || HaveEmpty(val) {
						break
					}

					var status *entities.Status
					if len(val) == 4 {
						status = entities.NewStatus(val[3].(string), startCell)
					} else {
						status = entities.NewStatus("", startCell)
					}

					yandexoid := entities.NewYandexoid(val[0].(string), val[1].(string), val[2].(string), status)
					event.Yandexoids[yandexoid.Login] = yandexoid
					startCell.AddRow(1)
				}
				logrus.Infof("Яндесоиды, зарегестрированные на %s, перенесены.", event.Name)

				mx.Lock()
				day.Events = append(day.Events, event)
				mx.Unlock()

				ch <- true
			}(rawCell)
		}

		wg.Wait()

		for i := 0; i < len(cells); i++ {
			el := <-ch
			if !el {
				logrus.Info("Информация о каком-то мероприятии не была перенесена успешно. Повторный ввод.")
				continue loop
			}
		}

		break
	}
	logrus.Info("Все мероприятия перенесены.")

	PlusOneEventsSetup(day.Events)

	day.Print()
	day.Service = service
	day.Client = &http.Client{}
	day.SheetTitle = sheetTitle

	return &day
}

func PlusOneEventsSetup(events []*entities.Event) error {
	poeFile, err := os.ReadFile(fmt.Sprintf("%s/%s", env.PathPrefix, env.PlusOneEventsFile))
	if err != nil {
		logrus.Errorf("Не получилось прочитать файл с событиями +1: %s.", err)
		return errors.New("Не получилось прочитать файл с событиями +1.")
	}

	poe := make([]string, 0)
	err = json.Unmarshal(poeFile, &poe)
	if err != nil {
		logrus.Errorf("Не получилось unmarshall события +1: %s.", err)
		return errors.New("Видимо, формат данных в файле c событиями +1 неверный.")
	}
	slices.Sort(poe)

	for _, event := range events {
		if _, find := slices.BinarySearch(poe, event.Name); !find {
			event.IsPlusOne = false
		}
	}

	return nil
}
