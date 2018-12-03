.PHONY: all serve restart kill before

dep:
	govendor sync

setup-db:
	@./scripts/setup-db.sh

pull-image:
	@$(foreach image,$(IMAGES), sudo docker-compose pull $(image);)

create-image:
	sudo docker-compose create
	sudo docker-compose start

start-image:
	sudo docker-compose start

migrate:
	go run migration/main.go

init: dep pull-image create-image setup-db migrate

test:
	@./scripts/coverage.sh
	go tool cover -html=coverage.out -o coverage.html

compile:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/app main.go

