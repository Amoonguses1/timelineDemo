APP_SERVER_NAME=timeline_demo_app

build:
	make down
	docker-compose build

up:
	docker-compose up -d

exec_app:
	docker-compose exec -it ${APP_SERVER_NAME} bash

down:
	docker-compose down --rmi all --volumes

stop:
	docker-compose stop
