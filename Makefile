build:
	go build --tags=debug -o bin/pasco cmd/main.go

debug: build
	LOG_LEVEL=DEBUG bin/pasco

dacha:
	env GOOS=darwin GOARCH=amd64 go build -o bin/pasco_darwin cmd/main.go