docker:
	docker build --rm -t easyalert .

build:
	go build -o cmd/easyalert/easyalert cmd/easyalert/main.go

migrate_test:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml run app gom init
	docker-compose -f docker-compose.yml -f docker-compose.test.yml run app gom migrate

test:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml run app go test ./...