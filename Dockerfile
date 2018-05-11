FROM golang:1.10-stretch

WORKDIR /go/src/github.com/bakku
COPY . .

RUN make build