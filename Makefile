go_build:
	go build -o cmd/easyalert/easyalert cmd/easyalert/main.go

build:
	docker-compose run app make

init:
	docker-compose build
	docker-compose up -d db
	# find a better solution
	sleep 5
	docker-compose run app gom init
	docker-compose run app gom migrate

migrate:
	docker-compose run app gom migrate

test:
	docker-compose run app go test ./...
