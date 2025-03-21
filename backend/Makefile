APP_SERVER_NAME=app-new

build:
	make down
	docker-compose build

up:
	make down
	docker-compose up -d

exec_app:
	docker-compose exec -it ${APP_SERVER_NAME} bash

down:
	docker-compose down --rmi all --volumes

stop:
	docker-compose stop


bench:
	make down
	MODE=bench docker-compose up -d


analysis:
	go run analysis.go
