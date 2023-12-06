PWD = ${CURDIR}
SERVICE_NAME = L0

# Запуск сервиса
.PHONY: build
build:
	go build -o bin/$(SERVICE_NAME)  $(PWD)/cmd/$(SERVICE_NAME)

# Запуск сервиса
.PHONY: run
run:
	go run $(PWD)/cmd/$(SERVICE_NAME)


# Запустить тесты
.PHONY: test
test:
	go test $(PWD)/... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html