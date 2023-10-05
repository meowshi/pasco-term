package entities

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/meowshi/pasco/internal/env"
	"github.com/meowshi/pasco/internal/utils"
	"github.com/rickb777/date"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"google.golang.org/api/sheets/v4"
)

var PrinterPool []string = []string{"1", "2"}

type Day struct {
	Date   date.Date
	Events []*Event

	PrinterId string

	SheetTitle string
	Service    *sheets.Service
	Client     *http.Client
}

type LockerEvent struct {
	Id   int    `json:"id"`
	Name string `json:"description"`
}

func (d *Day) BindLockers() error {
	req, err := http.NewRequest("GET", env.LockerApi+"/events/events/", nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("start_time__gte", fmt.Sprintf("%d-%.2d-%dT00:00:00Z", d.Date.Year(), d.Date.Month(), d.Date.Day()))
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", env.LockerToken)

	res, err := d.Client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	lockerEvents := make([]LockerEvent, 0)
	err = json.Unmarshal(body, &lockerEvents)
	if err != nil {
		return err
	}

	fmt.Println("Заведены следующие локеры:")
	for i, lockerEvent := range lockerEvents {
		fmt.Printf("%d. %s\n", i+1, lockerEvent.Name)
	}

	fmt.Println("Для каждой тренировки яндексоидов напиши цифру нужного локера. Если, допустим у события нет локеров - пиши '0'.")
	for i := 0; i < len(d.Events); {
		fmt.Printf("%s: ", d.Events[i].Name)

		lockerEventId := utils.ReadlnTrimmed()
		lockerEventIdInt, err := strconv.Atoi(lockerEventId)
		if lockerEventIdInt == 0 {
			d.Events[i].LockerEventId = "-1"
			i++
			continue
		} else if err != nil || lockerEventIdInt < 1 || lockerEventIdInt > len(lockerEvents) {
			fmt.Println("Ты ввел что-то не так. Дам еще один шанс.")
			continue
		}

		d.Events[i].LockerEventId = fmt.Sprintf("%d", lockerEvents[lockerEventIdInt-1].Id)

		i++
	}

	for {
		fmt.Print("Введи номер твоего принтера: ")
		printerNumber := utils.ReadlnTrimmed()
		_, isFinded := slices.BinarySearch(PrinterPool, printerNumber)
		if !isFinded {
			fmt.Println("Я не знаю о существовании такого принтера.")
			continue
		}

		d.PrinterId = printerNumber
		break
	}

	return nil
}

func (d *Day) Start() {
	err := d.BindLockers()
	if err != nil {
		color.Red("Что-то пошло не так при связывании локеров.")
		logrus.Errorf("Ошибка при связывании локеров: %s", err)
	}

	color.Green("\nГотовы встречать сотрудников!")

	for {
		fmt.Printf("Когда сюда после пика попадет id сотрудника, нажми Enter: ")

		var id string
		fmt.Scan(&id)
		key := utils.RfidToKey(id)

		login, err := d.GetYandexoidLogin(key)
		if err != nil {
			color.Red("Яндексоид не найден\n")
			logrus.Errorf("Ошибка при получении логина яндексоида: %s.", err)
			continue
		}
		color.Green("Ого! К нам пришел: %s\n", login)

		err = d.CheckYandexoid(login)
		if err != nil {
			color.Red(err.Error())
			fmt.Println("Будем давать подарки? (y/n): ")
			input := utils.Readln()
			switch input {
			case "y\n", "\n":
				break
			case "n\n":
				continue
			default:
				fmt.Println("Бог с тобой")
				continue
			}
		}

		err = d.GiveGift(login, key)
		if err != nil {
			color.Red("Когда отдавали подарок, что-то случилось: %s\n\n", err)
			continue
		}
		color.Green("Подарок улетел!")

		fmt.Println()
	}
}
func (d *Day) PrintBraclet(event *Event, count int) error {
	json := fmt.Sprintf("{\"event_id\": \"%s\", \"printer_id\": \"%s\", \"print_count\": %d}", event.LockerEventId, d.PrinterId, count)
	req, err := http.NewRequest("POST", env.LockerApi+"/bracelets/create/", bytes.NewReader([]byte(json)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Authorization", env.LockerToken)

	res, err := d.Client.Do(req)
	if err != nil && res.StatusCode != 200 {
		for i := 0; i < len(PrinterPool); i++ {
			if PrinterPool[i] != d.PrinterId {
				json := fmt.Sprintf("{\"event_id\": \"%s\", \"printer_id\": \"%s\", \"print_count\": %d}", event.LockerEventId, PrinterPool[i], count)
				req, err := http.NewRequest("POST", env.LockerApi+"/bracelets/create/", bytes.NewReader([]byte(json)))
				if err != nil {
					return err
				}
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")
				req.Header.Set("Authorization", env.LockerToken)

				res, err = d.Client.Do(req)
				if err != nil && res.StatusCode != 200 {
					return err
				}
			}
		}
	}
	return nil
}

func (d *Day) GiveGift(login, key string) error {
	giftUrl := env.GiftCollect + login + "/"

	form := url.Values{}
	form.Add("key", key)
	form.Add("count", "1")

	req, err := http.NewRequest(http.MethodPost, giftUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("что-то пошло не так. Status: %s", res.Status))
	}

	return nil
}

func (d *Day) GetYandexoidLogin(key string) (string, error) {
	req, err := http.NewRequest("GET", env.GiftAuth, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("key", key)
	req.URL.RawQuery = q.Encode()

	res, err := d.Client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode == 400 {
		return "", errors.New("HTTP response status - 400")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var giftRes GiftRes
	err = json.Unmarshal(body, &giftRes)
	if err != nil {
		return "", err
	}

	return giftRes.LuckyLogin, nil
}

type GiftRes struct {
	LuckyLogin     string      `json:"lucky_login"`
	Himself        bool        `json:"himself"`
	CollectorLogin string      `json:"collector_login"`
	CollectedCount int         `json:"collected_count"`
	FilterUrl      string      `json:"filter_url"`
	CollectedFor   []string    `json:"collected_for"`
	CollectorFor   interface{} `json:"collector_for"`
	Key            string      `json:"key"`
}

func (d *Day) GetYandexoidEntries(login string) []*Event {
	entries := make([]*Event, 0)

	for _, event := range d.Events {
		if event.IsRegistered(login) {
			entries = append(entries, event)
		}
	}

	return entries
}

func (d *Day) AskPeopleCount() int {
	var peopleCount int
loop:
	for {
		fmt.Print("Человек пришел один (y/n): ")

		reader := bufio.NewReader(os.Stdin)

		isSolo, _ := reader.ReadString('\n')
		isSolo = strings.ToLower(isSolo)

		switch isSolo {
		case "y\n", "\n":
			peopleCount = 1
			break loop
		case "n\n":
			peopleCount = 2
			break loop
		default:
			color.Red("Да или нет? y/<ENTER> или n?")
		}
	}
	return peopleCount
}

func (d *Day) UpdateYandexoidStatus(event *Event, login string, peopleCount int) error {
	yandexoid := event.Yandexoids[login]

	statusRange := fmt.Sprintf("%s!%s", d.SheetTitle, yandexoid.Status.Cell.ToString())
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, []interface{}{peopleCount})

	_, err := d.Service.Spreadsheets.Values.Update(env.SpreadsheetId, statusRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	yandexoid.Status.Value = fmt.Sprintf("%d", peopleCount)

	return nil
}

func (d *Day) CheckYandexoid(login string) error {
	entries := d.GetYandexoidEntries(login)

	var event *Event
	switch len(entries) {
	case 0:
		return errors.New("Яндексоид никуда не записан")
	case 1:
		event = entries[0]
		color.Green("Яндексоид записан на %s.\n", event.Name)
		break
	default:
		color.Green("Яндексоид записан на несколько мероприятий.")
		for i, e := range entries {
			fmt.Printf("%d. %s\n", i+1, e.Name)
		}

		fmt.Print("Введите цифру одного мероприятия: ")
		for {
			var choice int
			fmt.Scanf("%d", &choice)

			if choice > len(entries) || choice <= 0 {
				color.New(color.FgRed).Print("Что-то ты напутал. Попробуй снова: ")
				continue
			}

			event = entries[choice-1]

			break
		}
	}

	peopleCount := 1
	if event.IsPlusOne {
		peopleCount = d.AskPeopleCount()
	}

	err := d.PrintBraclet(event, peopleCount)
	if err != nil {
		logrus.Errorf("Ошибка при печати браслета: %s.", err)
		return errors.New("Не удалось распечатать браслет.")
	}

	return d.UpdateYandexoidStatus(event, login, peopleCount)
}

func (d *Day) Print() {
	for _, event := range d.Events {
		fmt.Printf("%s Записей - %d %s\n", event.Name, len(event.Yandexoids), func() string {
			if event.IsPlusOne {
				return "Разрешен +1"
			}
			return ""
		}())

		for _, y := range event.Yandexoids {
			fmt.Printf("\t%s %s %s %s\n", y.Login, y.Name, y.Surname, y.Status.Cell.ToString())
		}
	}
}
