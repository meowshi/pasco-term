package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/meowshi/pasco/internal/env"
	"github.com/meowshi/pasco/internal/log"
	"github.com/meowshi/pasco/internal/parser"
	"github.com/sirupsen/logrus"
)

func main() {
	log.InitLog()

	env.Init()

	service := SetupService()

	var parser parser.TableParser = &parser.CalibratedTableParser{}

	day := parser.Parse(service, os.Getenv(env.SpreadsheetId))
	day.Start()
}

func SetupService() *sheets.Service {
	ctx := context.Background()

	credentialsPath := fmt.Sprintf("%s/%s", env.PathPrefix, env.CredentialsFile)
	logrus.Infof("Credentials path: %s.", credentialsPath)

	credentials, err := os.ReadFile(credentialsPath)
	if err != nil {
		panic("Отсутсвует файл credentials.json")
	}

	config, err := google.JWTConfigFromJSON(credentials, env.Scope)
	if err != nil {
		panic("Что-то не так с credentials.json")
	}
	logrus.Info("JWT-конфиг сформирован.")

	client := config.Client(ctx)

	service, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		panic(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}
	logrus.Info("Sheets Service сформирован.")

	return service
}
