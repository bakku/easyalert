docker:
	docker build --rm -t easyalert .

build:
	go build -o cmd/easyalert/easyalert cmd/easyalert/main.go