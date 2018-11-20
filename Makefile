go_build:
	go build -o build/easyalert cmd/easyalert/main.go

build:
	docker-compose run app make

init:
	docker-compose up -d db
	# find a better solution
	sleep 5
	# setup dev database
	docker-compose run app gom init
	docker-compose run app gom migrate
	# setup test database
	docker-compose -f docker-compose.yml -f docker-compose.test.yml run app gom init
	docker-compose -f docker-compose.yml -f docker-compose.test.yml run app gom migrate

reset:
	docker-compose down
	docker volume rm easyalert_vendor
	docker volume rm easyalert_cmd

docker_build:
	docker-compose build

docker_clean:
	docker ps -aq --no-trunc -f status=exited | xargs docker rm

migrate:
	docker-compose run app gom migrate
	docker-compose -f docker-compose.yml -f docker-compose.test.yml run app gom migrate

test:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml run app go test ./...
