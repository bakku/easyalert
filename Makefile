build:
	go build -o cmd/easyalert/easyalert cmd/easyalert/main.go

docker:
	docker build --rm -t easyalert .

test:
	docker-compose run app sh -c 'DATABASE_URL=$$DATABASE_TEST_URL && gom init && gom migrate && go test ./...'
