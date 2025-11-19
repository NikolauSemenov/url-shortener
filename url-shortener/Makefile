include .docker-env
export

PROJECT ?= 'url-shortener'

help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	
## make all: команда для билда всех контейнеров проекта запуска теста и старта всех контейнеров
all: build up

## make build: команда для создания заданного контейнера проекта
build:
	docker-compose -p ${PROJECT} build

## make revision: команда для создания миграции
revision:
	migrate create -ext sql -dir db/migrations -seq ${MESSAGE}

# make-up: команда для выполнения миграций
migrate-up:
	docker-compose run --rm migrate \
	  -path=migrations \
	  -database ${DB_DSN} up

# migrate-down: команда для откатывания миграции на 1 версию назад
migrate-down:
	docker-compose run --rm migrate \
	  -path=migrations \
	  -database ${DB_DSN} down 1

## make up: команда для старта заданного проекта
up:
	docker-compose -p ${PROJECT} up -d

## make down: команда для остановки конкретного контейнера
stop:
	docker-compose stop ${PROJECT}

## make down: команда для остановки заданного контейнера
down:
	docker-compose -p ${PROJECT} down

## make logging: команда для отображения логов заданного приложения
logging:
	docker-compose logs -f --tail="50" ${PROJECT}

## make ps: команда для вывода списка всех запушенных контейнеров
ps:
	docker-compose ps

## make restart: команда для перезапуска заданного контейнера
restart:
	docker-compose restart ${PROJECT}

## make swag_init: команда для создания/обновления директории с документацией
swag_init:
	swag init -g cmd/app/main.go --output docs
