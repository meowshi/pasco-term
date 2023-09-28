package env

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	SpreadsheetId     string
	SheetId           string
	PathPrefix        string
	CredentialsFile   string
	Scope             string
	GiftAuth          string
	GiftCollect       string
	LockerToken       string
	LockerApi         string
	PlusOneEventsFile string
)

func Init() {
	var err error
	PathPrefix, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(fmt.Sprintf("я не знаю что произошло: %s", err))
	}

	godotenv.Load(fmt.Sprintf("%s/%s", PathPrefix, ".env"))

	SpreadsheetId = os.Getenv("SPREADSHEET_ID")
	SheetId = os.Getenv("SHEET_ID")
	CredentialsFile = os.Getenv("CREDENTIALS_FILE")
	Scope = os.Getenv("SCOPE")
	GiftAuth = os.Getenv("GIFT_AUTH")
	GiftCollect = os.Getenv("GIFT_COLLECT")
	LockerToken = os.Getenv("LOCKER_TOKEN")
	LockerApi = os.Getenv("LOCKER_API")
	PlusOneEventsFile = os.Getenv("PLUS_ONE_EVENTS_FILE")

	logrus.Info("Переменные окружения проинициализированы.")
}
